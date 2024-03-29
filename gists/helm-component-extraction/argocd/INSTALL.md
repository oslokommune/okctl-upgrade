# Installing ArgoCD

## Prerequisites

### Setup a [Github Personal Access Token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

1. Generate a fine grained PAT with `administration read/write permissions` for your IAC repository
2. Export the token as an environment variable
    ```shell
    export GITHUB_TOKEN=your-token
    ```

### Uninstall okctl provisioned ArgoCD

1. Edit the relevant `cluster.yaml` file and set `integrations.argocd` to `false`. If all integrations are commented
   out showing default values you need to uncomment the entire `integrations` section to avoid nil references.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.

## Installation

The installation consists of the following parts:

1. Setup a Cognito user pool client
2. Setup an SSL certificate
3. Install the Helm chart
4. Setup a [deploy key](https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys) for the repository

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
2. Run `make install-deploy-key-parameters` to store deploy key in SSM.
3. Run `make install-deploy-key` to install the key that ArgoCD uses to read private repositories

### Reapply applications

1. Run `make install-applications-snapshot` to reapply the applications.

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `applications-snapshot.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
		values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.

