package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/snapdocs/go-template/pkg/handlers"
)

var (
	port, host string
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // inefficient but pretty

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host = os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	http.HandleFunc("/foo", handlers.FooHandler)
	url := fmt.Sprintf("%s:%s", host, port)

	log.Info().Str("url", url).Msg("starting http listener")
	err := http.ListenAndServe(url, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("stopping http lister due to error")
	}
	log.Info().Msg("stopping http listener")
}
