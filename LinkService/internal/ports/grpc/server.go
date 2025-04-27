package grpc

import (
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	Register(server)
	return server
}
