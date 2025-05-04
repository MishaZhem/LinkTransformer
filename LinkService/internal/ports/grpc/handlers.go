package grpc

import (
	"LinkTransformer/internal/app"
	context "context"
	"errors"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	url, err := s.app.GenerateLink(ctx, req.Url)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	return urlToLinkResponse(url), nil
}

func (s *Server) RedirectLink(ctx context.Context, req *LinkRequest) (*LinkResponse, error) {
	url, err := s.app.RedirectLink(ctx, req.Url)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	return urlToLinkResponse(url), nil
}

func getStatusByError(err error) codes.Code {
	switch {
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, app.ErrBadRequest):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}
