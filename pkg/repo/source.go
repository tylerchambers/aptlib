package repo

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// SourceEntry represents an entry in a sources.list file.
type SourceEntry struct {
	RepoType     RepoType
	Options      map[string]string
	URI          string
	Distribution string
	Components   []string
}

// NewSource returns a pointer to a source type.
func NewSource(rt RepoType, options map[string]string, URI string, dist string, components []string) *SourceEntry {
	return &SourceEntry{rt, options, URI, dist, components}
}

// SourceFromString validates a sources.list entry and if valid returns a source type.
// TODO (tylerchambers): "distribution may also contain a variable, $(ARCH) which expands to the Debian architecture (i386, m68k, powerpc, ...) used on the system"
// I have never seen this in the wild, but we should probably support it eventually.
func SourceFromString(line string) (*SourceEntry, error) {
	entry := strings.Fields(line)
	if len(entry) < 4 {
		return nil, fmt.Errorf("not enough fields in source entry: %s", line)
	}
	rt := RepoType(entry[0])
	if !rt.IsValid() {
		return nil, fmt.Errorf("invalid repo type: %s", entry[0])
	}
	// Setup our map for potential options.
	// If this matches our current line specifies options.
	optionsRegex := regexp.MustCompile(`^deb?(-src)?\s(\[.*\])\s.*`)
	if optionsRegex.MatchString(line) {
		options := make(map[string]string)
		runes := []rune(line)
		optionsList := string(runes[strings.Index(line, "[")+1 : strings.LastIndex(line, "]")])
		if len(optionsList) < 3 {
			return nil, fmt.Errorf("options list cannot be empty: %s", line)
		}
		if strings.Contains(optionsList, "[") || strings.Contains(optionsList, "]") {
			return nil, fmt.Errorf("malformed options: %s", line)
		}
		for _, pair := range strings.Fields(optionsList) {
			if !strings.Contains(pair, "=") {
				return nil, fmt.Errorf("malformed option key value pair: %s", pair)
			}
			kv := strings.Split(pair, "=")
			options[kv[0]] = kv[1]
		}
		// We're done parsing the options. Pickup where we left off.
		rest := string(runes[strings.LastIndex(line, "]")+1:])
		restFields := strings.Fields(rest)
		URI, err := url.ParseRequestURI(restFields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid URI: %s", restFields[0])
		}
		dist := restFields[1]
		components := restFields[2:]
		return NewSource(rt, options, URI.String(), dist, components), nil

	}
	URI, err := url.ParseRequestURI(entry[1])
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %s", entry[1])
	}
	dist := entry[2]
	components := entry[3:]
	return NewSource(rt, nil, URI.String(), dist, components), nil
}
