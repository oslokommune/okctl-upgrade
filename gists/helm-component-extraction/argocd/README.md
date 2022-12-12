# ArgoCD

## Chart details

| Chart                                                               | Version    | App version |
| ---------------------------------------------------------------     | --------   | ----------- |
| [source](https://artifacthub.io/packages/helm/argo/argo-cd/3.26.12) | `v3.26.12` | `v2.1.7`    |

## Prerequisites

### Log into AWS

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
		```
2. Log in to the AWS account with `aws sso login`

### Log into okctl

1. Run `okctl venv -c <path to relevant cluster.yaml>` to log into the okctl cluster.

## Backup existing applications

1. Run `make applications-snapshot` to backup the existing applications.
2. Verify that there is an `applications.yaml` file in the current directory and that it contains the relevant YAML

## Uninstall okctl provisioned ArgoCD

1. Edit the relevant `cluster.yaml` file and set `integrations.argocd` to `false`. If all integrations are commented
   out showing default values you need to uncomment the entire `integrations` section to avoid nil references.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.

## Installation

The installation consists of three parts:

1. Setup a Cognito user pool client
2. Setup a TLS certificate
3. Install the Helm chart

### Initial

1. Place this directory (argocd) in `infrastructure/<cluster name>/helm-components/`
2. Edit the Makefile and configure the user input variables.

### Setup ArgoCD Cognito client

1. Run `make configure-cognito-client` to generate a Cognito client CloudFormation template.
2. Run `make install-cognito-client` to install the Cognito client CloudFormation stack.
3. Run `make install-cognito-parameters` to install the Cognito client parameters in AWS SSM Parameter Store.

### Setup ArgoCD Certificate

1. Run `make configure-certificate` to generate a Certificate CloudFormation template.
2. Run `make install-certificate` to install the Certificate CloudFormation stack.

### Install ArgoCD

1. Run `make configure-helm-chart` to generate necessary files to accompany the Helm chart.
2. Run `make install-helm-chart` to install the Helm chart.

### Setup deploy key

1. Run `make configure-deploy-key` to generate a deploy key.
2. Run `make install-deploy-key` to install the key that ArgoCD uses to read private repositories

### Reapply applications

1. Run `make install-applications-snapshot` to reapply the applications.

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `application.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
		values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.

Github fine grained PAT must have administration read/write 
