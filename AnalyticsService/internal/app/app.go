package app

import (
	"LinkTransformer/internal/adapters/repository"
	"context"
	"errors"
	"math/rand"
	"strings"
)

type Program struct {
	repository repository.Repository
}

type App interface {
	GenerateLink(ctx context.Context, url string) (string, error)
	RedirectLink(ctx context.Context, url string) (string, error)
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

const shortLinkLength = 8
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewApp(repository repository.Repository) App {
	return &Program{
		repository: repository,
	}
}

func generateShortKey(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func (r *Program) GenerateLink(ctx context.Context, url string) (string, error) {
	shortKey := generateShortKey(shortLinkLength)

	err := r.repository.SaveLink(ctx, shortKey, url)
	if err != nil {
		return "", err
	}

	return shortKey, nil
}

func (r *Program) RedirectLink(ctx context.Context, key string) (string, error) {
	url, err := r.repository.GetOriginalURL(ctx, key)
	if err != nil {
		return "", err
	}
	return url, nil
}
