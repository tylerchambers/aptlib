package aptlib

import (
	"log"
	"os"
)

// Client acts as an apt client.
type Client struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	// Arch contains a list of desired architectures.
	Arch string
	// Absolute paths to sources.list to be considered.
	SourceLists []string
	// The actual parsed entries
	SourceEntries []*SourceEntry
	// Where to store the package index files
	IndexLocation  string
	IndexGZStaging string
	RepoURIs       []string
}

// NewClient instantiates a new apt client.
func NewClient(infoLogger, warningLogger, errorLogger *log.Logger, sourceLists []string, arch, indexLocation, indexGZStaging string) *Client {
	c := &Client{
		InfoLogger:     infoLogger,
		WarningLogger:  warningLogger,
		ErrorLogger:    errorLogger,
		Arch:           arch,
		SourceLists:    sourceLists,
		IndexLocation:  indexLocation,
		IndexGZStaging: indexGZStaging,
	}
	return c
}

// Init initializes a client with sane defaults.
func (c *Client) Init() error {
	if c.InfoLogger == nil {
		c.InfoLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	}
	if c.WarningLogger == nil {
		c.WarningLogger = log.New(os.Stdout, "WARN: ", log.LstdFlags)
	}
	if c.ErrorLogger == nil {
		c.ErrorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
	}
	if c.Arch == "" {
		err := c.AutoDetectHostArch()
		if err != nil {
			return err
		}
	}
	if c.SourceLists == nil {
		err := c.AutoDetectSources()
		if err != nil {
			return err
		}
	}
	if c.IndexLocation == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.IndexLocation = cwd + "/index"
	}
	if c.IndexGZStaging == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.IndexGZStaging = cwd + "/index_staging"
	}
	return nil
}
