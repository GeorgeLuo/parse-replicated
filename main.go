package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GeorgeLuo/parse-replicated/environment"
)

// map of location paths to interested values.
// TODO: function to parse machine readable schema file and output map.
// TODO: account for multi level json Format
var locations = map[string]*environment.FileParams{
	"default/proc/cpuinfo":               &environment.FileParams{ParseForParams: []string{"cpu cores"}, Format: "cpuinfo"},
	"default/commands/loadavg/loadavg":   &environment.FileParams{ParseForParams: []string{"loadAvg"}, Format: "loadavg"},
	"default/commands/df/stdout":         &environment.FileParams{ParseForParams: []string{"diskUsage"}, Format: "df"},
	"default/docker/docker_version.json": &environment.FileParams{ParseForParams: []string{"Version"}, Format: "json"},
	"default/docker/docker_info.json":    &environment.FileParams{ParseForParams: []string{"Driver", "OperatingSystem"}, Format: "json"},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("bundle not provided")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(file)

	bundledParsedValues, errors := ParamsFromTar(&locations, reader)

	if len(errors) > 0 {
		fmt.Println(errors)
	}

	b, err := json.Marshal(bundledParsedValues)
	if err != nil {
		fmt.Printf("failed marshal results, error: %s'", err.Error())
		fmt.Printf("%+v\n", bundledParsedValues)
		return
	}
	fmt.Println(string(b))
}

// ParamsFromTar takes a reader to a tar file and the schema in the form of
// fileMap and returns the values delimited by fileMap. This is the logical
// access point of the program.
func ParamsFromTar(fileMap *map[string]*environment.FileParams,
	r io.Reader) (environment.BundleParsedValues, []error) {
	err := environment.GetFromUntarFiles(fileMap, r)
	if err != nil {
		return *new(environment.BundleParsedValues), []error{err}
	}
	return environment.GenerateBundleParsedValues(*fileMap)
}
