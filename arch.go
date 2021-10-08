package aptlib

import (
	"bytes"
	"os/exec"
	"strings"
)

// AutoDetectHostArch detects the host architecture using dpkg and sets the client to match.
func (c *Client) AutoDetectHostArch() error {
	cmd := exec.Command("dpkg-architecture", "-q", "DEB_HOST_ARCH")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return err
	}

	c.Arch = strings.TrimSuffix(out.String(), "\n")
	return nil
}
