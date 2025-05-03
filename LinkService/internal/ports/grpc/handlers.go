package grpc

import (
	"LinkTransformer/internal/app"
	context "context"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedLinkServiceServer
	app app.App
}

func Register(gRPC *grpc.Server, app app.App) {
	RegisterLinkServiceServer(gRPC, &Server{app: app})
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
