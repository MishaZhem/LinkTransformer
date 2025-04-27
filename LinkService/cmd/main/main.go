package main

import (
	grpcPort "LinkTransformer/internal/ports/grpc"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})
	lis, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Printf("can't create listener: %s\n", err.Error())
		return
	}
	grpcServer := grpcPort.NewGRPCServer()

	eg.Go(func() error {
		fmt.Println("starting GRPC server")
		errCh := make(chan error)

		defer func() {
			fmt.Println("stopping GRPC server")
			grpcServer.GracefulStop()
			_ = lis.Close()
			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("GRPC server error: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("servers shutdown: %s\n", err.Error())
	}
}
