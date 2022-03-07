package kubectl

import (
	"errors"

	"github.com/oslokommune/okctl-upgrade/upgrades/0.0.93-increase-argocd-memlimit/pkg/lib/logger"
	"k8s.io/client-go/kubernetes"
)

const (
	timeoutSeconds = 300
)

type Kubectl struct {
	log       logger.Logger
	clientSet *kubernetes.Clientset
}

var ErrNotFound = errors.New("nothing to do")
