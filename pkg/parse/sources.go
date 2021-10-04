package parse

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tylerchambers/goapt/pkg/repo"
)

// SourcesListFromPath parses sources in a given io.Reader.
func SourcesList(r io.Reader) ([]*repo.SourceEntry, error) {
	out := []*repo.SourceEntry{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			source, err := repo.SourceFromString(line)
			if err != nil {
				return nil, err
			}
			out = append(out, source)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// AllSources checks default locations for sources.list and parses everything it finds.
func AllSources() ([]*repo.SourceEntry, error) {
	out, err := SourcesFromFile("/etc/apt/sources.list")
	if err != nil {
		return nil, err
	}

	sources, err := ioutil.ReadDir("/etc/apt/sources.list.d/")
	if err != nil {
		return nil, err
	}

	for _, file := range sources {
		if strings.Contains(file.Name(), ".list") {
			s, err := SourcesFromFile("/etc/apt/sources.list.d/" + file.Name())
			if err != nil {
				return nil, err
			}
			out = append(out, s...)
		}
	}

	return out, nil
}

// SourcesFromFile parses a sources.list at a given path.
func SourcesFromFile(path string) ([]*repo.SourceEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	out, err := SourcesList(f)
	if err != nil {
		return nil, err
	}
	return out, nil
}
