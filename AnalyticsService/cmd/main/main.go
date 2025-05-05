package main

import (
	"AnalyticsService/internal/adapters/repository"
	"AnalyticsService/internal/app"
	grpcPort "AnalyticsService/internal/ports/grpc"
	"AnalyticsService/internal/ports/kafka"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})

	config, err := pgxpool.ParseConfig("postgres://postgres:123@localhost:5433/postgres")
	if err != nil {
		logger.WithError(err).Fatalf("can't parse pgxpool config")
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.WithError(err).Fatalf("can't create new pool")
	}

	repo := repository.NewRepository(pool, logger)
	consumer := kafka.NewConsumer(
		[]string{"localhost:9092"},
		"link-clicks",
		"analytics-group",
	)
	defer consumer.Close()

	a := app.NewApp(repo, consumer)

	go a.RunConsumer(context.Background())

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
	grpcServer := grpcPort.NewGRPCServer(a)

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
