package environment

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
)

// BundleParsedValues is the output of the program.
// Host OS: default/proc/version
// Host OS version: default/proc/version
// Number of cores: default/proc/cpuinfo
// Load average in seconds over the past 15 minutes: default/commands/loadavg/loadavg
//   0.75 0.35 0.25 (want this) 1/25 1747
// Disk usage in bytes on the root device: default/commands/df/stdout
// Docker version: default/docker/docker_version.json
// Docker storage driver: default/docker/docker_info.json
//   "Driver": "overlay2"
type BundleParsedValues struct {
	HostOS              string  `json:"host_os"`
	HostOSVersion       string  `json:"host_os_version,omitempty"`
	NumCores            int     `json:"num_cores"`
	LoadAverage         float64 `json:"load_average"`
	DiskUsage           float64 `json:"disk_usage"`
	DockerVersion       string  `json:"docker_version"`
	DockerStorageDriver string  `json:"docker_storage_driver"`
}

// FileParams encapsulates a target file by path.
type FileParams struct {
	ParseForParams []string
	Format         string
	File           *string
}

// GetFromUntarFiles takes a tar file and produces a map of file names
//  to files.
func GetFromUntarFiles(fileMap *map[string]*FileParams, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue

		default:
			if header.Typeflag == tar.TypeReg {
				if params, ok := (*fileMap)[header.Name]; ok {
					// localize file
					buf := new(bytes.Buffer)
					buf.ReadFrom(tr)
					s := buf.String()
					params.File = new(string)
					*params.File = s
				}
			}
		}
	}
}
