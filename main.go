package main

import (
	"net/http"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port string `envconfig:"PORT" required:"true"`
}

var s Settings
var db *pebble.DB
var log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})

func main() {
	http.HandleFunc("/", HandleProxy)

	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

	db, err = pebble.Open("db", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db.")
	}

	port := ":" + s.Port
	log.Printf("satcors listening on %s\n", port)
	log.Fatal().Err(http.ListenAndServe(port, nil)).Msg("stopped listening")
}
