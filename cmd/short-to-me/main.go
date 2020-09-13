package main

import (
	"github.com/gsiragusa/short-to-me/api"
	"github.com/gsiragusa/short-to-me/config"
	"github.com/gsiragusa/short-to-me/database"
	"github.com/gsiragusa/short-to-me/server"
	"github.com/gsiragusa/short-to-me/shortener"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger
	lgr := logrus.New()

	// configuration
	conf, err := config.Configure()
	if err != nil {
		lgr.WithError(err).Fatal("unable to load configuration")
	}

	// database
	store, err := database.NewMongoClient(conf)
	if err != nil {
		lgr.WithError(err).Fatal("unable to connect to Mongo")
	}

	// services
	shortenSvc := shortener.NewService(lgr, conf, store)

	// server
	srv := server.New(lgr, conf, api.NewAPI(lgr, conf, shortenSvc))

	if err := srv.ListenAndServe(); err != nil {
		lgr.WithError(err).Fatal("error starting server")
	}
}
