package aptlib

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// SourcesListFromPath parses sources in a given io.Reader.
func (c *Client) ParseSourcesList(r io.Reader) error {
	out := []*SourceEntry{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			source, err := SourceFromString(line)
			if err != nil {
				return err
			}
			out = append(out, source)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	c.SourceEntries = append(c.SourceEntries, out...)
	return nil
}

// AutoDetectSources checks default locations for sources.list and sets the client.
func (c *Client) AutoDetectSources() error {
	err := c.ParseSourcesFromFile("/etc/apt/sources.list")
	if err != nil {
		return err
	}

	sources, err := ioutil.ReadDir("/etc/apt/sources.list.d/")
	if err != nil {
		return err
	}

	for _, file := range sources {
		if strings.Contains(file.Name(), ".list") {
			err := c.ParseSourcesFromFile("/etc/apt/sources.list.d/" + file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SourcesFromFile parses a sources.list at a given path.
func (c *Client) ParseSourcesFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	err = c.ParseSourcesList(f)
	if err != nil {
		return err
	}
	return nil
}
