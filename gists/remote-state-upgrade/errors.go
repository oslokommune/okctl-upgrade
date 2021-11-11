package main

import "errors"

var (
	ErrNotFound         = errors.New("okctl/NotFound")
	ErrAlreadyExists    = errors.New("okctl/AlreadyExists")
	ErrNotAuthenticated = errors.New("okctl/NotAuthenticated")
	ErrBadRequest       = errors.New("okctl/BadRequest")
)

const (
	BadCurrentWorkingDirectory = iota
	BranchDirty
	BranchInvalid
	BranchOutdated
	BucketAlreadyExists
	ClusterManifestNotFound
	NotAuthenticated
)

func t(key int) string {
	switch key {
	case ClusterManifestNotFound:
		return "Cluster manifest not found at provided path, are you sure it is correct?"
	case BucketAlreadyExists:
		return "The bucket already exists. To not overwrite any important data, please contact #kjøremiljø"
	case NotAuthenticated:
		return "There seems to be a problem with the authentication. Please run 'okctl venv -c <cluster manifest>' again"
	case BranchInvalid:
		return "Invalid git branch: please switch to your IAC repo's main branch"
	case BranchOutdated:
		return "Outdated git branch: please ensure your branch is synced with the remote with 'git pull' and 'git push' before upgrading"
	case BranchDirty:
		return "Dirty git branch: please commit your local changes"
	case BadCurrentWorkingDirectory:
		return "Please run this upgrade in the IAC repo's root folder"
	default:
		return "missing translation"
	}
}
