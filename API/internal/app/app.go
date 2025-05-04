package app

import (
	"context"
	"errors"
)

type Program struct{}

type App interface {
	GenerateLink(ctx context.Context, url string) (string, error)
	RedirectLink(ctx context.Context, url string) (string, error)
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

func NewApp() App {
	return &Program{}
}

func (r *Program) GenerateLink(ctx context.Context, url string) (string, error) {
	return url, nil
}

func (r *Program) RedirectLink(ctx context.Context, key string) (string, error) {
	return key, nil
}
