package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

const NoveltyContext = "testing-1"

func (a App) Run(v turbine.Turbine) error {
	// Identify an upstream data store for your data app
	// with the `Resources` function
	db, err := v.Resources("noveltydb")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	rr, err := db.Records("user_activity", nil)
	if err != nil {
		return err
	}

	// Register the Novelty Server URL as a secret so that it will
	// be available for use by the DetectAnomaly function
	err = v.RegisterSecret("NOVELTY_SERVER_URL") // makes env var available to data app
	if err != nil {
		return err
	}
	err = v.RegisterSecret("NOVELTY_AUTH") // makes env var available to data app
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	res, _ := v.Process(rr, DetectAnomaly{})

	// write the augmented records (including anomaly data) back
	// into the same database, but in a different table
	err = db.Write(res, "user_activity_enriched")
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
			log.Printf("error in process: %s", err.Error())
		}

		// embed novelty score as string
		resString, err := json.Marshal(res)
		if err != nil {
			log.Printf("error marshaling novelty response: %s", err.Error())
			return nil, nil
		}
		r.Payload.Set("novelty", resString)
		stream[i] = r
	}
	return stream, nil
}

// formatObservation takes a map[string]interface{} and flattens it into a []string
func formatObservation(r turbine.Record) []string {
	country := r.Payload.Get("country").(string)
	city := r.Payload.Get("city").(string)
	email := r.Payload.Get("email").(string)
	userID := r.Payload.Get("user_id").(float64)
	tsFloat := r.Payload.Get("timestamp").(float64)
	tod, err := timeOfDay(fmt.Sprint(int(tsFloat)))
	log.Printf("tod: %+v", tod)
	if err != nil {
		log.Printf("error in formatObservation: %s", err.Error())
		return nil
	}

	obs := []string{tod, country, city, email, fmt.Sprint(userID)}

	log.Printf("obs: %+v", obs)

	return obs
}

// map timestamp to a time of day i.e. morning, afternoon, evening, night
func timeOfDay(t string) (string, error) {
	intTime, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return "", err
	}

	ts := time.Unix(intTime, 0)

	splitAfternoon := 12
	splitEvening := 17
	splitNight := 21

	if ts.Hour() < splitAfternoon {
		return "morning", nil
	}

	if ts.Hour() >= splitAfternoon && ts.Hour() < splitEvening {
		return "afternoon", nil
	}

	if ts.Hour() >= splitEvening && ts.Hour() < splitNight {
		return "evening", nil
	}

	return "night", nil
}
