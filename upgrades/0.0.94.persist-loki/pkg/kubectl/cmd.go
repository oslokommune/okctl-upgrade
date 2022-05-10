package kubectl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os/exec"
	"path"
	"runtime"

	"github.com/Masterminds/semver"
	"github.com/spf13/afero"
)

func acquireBinaryPath(Fs *afero.Afero, homeDirFn func() (string, error)) (string, error) {
	homeDir, err := homeDirFn()
	if err != nil {
		return "", fmt.Errorf("acquiring home directory: %w", err)
	}

	binaryDir := path.Join(homeDir, defaultOkctlConfigDirName, defaultOkctlBinariesDirName, defaultBinaryName)

	exists, err := Fs.DirExists(binaryDir)
	if err != nil {
		return "", fmt.Errorf("checking binary directory existence: %w", err)
	}

	if !exists {
		return "", errors.New("missing binary directory")
	}

	versions, err := gatherVersions(Fs, binaryDir)
	if err != nil {
		return "", fmt.Errorf("gathering versions: %w", err)
	}

	var latest semver.Version
	for version := range versions {
		if version.GreaterThan(&latest) {
			latest = version
		}
	}

	return path.Join(binaryDir, latest.String(), runtime.GOOS, runtime.GOARCH, defaultBinaryName), nil
}

func runCommand(binaryPath string, args ...string) (io.Reader, error) {
	cmd := exec.Command(binaryPath, args...) //nolint:gosec

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	return &stdout, nil
}

func gatherVersions(Fs *afero.Afero, baseDir string) (map[semver.Version]interface{}, error) {
	versions := make(map[semver.Version]interface{})

	err := Fs.Walk(baseDir, func(currentPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isDirectory, err := Fs.IsDir(currentPath)
		if err != nil {
			return fmt.Errorf("checking path type: %w", err)
		}

		if !isDirectory {
			return nil
		}

		version, err := semver.NewVersion(path.Base(currentPath))
		if err == nil {
			versions[*version] = true
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return versions, nil
}
