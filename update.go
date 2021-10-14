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

type result struct {
	message string
	err     error
}

// BuildRepoURIs generates valid apt repo URIs from a sources.list.
// TODO (tylerchambers): This should properly parse the URI
// and support file,cdrom,http,ftp,copy,rsh/ssh,etc.
// For now we're focusing on HTTP because it's the most common.
func (c *Client) BuildRepoURIs() {
	for _, s := range c.SourceEntries {
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
					arches := strings.Split(val, ",")
					base := repoURIs[i]
					for j, v := range arches {
						if j == 0 {
							repoURIs[i] = base + v
						} else {
							repoURIs = append(repoURIs, base+v)
						}
					}
				} else {
					arch := c.Arch
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
		// then ungzip to
		c.RepoURIs = append(c.RepoURIs, repoURIs...)

	}
}

// DownloadIndiciesAsync concurrently downloads from URIs.
// cons is the number of simultaneous connections. If it's < 0 a new connection
// is established for each URI.
func (c *Client) DownloadIndiciesAsync() error {
	c.InfoLogger.Println("starting index download.")
	err := os.MkdirAll(c.IndexGZStaging, os.ModePerm)
	if err != nil {
		c.ErrorLogger.Fatalf("could not create directory %s for index staging %v", c.IndexGZStaging, err)
	}

	if len(c.RepoURIs) == 0 {
		return nil
	}
	ch := make(chan result)

	for i, v := range c.RepoURIs {
		fmt.Println(v)
		c.RepoURIs[i] = v + "/Packages.gz"
		filename := strings.TrimPrefix(v, "http://")
		filename = strings.TrimPrefix(filename, "https://")
		filename = strings.ReplaceAll(filename, "/", "_")
		filename = filename + ".gz"
		c.InfoLogger.Printf("downloading: %s\tto: %s\n", c.RepoURIs[i], c.IndexGZStaging+"/"+filename)
		go c.downloadFileAsync(c.IndexGZStaging+"/"+filename, c.RepoURIs[i], ch)
	}
	for i := 0; i < len(c.RepoURIs); i++ {
		r := <-ch
		if r.err != nil {
			c.ErrorLogger.Printf(r.message)
			return err
		} else {
			c.InfoLogger.Println(r.message)
		}
	}
	c.InfoLogger.Println("finished downloading fresh package indicies")
	return nil
}

// DownloadIndicies synchronously downloads package Packages.gz from the specified repositories.
func (c *Client) DownloadIndicies() error {
	err := os.MkdirAll(c.IndexGZStaging, os.ModePerm)
	if err != nil {
		return err
	}

	if len(c.RepoURIs) == 0 {
		return nil
	}

	for i, v := range c.RepoURIs {
		fmt.Println(v)
		c.RepoURIs[i] = v + "/Packages.gz"
		filename := strings.TrimPrefix(v, "http://")
		filename = strings.TrimPrefix(filename, "https://")
		filename = strings.ReplaceAll(filename, "/", "_")
		filename = filename + ".gz"
		c.InfoLogger.Printf("downloading: %s\tto: %s\n", c.RepoURIs[i], c.IndexGZStaging+"/"+filename)
		err = c.downloadFile(c.IndexGZStaging+"/"+filename, c.RepoURIs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) downloadFileAsync(filepath string, url string, ch chan result) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		ch <- result{message: fmt.Sprintf("download failed for %s: %v", url, err), err: err}
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		ch <- result{message: fmt.Sprintf("failed to create file: %s %v\n", filepath, err), err: err}
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		ch <- result{message: fmt.Sprintf("failed to write to disk: %s %v\n", out.Name(), err), err: err}
		return
	}
	ch <- result{message: fmt.Sprintf("successfully downloaded: %s\n", url), err: nil}
}

func (c *Client) downloadFile(filepath string, url string) error {
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
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UnpackIndiciesAsync() {
	err := os.MkdirAll(c.IndexLocation, os.ModePerm)
	if err != nil {
		c.ErrorLogger.Fatalf("could not create directory %s for indexes %v", c.IndexLocation, err)
	}

	sources, err := ioutil.ReadDir(c.IndexGZStaging)
	if err != nil {
		c.ErrorLogger.Printf("could not open directory %s\n", c.IndexGZStaging)
	}
	var wg sync.WaitGroup

	for _, file := range sources {
		wg.Add(1)
		go c.readAndExtractIndexGzAsync(file, &wg)
	}

	wg.Wait()
}

func (c *Client) UnpackIndicies() error {
	err := os.MkdirAll(c.IndexLocation, os.ModePerm)
	if err != nil {
		return err
	}

	sources, err := ioutil.ReadDir(c.IndexGZStaging)
	if err != nil {
		return err
	}
	for _, file := range sources {
		err = c.readAndExtractIndexGz(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) readAndExtractIndexGzAsync(file fs.FileInfo, wg *sync.WaitGroup) {
	err := os.MkdirAll(c.IndexLocation, os.ModePerm)
	if err != nil {
		c.ErrorLogger.Fatalf("could not create directory %s for indexes %v", c.IndexLocation, err)
	}

	out, err := os.Create(c.IndexLocation + "/" + file.Name())
	if err != nil {
		c.ErrorLogger.Printf("could not create file %s: %v\n", c.IndexLocation+"/"+file.Name(), err)
	}
	fbytes, err := os.ReadFile(c.IndexGZStaging + "/" + file.Name())
	if err != nil {
		c.ErrorLogger.Printf("could not locate file %s: %v", c.IndexGZStaging+" / "+file.Name(), err)
	}
	err = gunzipWrite(out, fbytes)
	if err != nil {
		c.ErrorLogger.Printf("could not unpack %s: %v", out.Name(), err)
	}
	c.InfoLogger.Printf("finished extracting %s", out.Name())
	wg.Done()
}

func (c *Client) readAndExtractIndexGz(file fs.FileInfo) error {
	err := os.MkdirAll(c.IndexLocation, os.ModePerm)
	if err != nil {
		return err
	}
	out, err := os.Create(c.IndexLocation + "/" + file.Name())
	if err != nil {
		return err
	}
	fbytes, err := os.ReadFile(c.IndexGZStaging + "/" + file.Name())
	if err != nil {
		return err
	}
	err = gunzipWrite(out, fbytes)
	if err != nil {
		return err
	}
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
