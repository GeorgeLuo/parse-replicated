package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// GenerateBundleParsedValues iterates through parseMap and populates a
// BundleParsedValues from the parameters.
func GenerateBundleParsedValues(parseMap map[string]*FileParams) (BundleParsedValues, []error) {
	var parsedBundle BundleParsedValues
	var errors []error

	for _, fileParam := range parseMap {

		switch fileParam.Format {
		case "loadavg":
			loadAverage, err := parseLoadAverage15Min(fileParam.File)
			if err != nil {
				errors = append(errors, err)
			} else {
				parsedBundle.LoadAverage = loadAverage
			}
			continue
		case "df":
			diskUsage, err := parseDiskUsage(fileParam.File)
			if err != nil {
				errors = append(errors, err)
			} else {
				parsedBundle.DiskUsage = diskUsage
			}
			continue
		case "json":
			for _, param := range fileParam.ParseForParams {
				value, err := parseJSONForValue(fileParam.File, param)
				if err != nil {
					errors = append(errors, err)
				} else {
					if param == "Version" {
						parsedBundle.DockerVersion = value
					} else if param == "Driver" {
						parsedBundle.DockerStorageDriver = value
					} else if param == "OperatingSystem" {
						parsedBundle.HostOS, parsedBundle.HostOSVersion, err =
							parseOperatingSystem(value)
						if err != nil {
							errors = append(errors, err)
						}
					}
				}
			}
			continue
		case "cpuinfo":
			for _, param := range fileParam.ParseForParams {
				value, err := parseCpuinfoForInt(fileParam.File, param)
				if err != nil {
					errors = append(errors, err)
				} else {
					if param == "cpu cores" {
						parsedBundle.NumCores = value
					}
				}
			}
		}
	}
	return parsedBundle, errors
}

// ParseDiskUsage parses df output and finds summation of disk usage.
func parseDiskUsage(dfOut *string) (float64, error) {
	splitByLines := strings.Split(*dfOut, "\n")
	l := len(splitByLines)
	var used float64

	// skip column titles
	for _, line := range splitByLines[1 : l-1] {
		splitBySpace := strings.Fields(line)
		u, err := strconv.ParseFloat(splitBySpace[2], 64)
		if err != nil {
			return 0, err
		}
		used += u
	}
	return used, nil
}

// ParseLoadAverage15Min returns the 3rd column of the loadavg output.
// TODO: map column to timespan for generic method
func parseLoadAverage15Min(loadavgOut *string) (float64, error) {
	loadVals := strings.Fields(*loadavgOut)
	return strconv.ParseFloat(loadVals[2], 64)
}

// ParseJSONForValue is a generic parser for key in json file.
// TODO: pass a slice of values to parse for.
func parseJSONForValue(jsonString *string, key string) (string, error) {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(*jsonString), &jsonMap)
	if val, ok := jsonMap[key].(string); ok {
		return val, nil
	}

	s := fmt.Sprintf("key %s not found", key)
	return "", errors.New(s)
}

// ParseOperatingSystem returns the Host OS and OS version from a full
// OperatingSystem denomination.
func parseOperatingSystem(hostString string) (string, string, error) {

	var idx int
	if idx = strings.IndexByte(hostString, ' '); idx < 0 {
		s := fmt.Sprintf("failed to parse OperatingSystem string: %s", hostString)
		return "", "", errors.New(s)
	}

	OS := hostString[:idx]
	Version := hostString[(idx + 1):]

	if unicode.IsLetter(rune(Version[0])) {
		s := fmt.Sprintf("failed to find OS Version, OperatingSystem: %s", hostString)
		OS = hostString
		return OS, "", errors.New(s)
	}

	return OS, Version, nil
}

// ParseCpuinfoForInt returns the queried value from the cpuinfo file.
// This function will return the total of whatever parameter is queried.
// TODO: pass a slice of values to parse for, return an interface.
func parseCpuinfoForInt(cpuinfo *string, param string) (int, error) {
	splitByLines := strings.Split(*cpuinfo, "\n")
	l := len(splitByLines)
	var total int

	for idx, line := range splitByLines[0 : l-1] {
		s := strings.Split(line, ":")
		if strings.TrimSpace(s[0]) == "cpu cores" {
			var cores int
			cores, err := strconv.Atoi(strings.TrimSpace(s[1]))
			if err != nil {
				s := fmt.Sprintf("failed to parse cpuinfo, line: %d, text: '%s'", idx, line)
				return total, errors.New(s)
			}
			total += cores
		}
	}

	return total, nil
}
