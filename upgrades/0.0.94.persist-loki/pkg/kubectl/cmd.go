package kubectl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

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

func runCommand(fs *afero.Afero, kubeconfigPath string, args ...string) (io.Reader, error) {
	cmd := exec.Command(defaultBinaryName, args...) //nolint:gosec

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	env, err := generateEnv(fs, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("generating env: %w", err)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = env

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	return &stdout, nil
}

func generateEnv(fs *afero.Afero, kubeconfigPath string) ([]string, error) {
	envMap := make(map[string]string)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("acquiring home directory: %w", err)
	}

	binariesDir := path.Join(homeDir, defaultOkctlConfigDirName, defaultOkctlBinariesDirName)
	toolDirectories := make([]string, 0)

	for _, tool := range []string{defaultBinaryName, defaultAWSIAMAuthenticatorName} {
		toolBaseDir := path.Join(binariesDir, tool)

		versions, err := gatherVersions(fs, toolBaseDir)
		if err != nil {
			return nil, fmt.Errorf("gathering versions: %w", err)
		}

		latestVersion := getLatestVersion(versions)

		toolDirectories = append(toolDirectories, path.Join(
			toolBaseDir,
			latestVersion.String(),
			runtime.GOOS,
			defaultArch,
			tool,
		))
	}

	envMap["HOME"] = homeDir
	envMap["KUBECONFIG"] = kubeconfigPath
	envMap["PATH"] = strings.Join(toolDirectories, ":")

	return envAsArray(envMap), nil
}

// envAsArray converts a map to a string array of KEY=VALUE pairs
func envAsArray(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return result
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

func getLatestVersion(versions map[semver.Version]interface{}) semver.Version {
	var latest semver.Version

	for version := range versions {
		if version.GreaterThan(&latest) {
			latest = version
		}
	}

	return latest
}
