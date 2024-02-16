package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	day                 = time.Hour * 24
	defaultDownloadTime = "04h00s"
)

// Configuration holds all the configuration for the server.
type Configuration struct {
	DatabaseDSN       string
	BaseURL           string
	downloadStartTime time.Time
}

func getConfiguration() Configuration {
	dsn, ok := os.LookupEnv("TWOMES_DSN")
	if !ok {
		logrus.Fatal("TWOMES_DSN was not set")
	}

	baseURL, ok := os.LookupEnv("TWOMES_BASE_URL")
	if !ok {
		logrus.Fatal("TWOMES_BASE_URL was not set")
	}

	downloadTime, ok := os.LookupEnv("TWOMES_DOWNLOAD_TIME")
	if !ok {
		logrus.Warning("TWOMES_DOWNLOAD_TIME was not set. defaulting to", defaultDownloadTime)
		downloadTime = defaultDownloadTime
	}

	duration, err := time.ParseDuration(downloadTime)
	if err != nil {
		logrus.Fatal(err)
	}

	downloadStartTime := time.Now().Truncate(day)
	downloadStartTime = downloadStartTime.Add(duration)
	// If time is in the past, add 1 day.
	if downloadStartTime.Before(time.Now()) {
		downloadStartTime = downloadStartTime.Add(day)
	}

	return Configuration{
		DatabaseDSN:       dsn,
		BaseURL:           baseURL,
		downloadStartTime: downloadStartTime,
	}
}
