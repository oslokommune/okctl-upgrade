package eksctl

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func runCommand(stdin io.Reader, args ...string) (io.Reader, error) {
	cmd := exec.Command("eksctl", args...)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = stdin

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	return &stdout, nil
}
