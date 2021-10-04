package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tylerchambers/goapt/pkg/dpkg"
	"github.com/tylerchambers/goapt/pkg/repo"
)

// TODO (tylerchambers): This should properly parse the URI
// and support file,cdrom,http,ftp,copy,rsh/ssh,etc.
// For now we're focusing on HTTP because it's the most common.
func GetPackageIndices(s *repo.SourceEntry) error {
	repoURIs := []string{}

	for _, v := range s.Components {
		repoURIs = append(repoURIs, s.URI+"dists/"+s.Distribution+"/"+v+"/")
	}
	if s.RepoType == "deb" {
		for i, v := range repoURIs {

			repoURIs[i] = v + "binary-"
			if val, ok := s.Options["arch"]; ok {
				repoURIs[i] = repoURIs[i] + val
			} else {
				arch, err := dpkg.GetHostArch()
				if err != nil {
					return fmt.Errorf("could not identify host arch: %v", err)
				}
				repoURIs[i] = repoURIs[i] + arch
			}
		}
	}
	if s.RepoType == "deb-src" {
		for i, v := range repoURIs {
			repoURIs[i] = v + "source/"
		}
	}

	for i, v := range repoURIs {
		repoURIs[i] = v + "/Packages.gz"
	}

	// TODO (tylerchambers): Support all valid options
	// lang, target, pdiffs, by-hash, allow-insecure, allow-weak,
	// allow-downgrade-to-insecure, trusted, signed-by, check-valid-until,
	// valid-until-min, check-date, date-max-future, inrelease-path
	for _, v := range repoURIs {
		// TODO (tylerchambers): add logging
		// TODO (tylerchambers): store the gzipped files in tmp or something
		// then ungzip to /var/lib/goapt/lists/
		filename := strings.TrimPrefix(v, "http://")
		filename = strings.ReplaceAll(filename, "/", "_")
		err := downloadFile(filename, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
