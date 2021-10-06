package aptlib

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetHostArch gets the host architecture from dpkg.
func GetHostArch() (string, error) {
	cmd := exec.Command("dpkg-architecture", "-q", "DEB_HOST_ARCH")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(out.String(), "\n"), nil
}
