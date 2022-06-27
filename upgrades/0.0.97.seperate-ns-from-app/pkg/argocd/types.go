package argocd

type debugLogger interface {
	Debug(args ...interface{})
}
