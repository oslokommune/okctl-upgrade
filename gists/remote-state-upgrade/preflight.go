package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func preflightCheck() {
	err := ensureMainBranch()
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			log.Fatal(t(BranchInvalid))
		default:
			log.Fatal(fmt.Errorf("checking current branch: %w", err))
		}
	}

	err = ensureUpToDate()
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			log.Fatal(t(BranchOutdated))
		default:
			log.Fatal(fmt.Errorf("checking branch up to date status: %w", err))
		}
	}

	err = ensureCleanBranch()
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			log.Fatal(t(BranchDirty))
		default:
			log.Fatal(fmt.Errorf("checking branch cleanliness: %w", err))
		}
	}

	err = ensureRepositoryRootDirectory()
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			log.Fatal(t(BadCurrentWorkingDirectory))
		default:
			log.Fatal(fmt.Errorf("checking current directory: %w", err))
		}
	}
}

func ensureMainBranch() error {
	stdout := bytes.Buffer{}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	switch strings.TrimSpace(stdout.String()) {
	case "main":
		return nil
	case "master":
		return nil
	default:
		return ErrBadRequest
	}
}

func ensureUpToDate() error {
	buf := bytes.Buffer{}

	cmd := exec.Command("git", "fetch", "--dry-run")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	if buf.Len() != 0 {
		return ErrBadRequest
	}

	buf = bytes.Buffer{}

	cmd = exec.Command("git", "push", "--dry-run")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	if strings.TrimSpace(buf.String()) != "Everything up-to-date" {
		return ErrBadRequest
	}

	return nil
}

func ensureCleanBranch() error {
	stdout := bytes.Buffer{}

	cmd := exec.Command("git", "status", "-s")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	switch stdout.Len() {
	case 0:
		return nil
	default:
		return ErrBadRequest
	}
}

func ensureRepositoryRootDirectory() error {
	stdout := bytes.Buffer{}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %w", err)
	}

	if currentDirectory != strings.TrimSpace(stdout.String()) {
		return ErrBadRequest
	}

	return nil
}
