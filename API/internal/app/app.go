package app

import (
	grpcPort "LinkTransformer/internal/ports/grpc"

	"context"
	"errors"
)

type Program struct {
	linkServiceClientClient grpcPort.LinkServiceClient
}

type App interface {
	GenerateLink(ctx context.Context, url string) (string, error)
	RedirectLink(ctx context.Context, url string) (string, error)
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

func NewApp(linkServiceClientClient grpcPort.LinkServiceClient) App {
	return &Program{linkServiceClientClient: linkServiceClientClient}
}

func (r *Program) GenerateLink(ctx context.Context, url string) (string, error) {
	link, err := r.linkServiceClientClient.GenerateLink(ctx, &grpcPort.LinkRequest{
		Url: url,
	})
	if err != nil {
		return "", err
	}
	return link.Url, nil
}

func (r *Program) RedirectLink(ctx context.Context, key string) (string, error) {
	link, err := r.linkServiceClientClient.RedirectLink(ctx, &grpcPort.LinkRequest{
		Url: key,
	})
	if err != nil {
		return "", err
	}
	return link.Url, nil
}
