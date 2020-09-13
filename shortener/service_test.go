package shortener

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gsiragusa/short-to-me/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testUrl = "http://www.test.com"
	shortId = "RMAp1Vz"
)

func MakeTestService(t *testing.T) (Service, *MockStore) {
	log := logrus.New()
	log.Out = ioutil.Discard // silent logger

	conf, err := config.Configure()
	require.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := NewMockStore(ctrl)

	return NewService(log, conf, store), store
}

func TestService_ShortenUrl(t *testing.T) {
	svc, store := MakeTestService(t)
	ctx := context.Background()

	store.EXPECT().FindUrl(ctx, testUrl).Return(nil, errors.New(""))
	store.EXPECT().StoreUrl(ctx, gomock.Any())

	res, err := svc.ShortenUrl(ctx, testUrl)
	require.Nil(t, err)
	require.NotEmpty(t, res)
}

func TestService_RetrieveUrl(t *testing.T) {
	svc, store := MakeTestService(t)
	ctx := context.Background()

	expected := &ModelShorten{
		Id:  shortId,
		Url: testUrl,
	}

	store.EXPECT().FindById(ctx, shortId).Return(expected, nil)

	res, err := svc.RetrieveUrl(ctx, shortId)
	require.Nil(t, err)
	require.Equal(t, testUrl, res)
}

func TestService_DeleteUrl(t *testing.T) {
	svc, store := MakeTestService(t)
	ctx := context.Background()

	store.EXPECT().DeleteById(ctx, shortId)

	err := svc.DeleteUrl(ctx, shortId)
	require.Nil(t, err)
}

func TestService_CountRedirects(t *testing.T) {
	svc, store := MakeTestService(t)
	ctx := context.Background()

	expected := &ModelShorten{
		Id:    shortId,
		Url:   testUrl,
		Count: 10,
	}

	store.EXPECT().FindById(ctx, shortId).Return(expected, nil)

	res, err := svc.CountRedirects(ctx, shortId)
	require.Nil(t, err)
	require.Equal(t, int64(10), res)
}

func TestService_IncrementRedirect(t *testing.T) {
	svc, store := MakeTestService(t)
	ctx := context.Background()

	expected := &ModelShorten{
		Id:    shortId,
		Url:   testUrl,
		Count: 10,
	}

	store.EXPECT().IncrementCount(ctx, shortId).Return(expected, nil)

	res, err := svc.IncrementRedirect(ctx, shortId)
	require.Nil(t,err)
	require.Equal(t, testUrl, res)
}