package shortener

import (
	"context"

	"github.com/gsiragusa/short-to-me/errors"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=shortener
type Service interface {
	ShortenUrl(ctx context.Context, url string) (string, *errors.Error)
	RetrieveUrl(ctx context.Context, url string) (string, *errors.Error)
	DeleteUrl(ctx context.Context, url string) *errors.Error
	CountRedirects(ctx context.Context, url string) (int64, *errors.Error)
	IncrementRedirect(ctx context.Context, id string) (string, *errors.Error)
}

type Store interface {
	StoreUrl(ctx context.Context, document interface{}) error
	FindUrl(ctx context.Context, url string) (*ModelShorten, error)
	FindById(ctx context.Context, id string) (*ModelShorten, error)
	DeleteById(ctx context.Context, id string) error
	IncrementCount(ctx context.Context, id string) (*ModelShorten, error)
}
