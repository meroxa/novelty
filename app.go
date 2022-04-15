package main

import (
	"fmt"
	// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
	"log"
	"os"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

const NoveltyContext = "meroxa-demo-1"

func (a App) Run(v turbine.Turbine) error {
	// Identify an upstream data store for your data app
	// with the `Resources` function
	db, err := v.Resources("demodb")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	rr, err := db.Records("events", nil)
	if err != nil {
		return err
	}

	// Register the Novelty Server URL as a secret so that it will
	// be available for use by the DetectAnomaly function
	err = v.RegisterSecret("NOVELTY_SERVER_URL") // makes env var available to data app
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	res, _ := v.Process(rr, DetectAnomaly{})

	// write the augmented records (including anomaly data) back
	// into the same database, but in a different table
	err = db.Write(res, "events_novelty", nil)
	if err != nil {
		return err
	}

	return nil
}

type DetectAnomaly struct{}

func (f DetectAnomaly) Process(stream []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	for i, r := range stream {
		serverURL := os.Getenv("NOVELTY_SERVER_URL")
		nClient, err := NewNoveltyClient(serverURL)
		res, err := nClient.Observe(NoveltyContext, formatObservation(r))
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		r.Payload.Set("novelty", res)
		stream[i] = r
	}
	return stream, nil
}

// formatObservation takes a map[string]interface{} and flattens it into a []string
func formatObservation(r turbine.Record) []string {
	var obs []string
	for _, v := range r.Payload {
		obs = append(obs, fmt.Sprint(v))
	}

	return obs
}
