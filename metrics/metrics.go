package passengermetrics

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"os/exec"
)

type PassengerCollection struct {
	RawOutput    []byte
	ParsedOutput *PassengerData
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
	// If xmllint is not installed it has this silly new line appended to the output
	// *** Tip: if you install the 'xmllint' command then the XML output will be indented.
	// While we can install xmllint, this hack will allow this to work regardless
	// https://github.com/phusion/passenger/blob/085504c2e00b6e6322bb0ec97aa5ff43d037c729/bin/passenger-status#L174
	output, err := exec.Command("passenger-status", "--show=xml", "|", "grep", "-v", "'***Tip:'").Output()
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
		return fmt.Errorf("Cannot parse input with error: %v\n", err)
	}
	p.ParsedOutput = out
	return nil
}
