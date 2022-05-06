package v1alpha1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ClusterKind is a string value that represents the resource type
	ClusterKind = "Cluster"
	// ClusterAPIVersion defines the versioned schema of this representation
	ClusterAPIVersion = "okctl.io/v1alpha1"
)

// Cluster is a unique Kubernetes cluster with a set of integrations that
// can be enabled or disabled.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`

	// Metadata uniquely identifies a cluster.
	Metadata ClusterMeta `json:"metadata"`

	// Github defines what organisation, repository, etc. that
	// this cluster will integrate with.
	Github ClusterGithub `json:"github"`

	// ClusterRootDomain defines the main primary zone to associate with this
	// cluster. This will be the zone that we will use to create subdomains
	// for auth, ArgoCD, etc.
	ClusterRootDomain string `json:"clusterRootDomain"`

	// VPC defines how we configure the VPC for the cluster
	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// Integrations defines what cluster integrations we deploy to the
	// cluster
	// +optional
	Integrations *ClusterIntegrations `json:"integrations,omitempty"`

	// DNSZones is an optional list of DNS zones managed or associated with
	// this cluster.
	// +optional
	DNSZones []ClusterDNSZone `json:"dnsZones,omitempty"`

	// Users is an optional list of email addresses
	// +optional
	Users []ClusterUser `json:"users,omitempty"`

	// Databases is an optional list of databases
	// +optional
	Databases *ClusterDatabases `json:"databases,omitempty"`

	// Experimental is an optional section for testing
	// +optional
	Experimental *ClusterExperimental `json:"experimental,omitempty"`
}

// ClusterMeta describes a unique cluster
type ClusterMeta struct {
	// Name is a descriptive value given to the cluster, e.g., the name
	// of the team, product, project, etc.
	Name string `json:"name"`

	// Region specifies the AWS region the cluster should be created in
	// https://aws.amazon.com/about-aws/global-infrastructure/regions_az/
	Region string `json:"region"`

	// AccountID specifies the AWS Account ID
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html
	AccountID string `json:"accountID"`
}

// String returns a unique identifier for a cluster
// Not sure about this..
func (receiver *ClusterMeta) String() string {
	return fmt.Sprintf("%s.%s.okctl.io/%s", receiver.Name, receiver.Region, receiver.AccountID)
}

// ClusterVPC is a definition of the VPC we create for the EKS cluster
type ClusterVPC struct {
	// CIDR is the IP-address range to associate with the VPC
	// https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing.
	// The VPC CIDR must be compatible with EKS: https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html
	// +optional
	CIDR string `json:"cidr,omitempty"`

	// HighAvailability means we create redundancy in the network setup. If set to
	// true we will create a NAT gateway per public subnet, instead of routing
	// all traffic through one.
	// +optional
	HighAvailability bool `json:"highAvailability,omitempty"`
}

// ClusterDNSZone is analogous to a DNS Zone file (https://en.wikipedia.org/wiki/Zone_file).
// A DNS Zone represents a subset, in form of a single parent domain, of the hierarchical
// domain name structure. In AWS, we map this data to a Route53 HostedZone.
type ClusterDNSZone struct {
	// ParentDomain is the root domain for all DNS records of this
	// DNS zone delegation, e.g., `{team-name}.oslo.systems`
	ParentDomain string `json:"parentDomain"`

	// ReuseExisting determines if we should look for an existing DNS zone
	// or create a new one. If set to true, we will not attempt to create a
	// new DNS zone.
	ReuseExisting bool `json:"managedZone"`
}

// ClusterGithub identifies a repository and path on github.com where
// we can set up an integration with Argo CD, among other things.
type ClusterGithub struct {
	// Organisation name on github.com, e.g., "oslokommune"
	Organisation string `json:"organisation"`

	// Repository name on github.com, e.g., "okctl". The repository
	// you specify here must be owned by the organisation specified above.
	Repository string `json:"repository"`

	// OutputPath is a path from the root of the org/repository where
	// we can store generated output files
	OutputPath string `json:"outputPath"`
}

// Path returns the Github repository URL path
func (c ClusterGithub) Path() string {
	return fmt.Sprintf("%s/%s", c.Organisation, c.Repository)
}

// URL returns the Github IAC repository URL
func (c ClusterGithub) URL() string {
	return fmt.Sprintf("git@github.com:%s", c.Path())
}

// ClusterIntegrations ...
type ClusterIntegrations struct {
	// AWSLoadBalancerController if set to true will install the AWS load balancer controller
	// +optional
	AWSLoadBalancerController bool `json:"awsLoadBalancerController"`

	// ExternalDNS if set to true will install the external-dns controller into the cluster
	// +optional
	ExternalDNS bool `json:"externalDNS,omitempty"`

	// ExternalSecrets if set to true will install the external-secrets controller into the cluster
	// +optional
	ExternalSecrets bool `json:"externalSecrets,omitempty"`

	// Autoscaler if set to true will install the cluster autoscaler into the cluster
	// +optional
	Autoscaler bool `json:"autoscaler,omitempty"`

	// KubePromStack if set to true will install the kubernetes-prometheus-stack into the cluster
	// We should probably give this a better name, something more related to monitoring, but
	// we can think about that down the road.
	// +optional
	KubePromStack bool `json:"kubePromStack,omitempty"`

	// Loki if set to true will install the Loki log collector and data source for grafana into
	// the cluster.
	Loki bool `json:"loki,omitempty"`

	// Promtail if set to true will install the Promtail log scraper
	Promtail bool `json:"promtail,omitempty"`

	// Tempo if set to true will install tempo for trace ingestion
	Tempo bool `json:"tempo,omitempty"`

	// Blockstorage if set to true will install the EBS CSI block storage driver into the
	// cluster, which makes it possible to create PersistentVolumeClaims in AWS
	// +optional
	Blockstorage bool `json:"blockstorage,omitempty"`

	// Cognito if set to true will install the Cognito user pool into the cluster.
	// Might want to make this one more fine-grained, so that the teams can more easily
	// give access to their admin APIs or whatever. Might not be required for now.
	// +optional
	Cognito bool `json:"cognito,omitempty"`

	// ArgoCD if set to true will install the ArgoCD deployment setup into the cluster. This
	// integration requires ALBIngressController, ExternalDNS and Cognito.
	// +optional
	ArgoCD bool `json:"argoCD,omitempty"`
}

// ClusterUser represents the identity of a user
// that should have access to the cluster
type ClusterUser struct {
	// Email is the valid email address of the user
	Email string `json:"email"`
}

// ClusterDatabases contains the declaration of
// different types of databases
type ClusterDatabases struct {
	// Postgres contains the declared list of postgres databases
	// +optional
	Postgres []ClusterDatabasesPostgres `json:"postgres"`
}

// ClusterDatabasesPostgres contains the declaration of a postgres database
type ClusterDatabasesPostgres struct {
	// Name we should give to the database
	Name string `json:"name"`

	// User is the name we give to the admin user,
	// you can not set this to `admin` as that is a reserved
	// word
	User string `json:"user"`

	// Namespace determines where we will write the
	// Kubernetes ConfigMap and Secret; for easily
	// accessing the database
	Namespace string `json:"namespace"`
}

// ClusterExperimental contains experimental fields
type ClusterExperimental struct {
	// AutomatizeZoneDelegation will automatically merge the delegation
	// pull requests when set to true
	// +optional
	AutomatizeZoneDelegation bool `json:"automatizeZoneDelegation"`
}

// Validate the content of cluster experimental
func (e ClusterExperimental) Validate() error {
	return nil
}

// ClusterTypeMeta returns an initialised TypeMeta object
// for a Cluster
func ClusterTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterKind,
		APIVersion: ClusterAPIVersion,
	}
}

// NewCluster returns a Cluster with sensible defaults
func NewCluster() Cluster {
	return Cluster{
		TypeMeta: ClusterTypeMeta(),
		Metadata: ClusterMeta{
			Name:      "",
			Region:    "eu-west-1",
			AccountID: "",
		},
		Github: ClusterGithub{
			Organisation: "oslokommune",
			Repository:   "",
			OutputPath:   "infrastructure",
		},
		ClusterRootDomain: "",
		VPC: &ClusterVPC{
			CIDR:             "",
			HighAvailability: true,
		},
		Integrations: &ClusterIntegrations{
			AWSLoadBalancerController: true,
			ExternalDNS:               true,
			ExternalSecrets:           true,
			Autoscaler:                true,
			KubePromStack:             true,
			Loki:                      true,
			Promtail:                  true,
			Tempo:                     true,
			Blockstorage:              true,
			Cognito:                   true,
			ArgoCD:                    true,
		},
		Databases: &ClusterDatabases{},
		Experimental: &ClusterExperimental{
			AutomatizeZoneDelegation: false,
		},
	}
}
