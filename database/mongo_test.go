// +build integration

package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/shortener"
	"github.com/stretchr/testify/require"
)

var (
	client *Client
	ctx    context.Context
	doc    = shortener.ModelShorten{
		Id:    "RMAp1Vz",
		Url:   "http://www.test.com",
		Count: 10,
	}
)

func TestMain(m *testing.M) {
	var err error
	_ = os.Setenv("MONGO_DB_NAME", "test")
	conf, err := config.Configure()
	if err != nil {
		log.Fatal(err)
	}
	client, err = NewMongoClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func clearCollection() {
	if err := client.db.Collection(CollShortUrls).Drop(ctx); err != nil {
		panic(err)
	}
}

func addDocument(t *testing.T) {
	_, err := client.db.Collection(CollShortUrls).InsertOne(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_FindUrl(t *testing.T) {
	clearCollection()
	addDocument(t)

	res, err := client.FindUrl(ctx, doc.Url)

	require.Nil(t, err)
	require.Equal(t, doc.Id, res.Id)
}

func TestClient_FindById(t *testing.T) {
	clearCollection()
	addDocument(t)

	res, err := client.FindById(ctx, doc.Id)

	require.Nil(t, err)
	require.Equal(t, doc.Url, res.Url)
}

func TestClient_StoreUrl(t *testing.T) {
	clearCollection()

	err := client.StoreUrl(ctx, doc)

	require.Nil(t, err)

	res, err := client.FindById(ctx, doc.Id)

	require.Nil(t, err)
	require.Equal(t, doc.Url, res.Url)
}

func TestClient_DeleteById(t *testing.T) {
	clearCollection()
	addDocument(t)

	err := client.DeleteById(ctx, doc.Id)

	require.Nil(t, err)

	res, err := client.FindById(ctx, doc.Id)

	require.Nil(t, res)
	require.NotNil(t, err)
}

func TestClient_IncrementCount(t *testing.T) {
	clearCollection()
	addDocument(t)

	_, err := client.IncrementCount(ctx, doc.Id)

	require.Nil(t, err)

	res, err := client.FindById(ctx, doc.Id)

	require.Nil(t, err)
	require.Equal(t, doc.Count+1, res.Count)
}
