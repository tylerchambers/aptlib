package aptlib

import (
	"log"
)

// AptClient stores client information.
type AptClient struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	// Arch contains a list of desired architectures.
	Arch []string
	// Absolute paths to sources.list to be considered.
	SourceLists []string
	// The actual parsed entries
	SourceEntries []*SourceEntry
	// Where to store the pacakge index files
	IndexLocation string
}
