package aptlib

// RepoType represents the types of repos in apt.
type RepoType string

// There are two valid apt repo types:
// binary (represented as "deb")
// source (represented as "deb-src").
const (
	BinaryRepo RepoType = "deb"
	SrcRepo    RepoType = "deb-src"
)

// IsValid returns true when the given RepoType is valid.
func (rt RepoType) IsValid() bool {
	switch rt {
	case BinaryRepo, SrcRepo:
		return true
	}
	return false
}
