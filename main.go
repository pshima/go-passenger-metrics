package main

import (
	"fmt"
	statsite "github.com/armon/go-metrics"
	"github.com/pshima/go-passenger-metrics/metrics"
	"log"
	"os"
	"time"
)

const (
	passengerPath string = "/usr/sbin/passenger-status"
	appname       string = "go-passenger-metrics"
)

var quiet bool = false

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-quiet" {
			quiet = true
		}
	}

	statsite_host := os.Getenv("GRAPHITE_HOST")
	statsite_port := os.Getenv("GRAPHITE_PORT")
	if statsite_host == "" {
		log.Fatal("GRAPHITE_HOST is empty, exiting")
	}
	if statsite_port == "" {
		log.Fatal("GRAPHITE_PORT is empty, exiting")
	}
	statsite_addr := fmt.Sprintf("%s:%s", statsite_host, statsite_port)
	log.Printf("%s Starting metrics output to %s", appname, statsite_addr)
	sink, err := statsite.NewStatsiteSink(statsite_addr)
	if err != nil {
		log.Fatalf("Error connecting", err)
	}
	if _, err := statsite.NewGlobal(statsite.DefaultConfig(appname), sink); err != nil {
		log.Fatalf("Error starting metrics layer")
	}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		run()
	}

}

func run() {
	p := &passengermetrics.PassengerCollection{}
	p.PassengerPath = passengerPath
	if err := p.RunPassengerStatus(); err != nil {
		log.Fatalf("%s Error running passenger status: %v", appname, err)
	}
	if err := p.ParseRawOutput(); err != nil {
		log.Fatalf("%s %v", appname, err)
	}

	if quiet != true {
		log.Printf("%s Queue Depth: %v", appname, p.ParsedOutput.QueueLength)
	}

	statsite.SetGauge([]string{"passenger-queue-depth"}, float32(p.ParsedOutput.QueueLength))
}
