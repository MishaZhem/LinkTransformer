package app

import (
	"LinkTransformer/internal/adapters/repository"
	"LinkTransformer/internal/ports/kafka"
	"context"
	"errors"
	"log"
	"math/rand"
	"strings"
)

type Program struct {
	repository repository.Repository
	producer   kafka.Producer
}

type App interface {
	GenerateLink(ctx context.Context, url string) (string, error)
	RedirectLink(ctx context.Context, url string, ipAddress string, userAgent string) (string, error)
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

const shortLinkLength = 8
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewApp(repository repository.Repository, producer kafka.Producer) App {
	return &Program{
		repository: repository,
		producer:   producer,
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

func (r *Program) RedirectLink(ctx context.Context, key string, ipAddress string, userAgent string) (string, error) {
	url, err := r.repository.GetOriginalURL(ctx, key)
	if err != nil {
		return "", err
	}

	err = r.producer.SendClickEvent(ctx, key, ipAddress, userAgent)
	if err != nil {
		// fmt.Print("Error with Kafka: \n")
		log.Println("Error with Kafka:", err)
	}

	return url, nil
}
