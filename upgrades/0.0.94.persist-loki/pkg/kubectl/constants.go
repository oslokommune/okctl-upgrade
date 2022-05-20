package kubectl

const (
	AvailabilityZoneLabelKey       = "failure-domain.beta.kubernetes.io/zone"
	defaultOkctlConfigDirName      = ".okctl"
	defaultOkctlBinariesDirName    = "binaries"
	defaultBinaryName              = "kubectl"
	defaultMonitoringNamespace     = "monitoring"
	defaultArch                    = "amd64"
	defaultAWSIAMAuthenticatorName = "aws-iam-authenticator"
	defaultLokiPodName             = "loki-0"
	defaultLokiConfigSecretKey     = "loki.yaml"

	// Kubernetes kinds
	persistentVolumeClaimResourceKind = "persistentvolumeclaim"
	persistentVolumeResourceKind      = "persistentvolume"
	secretResourceKind                = "secret"
	podResourceKind                   = "pod"
	statefulSetResourceKind           = "statefulset"
	serviceAccountResourceKind        = "serviceaccount"
)
