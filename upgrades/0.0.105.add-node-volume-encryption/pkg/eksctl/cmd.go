package eksctl

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func runCommand(args ...string) (io.Reader, error) {
	cmd := exec.Command("eksctl", args...)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	return &stdout, nil
}
