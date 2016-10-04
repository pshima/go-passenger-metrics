package main

import (
	"fmt"
	"log"
	"os"
	"time"

	statsite "github.com/armon/go-metrics"
	"github.com/pshima/go-passenger-metrics/metrics"
)

const (
	passengerPath string = "/usr/sbin/passenger-status"
	appName       string = "go-passenger-metrics"
)

var quiet bool

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-quiet" {
			quiet = true
		}
	}

	statsiteHost := os.Getenv("GRAPHITE_HOST")
	statsitePort := os.Getenv("GRAPHITE_PORT")
	if statsiteHost == "" {
		log.Fatal("GRAPHITE_HOST is empty, exiting")
	}
	if statsitePort == "" {
		log.Fatal("GRAPHITE_PORT is empty, exiting")
	}
	statsiteAddr := fmt.Sprintf("%s:%s", statsiteHost, statsitePort)
	log.Printf("%s Starting metrics output to %s", appName, statsiteAddr)
	sink, err := statsite.NewStatsiteSink(statsiteAddr)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	if _, err := statsite.NewGlobal(statsite.DefaultConfig(appName), sink); err != nil {
		log.Fatalf("Error starting metrics layer")
	}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		run()
	}

}

func run() {
	p := passengermetrics.PassengerCollection{
		PassengerPath: passengerPath,
	}
	if err := p.RunPassengerStatus(); err != nil {
		log.Fatalf("%s Error running passenger status: %v", appName, err)
	}
	if err := p.ParseRawOutput(); err != nil {
		log.Fatalf("%s %v", appName, err)
	}

	if quiet != true {
		log.Printf("%s Queue Depth: %v", appName, p.ParsedOutput.QueueLength)
	}

	statsite.SetGauge([]string{"passenger-queue-depth"}, float32(p.ParsedOutput.QueueLength))
	statsite.SetGauge([]string{"passenger-stats-length"}, float32(len(p.RawOutput)))
}
