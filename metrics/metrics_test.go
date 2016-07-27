package passengermetrics

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	tmpdir      = "/tmp"
	exampledata = `
<?xml version="1.0" encoding="iso8859-1" ?>
<info version="3"><passenger_version>5.0.7</passenger_version><process_count>2</process_count><max>6</max><capacity_used>4</capacity_used><get_wait_list_size>0</get_wait_list_size><supergroups><supergroup><name>/asdf</name><state>READY</state><get_wait_list_size>0</get_wait_list_size><capacity_used>1</capacity_used><group default="true"><name>/asdf;&#35;default</name><component_name>default</component_name><app_root>/asdf</app_root><app_type>rack</app_type><environment>production</environment><uuid>asdfasdf</uuid><enabled_process_count>0</enabled_process_count><disabling_process_count>0</disabling_process_count><disabled_process_count>0</disabled_process_count><capacity_used>1</capacity_used><get_wait_list_size>1</get_wait_list_size><disable_wait_list_size>0</disable_wait_list_size><processes_being_spawned>1</processes_being_spawned><spawning/><life_status>ALIVE</life_status><options><app_root>/asdf</app_root><app_group_name>/asdf;production&#41;</app_group_name><app_type>rack</app_type><start_command>/usr/bin/ruby2.1&#9;/usr/share/passenger/helper-scripts/rack-loader.rb</start_command><startup_file>config.ru</startup_file><process_title>Passenger RubyApp</process_title><log_level>3</log_level><start_timeout>90000</start_timeout><environment>production</environment><base_uri>/</base_uri><spawn_method>smart</spawn_method><user>root</user><default_user>root</default_user><default_group>root</default_group><ruby>/usr/bin/ruby2.1</ruby><python>python</python><nodejs>node</nodejs><logging_agent_address>unix:/tmp/passenger.asdf/agents.s/logging</logging_agent_address><logging_agent_username>logging</logging_agent_username><logging_agent_password>asdf</logging_agent_password><debugger>false</debugger><analytics>false</analytics><group_secret>asdf</group_secret><min_processes>1</min_processes><max_processes>0</max_processes><max_preloader_idle_time>-1</max_preloader_idle_time><max_out_of_band_work_instances>1</max_out_of_band_work_instances></options><processes></processes></group></supergroup></supergroups></info>
`
)

func testingConfig() *PassengerCollection {
	p := &PassengerCollection{}
	return p
}

func writeTestFile(filename string, contents []byte) error {
	writepath := filepath.Join(tmpdir, filename)

	handle, err := os.OpenFile(writepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("Could not open local file for writing!: %v", err)
	}
	defer handle.Close()

	_, err = handle.Write(contents)
	if err != nil {
		return fmt.Errorf("Could not write data to file!: %v", err)
	}
	return nil
}

func TestRunPassengerFail(t *testing.T) {
	p := testingConfig()
	p.PassengerPath = "/dev/null"
	err := p.RunPassengerStatus()
	if err == nil {
		t.Errorf("Running passenger from a bad path was successful! %v", err)
	}
}

func TestRunPassengerFailNotFound(t *testing.T) {
	p := testingConfig()
	p.PassengerPath = "/tmp/asdfasdfasdf1234"
	err := p.RunPassengerStatus()
	if err == nil {
		t.Errorf("Running passenger from a bad path was successful! %v", err)
	}
}

func TestRunPassengerSuccess(t *testing.T) {
	p := testingConfig()

	contents := []byte("#!/bin/bash\necho 'test'")
	writeTestFile("passengertest", contents)

	p.PassengerPath = filepath.Join(tmpdir, "passengertest")
	err := p.RunPassengerStatus()
	if err != nil {
		t.Errorf("Running passenger failed! %v", err)
	}

	if string(p.RawOutput) != "test\n" {
		t.Errorf("Input of passenger status does not equal output! %v", err)
	}
}

func TestParseOutputSuccess(t *testing.T) {
	p := testingConfig()
	p.RawOutput = []byte("<xml>yay</xml>")
	if err := p.ParseRawOutput(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestParseOutputFail(t *testing.T) {
	p := testingConfig()
	p.RawOutput = []byte("yay")
	if err := p.ParseRawOutput(); err == nil {
		t.Errorf("%s", err)
	}
}

func TestParseOutputValues(t *testing.T) {
	p := testingConfig()
	p.RawOutput = []byte(exampledata)
	p.ParseRawOutput()
	if p.ParsedOutput.PassengerVersion != "5.0.7" {
		t.Errorf("Bad parse of passenger version %v", p.ParsedOutput.PassengerVersion)
	}

	if p.ParsedOutput.QueueLength != 0 {
		t.Errorf("Bad parse of queue length %v", p.ParsedOutput.QueueLength)
	}

	if p.ParsedOutput.ProcessCount != 2 {
		t.Errorf("Bad parse of process count %v", p.ParsedOutput.ProcessCount)
	}

}
