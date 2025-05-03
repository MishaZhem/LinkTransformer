package grpc

import (
	"LinkTransformer/internal/app"

	"google.golang.org/grpc"
)

func NewGRPCServer(app app.App) *grpc.Server {
	server := grpc.NewServer()
	Register(server, app)
	return server
}
