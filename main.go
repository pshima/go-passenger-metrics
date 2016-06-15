package main

import (
	"fmt"
	statsite "github.com/armon/go-metrics"
	"github.com/pshima/go-passenger-metrics/metrics"
	"log"
	"os"
	"time"
)

func main() {
	statsite_host := os.Getenv("GRAPHITE_HOST")
	statsite_port := os.Getenv("GRAPHITE_PORT")
	if statsite_host == "" {
		log.Fatal("GRAPHITE_HOST is empty, exiting")
	}
	if statsite_port == "" {
		log.Fatal("GRAPHITE_PORT is empty, exiting")
	}
	statsite_addr := fmt.Sprintf("%s:%s", statsite_host, statsite_port)
	log.Printf("Starting metrics output to %s", statsite_addr)
	sink, err := statsite.NewStatsiteSink(statsite_addr)
	if err != nil {
		log.Fatalf("Error connecting", err)
	}
	if _, err := statsite.NewGlobal(statsite.DefaultConfig("go-passenger-metrics"), sink); err != nil {
		log.Fatalf("Error starting metrics layer")
	}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		run()
	}

}

func run() {
	p := &passengermetrics.PassengerCollection{}
	p.RunPassengerStatus()
	err := p.ParseRawOutput()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Queue Depth: %v", p.ParsedOutput.QueueLength)

	statsite.SetGauge([]string{"passenger-queue-depth"}, float32(p.ParsedOutput.QueueLength))
}
