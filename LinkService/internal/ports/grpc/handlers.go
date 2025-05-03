package grpc

import (
	context "context"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedLinkServiceServer
}

func Register(gRPC *grpc.Server) {
	RegisterLinkServiceServer(gRPC, &Server{})
}

func urlToLinkResponse(url string) *LinkResponse {
	return &LinkResponse{
		Url: url,
	}
}

func (s *Server) GenerateLink(ctx context.Context, req *LinkRequest) (*LinkResponse, error) {
	return urlToLinkResponse(req.Url), nil
}

func (s *Server) RedirectLink(ctx context.Context, req *LinkRequest) (*LinkResponse, error) {
	return urlToLinkResponse(req.Url), nil
}
