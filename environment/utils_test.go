// handlers_test.go
package environment

import (
	"io/ioutil"
	"testing"
)

// TODO: test cases for invalid input

// TestParseDiskUsageBasic tests case valid df file
func TestParseDiskUsageBasic(t *testing.T) {
	content, err := ioutil.ReadFile("../test_resources/stdout")
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)

	diskUsage, _ := parseDiskUsage(&s)
	if diskUsage != 1111.0 {
		t.Errorf("json value returned unexpected: got %f want %f",
			diskUsage, 1111.0)
	}
}

// ParseLoadAverage15MinBasic tests case of valid ld file
func TestParseLoadAverage15MinBasic(t *testing.T) {
	content, err := ioutil.ReadFile("../test_resources/loadavg")
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)

	loadavg, _ := parseLoadAverage15Min(&s)

	if loadavg != 0.05 {
		t.Errorf("json value returned unexpected: got %f want %f",
			loadavg, 0.05)
	}
}

// TestParseJSONForValueBasic tests case of valid json file
func TestParseJSONForValueBasic(t *testing.T) {
	content, err := ioutil.ReadFile("../test_resources/docker_info.json")
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)

	value, _ := parseJSONForValue(&s, "OperatingSystem")

	if value != "Ubuntu 18.04.2 LTS" {
		t.Errorf("json value returned unexpected: got %s want %s",
			value, "Ubuntu 18.04.2 LTS")
	}
}

// TestParseOperatingSystemBasic tests case of OperatingSystem string with version
func TestParseOperatingSystemBasic(t *testing.T) {

	valid := "Ubuntu 18.04.2 LTS"
	os, version, _ := parseOperatingSystem(valid)

	if os != "Ubuntu" {
		t.Errorf("Host OS returned unexpected: got %s want %s",
			os, "Ubuntu")
	}

	if version != "18.04.2 LTS" {
		t.Errorf("Host OS returned unexpected: got %s want %s",
			version, "18.04.2 LTS")
	}
}

// TestParseCpuinfoForIntBasic tests case of valid cpuinfo file
func TestParseCpuinfoForIntBasic(t *testing.T) {

	content, err := ioutil.ReadFile("../test_resources/cpuinfo")
	if err != nil {
		t.Fatal(err)
	}
	s := string(content)

	cores, _ := parseCpuinfoForInt(&s, "cpu cores")

	if cores != 10 {
		t.Errorf("num cores returned unexpected: got %d want %d",
			cores, 10)
	}
}
