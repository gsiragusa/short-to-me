package shortener

import (
	"context"
	"strings"
	"time"

	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/errors"
	"github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
)

type service struct {
	le     *logrus.Logger
	config *config.AppConfig
	store  Store
}

func NewService(le *logrus.Logger, config *config.AppConfig, store Store) Service {
	return &service{
		le:     le,
		config: config,
		store:  store,
	}
}

var (
	errorNotFound       = errors.NewErrorNotFound()
	internalServerError = errors.NewInternalServerError()
)

// service method that returns the url encoded
func (s *service) ShortenUrl(ctx context.Context, url string) (string, *errors.Error) {
	le := s.le.WithField("url", url)
	le.Info("requested short url")

	// check url was already stored
	existing, err := s.store.FindUrl(ctx, url)
	if err == nil {
		le.Infof("already existing: %s", existing.Id)
		return existing.Id, nil
	}

	// hash current timestamp for url encoding
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		le.WithError(err).Error("error creating hash")
		return "", &internalServerError
	}
	encoded, err := h.Encode([]int{int(time.Now().Unix())})
	if err != nil {
		le.WithError(err).Error("error encoding")
		return "", &internalServerError
	}

	// create the model using the encoded timestamp as id
	// it identifies the short url
	res := &ModelShorten{
		Id:  encoded,
		Url: url,
	}

	// store the object
	le.Info("store short url")
	if err := s.store.StoreUrl(ctx, res); err != nil {
		le.WithError(err).Error("unable to store short url")
		return "", &internalServerError
	}

	le.Infof("created id: %s", res.Id)
	return res.Id, nil
}

// service method that retrieves an existing extended url given a short url
func (s *service) RetrieveUrl(ctx context.Context, url string) (string, *errors.Error) {
	le := s.le.WithField("url", url)
	le.Info("requested read url")

	// get the id from the last part of the url
	split := strings.Split(url, "/")
	id := split[len(split)-1]

	// look for url
	existing, err := s.store.FindById(ctx, id)
	if err != nil {
		le.WithError(err).Error("url was not found")
		return "", &errorNotFound
	}

	le.Infof("found: %s", existing.Url)
	return existing.Url, nil
}

// service method that deletes an existing short url
func (s *service) DeleteUrl(ctx context.Context, url string) *errors.Error {
	le := s.le.WithField("url", url)
	le.Info("requested delete url")

	// get the id from the last part of the url
	split := strings.Split(url, "/")
	id := split[len(split)-1]

	// delete url
	if err := s.store.DeleteById(ctx, id); err != nil {
		le.WithError(err).Error("url was not found")
		return &errorNotFound
	}

	le.Info("url deleted")
	return nil
}

// service method that returns the count of redirects for a given url
func (s *service) CountRedirects(ctx context.Context, url string) (int64, *errors.Error) {
	le := s.le.WithField("url", url)
	le.Info("requested count redirects")

	// get the id from the last part of the url
	split := strings.Split(url, "/")
	id := split[len(split)-1]

	// look for url
	existing, err := s.store.FindById(ctx, id)
	if err != nil {
		le.WithError(err).Error("url was not found")
		return 0, &errorNotFound
	}

	le.Infof("returning count: %d", existing.Count)
	return existing.Count, nil
}

// service method that, given a short url id, increments the count of redirects
// and returns the redirect url
func (s *service) IncrementRedirect(ctx context.Context, id string) (string, *errors.Error) {
	le := s.le.WithField("id", id)
	le.Info("increment redirect count")

	existing, err := s.store.IncrementCount(ctx, id)
	if err != nil {
		le.WithError(err).Error("url was not found")
		return "", &errorNotFound
	}

	le.Info("count incremented")
	return existing.Url, nil
}
