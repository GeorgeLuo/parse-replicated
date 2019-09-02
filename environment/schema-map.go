package environment

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
)

// BundleParsedValues is the output of the program.
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
	// indicates the fields of importance.
	ParseForParams []string
	// indicates the approach of the parsing function.
	Format string
	// the contents of the file
	File *string
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
