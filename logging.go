package main

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/snapdocs/go-common/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func testDatadogAndAirbrakes() {
	log.Logger = log.Logger.With().Str("service", "nickbehrens-go-playground-v0.0.0").Logger()

	tracer.Start(
		tracer.WithEnv("dev"),
		tracer.WithService("nickbehrens-go-playground"),
		tracer.WithServiceVersion("version-0.0.1"),
	)

	importantGoStuff()

	err := ohNoAnAirbrake()

	if err != nil {
		log.Fatal().Err(err).Msg("Error with airbrake function.")
	}

	// This line stopps the tracer and flushes out the Datadog agent.
	defer tracer.Stop()
}

func importantGoStuff() {
	span := tracer.StartSpan("important.go.method")

	span.SetTag("I a tag", "im a value for that tag")

	defer span.Finish()

	log.Info().Msg("my log message to datadog")

	log.Printf("my log message %v", span)
}

func ohNoAnAirbrake() error {
	config := &logging.Config{
		ProjectId:     306382,
		ProjectKey:    "720f6524fc1a5df11d246933a1eaef03",
		Environment:   "development",
		BreakLogLevel: "info",
		LogLevel:      "info",
	}

	err := config.SetupWithZeroLog()

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to setup with ZeroLog.")
		return err
	}

	err = errors.New("nickbehrens-development-operation-failed")

	log.Fatal().Err(err).Msg("Failed my operation.")

	defer logging.Close()

	return nil
}
