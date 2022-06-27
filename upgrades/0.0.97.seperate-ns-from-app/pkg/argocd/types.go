package argocd

import "io"

type debugLogger interface {
	Debug(args ...interface{})
}

type applier interface {
	Apply(io.Reader) error
}
