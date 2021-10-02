package parse

import (
	"bufio"
	"os"
	"strings"

	"github.com/tylerchambers/goapt/pkg/repo"
)

// SourcesListFromPath parses the sources.list specified at a given path.
func SourcesList(path string) ([]*repo.SourceEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	out := []*repo.SourceEntry{}
	scanner := bufio.NewScanner(f)

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
