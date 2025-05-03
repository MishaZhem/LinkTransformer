package app

import (
	"LinkTransformer/internal/adapters/repository"
	"errors"
)

type Program struct {
	repository repository.Repository
}

type App interface {
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

func NewApp(repository repository.Repository) App {
	return &Program{
		repository: repository,
	}
}
