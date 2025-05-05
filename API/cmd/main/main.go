package main

import (
	"LinkTransformer/internal/app"
	grpcAnalyticsService "LinkTransformer/internal/ports/grpc/AnalyticsService"
	grpcLinkService "LinkTransformer/internal/ports/grpc/LinkService"

	"LinkTransformer/internal/ports/httpgin"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// connect to GRPC server
	conn, err := grpc.DialContext(context.Background(), "localhost:1080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	linkServiceClientClient := grpcLinkService.NewLinkServiceClient(conn)
	analyticsServiceClientClient := grpcAnalyticsService.NewAnalyticsServiceClient(conn)

	a := app.NewApp(linkServiceClientClient, analyticsServiceClientClient)

	// start HTTP server
	httpServer := httpgin.NewHTTPServer(":18080", a)

	eg.Go(func() error {
		fmt.Println("starting HTTP server")
		errCh := make(chan error)

		defer func() {
			fmt.Println("stopping HTTP server")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("error on HTTP server closing occurred: %s", err.Error())
			}
			close(errCh)
		}()

		go func() {
			if err := httpServer.Listen(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("HTTP server error: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("servers shutdown: %s\n", err.Error())
	}
}
