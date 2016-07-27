package passengermetrics

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"os"
	"os/exec"
)

type PassengerCollection struct {
	RawOutput     []byte
	ParsedOutput  *PassengerData
	PassengerPath string
}

type PassengerData struct {
	PassengerVersion string    `xml:"passenger_version"`
	ProcessCount     int       `xml:"process_count"`
	QueueLength      int       `xml:"get_wait_list_size"`
	Processes        []Process `xml:"supergroups>supergroup>group>processes>process"`
}

type Process struct {
	PID       int    `xml:"pid"`
	Processed int    `xml:"processed"`
	RSS       int    `xml:"rss"`
	Uptime    string `xml:"uptime"`
}

func (p *PassengerCollection) RunPassengerStatus() error {
	if _, err := os.Stat(p.PassengerPath); os.IsNotExist(err) {
		return fmt.Errorf("Passenger not found at: %v", p.PassengerPath)
	}
	output, err := exec.Command(p.PassengerPath, "--show=xml").Output()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	p.RawOutput = output
	return nil
}

func (p *PassengerCollection) ParseRawOutput() error {
	out := &PassengerData{}
	dec := xml.NewDecoder(bytes.NewReader(p.RawOutput))
	dec.CharsetReader = charset.NewReaderLabel
	err := dec.Decode(out)
	if err != nil {
		return fmt.Errorf("Cannot parse passenger-status output with error: %v\n", err)
	}
	p.ParsedOutput = out
	return nil
}
