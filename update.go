package aptlib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// BuildRepoURIs generates valid apt repo URIs from a sources.list.
// TODO (tylerchambers): This should properly parse the URI
// and support file,cdrom,http,ftp,copy,rsh/ssh,etc.
// For now we're focusing on HTTP because it's the most common.
func BuildRepoURIs(s *SourceEntry) ([]string, error) {
	repoURIs := []string{}

	for _, v := range s.Components {
		if strings.HasSuffix("/", s.URI) {
			repoURIs = append(repoURIs, s.URI+"dists/"+s.Distribution+"/"+v+"/")
		} else {
			repoURIs = append(repoURIs, s.URI+"/dists/"+s.Distribution+"/"+v+"/")
		}
	}
	if s.RepoType == "deb" {
		for i, v := range repoURIs {

			repoURIs[i] = v + "binary-"
			if val, ok := s.Options["arch"]; ok {
				repoURIs[i] = repoURIs[i] + val
			} else {
				arch, err := GetHostArch()
				if err != nil {
					return nil, fmt.Errorf("could not identify host arch: %v", err)
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

	// TODO (tylerchambers): Support all valid options
	// lang, target, pdiffs, by-hash, allow-insecure, allow-weak,
	// allow-downgrade-to-insecure, trusted, signed-by, check-valid-until,
	// valid-until-min, check-date, date-max-future, inrelease-path
	// TODO (tylerchambers): add logging
	// TODO (tylerchambers): store the gzipped files in tmp or something
	// then ungzip to /var/lib/goapt/lists/
	return repoURIs, nil
}

// DownloadIndicies concurrently downloads from URIs.
// cons is the number of simultanerous connections. If it's < 0 a new connection
// is established for each URI.
func DownloadIndicies(repoURIs []string) error {

	var wg sync.WaitGroup

	for i, v := range repoURIs {
		wg.Add(1)
		repoURIs[i] = v + "/Packages.gz"
		fmt.Println(repoURIs[i])
		filename := strings.TrimPrefix(v, "http://")
		filename = strings.ReplaceAll(filename, "/", "_")
		fmt.Printf("downloading: %s\tto: %s\n", repoURIs[i], "/tmp/goapt/"+filename)
		go downloadFile("/tmp/goapt/indexgzs/"+filename, repoURIs[i], &wg)
	}
	wg.Wait()
	return nil
}

func downloadFile(filepath string, url string, wg *sync.WaitGroup) error {
	// Get the data
	fmt.Printf("get: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("get err: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("create err: %v", err)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	fmt.Printf("write err: %v", err)
	wg.Done()
	return err
}

func UnpackIndicies() error {
	sources, err := ioutil.ReadDir("/tmp/goapt/indexgzs")
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, file := range sources {
		wg.Add(1)
		go readAndExtract(file, &wg)
	}

	wg.Wait()
	return nil
}

func readAndExtract(file fs.FileInfo, wg *sync.WaitGroup) error {
	out, err := os.Create("/tmp/goapt/" + file.Name())
	if err != nil {
		return err
	}
	fbytes, err := os.ReadFile("/tmp/goapt/indexgzs/" + file.Name())
	if err != nil {
		return err
	}
	err = gunzipWrite(out, fbytes)
	if err != nil {
		return err
	}
	wg.Done()
	return nil
}

func gunzipWrite(w io.Writer, data []byte) error {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
