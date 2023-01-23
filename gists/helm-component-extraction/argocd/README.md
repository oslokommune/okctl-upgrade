# ArgoCD

## Chart details

| Chart                                                               | Version    | App version |
| ---------------------------------------------------------------     | --------   | ----------- |
| [source](https://artifacthub.io/packages/helm/argo/argo-cd)         | `v5.17.1`  | `v2.5.6`    |

## Preparation

### Log into AWS

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```shell
    export AWS_PROFILE=your-profile-name
    ```
2. Log in to the AWS account with `aws sso login`
3. Verify that you are logged in by running `aws s3 ls`

### Log into okctl

1. Run `okctl venv -c <path to relevant cluster.yaml>` to log into the okctl cluster.
2. Verify that you are logged in by running `kubectl get pods -A`

### Backup existing applications

1. Run `make applications-snapshot` to backup the existing applications.
2. Verify that there is an `applications-snapshot.yaml` file in the current directory and that it contains the relevant YAML

## Installation or upgrade

Depending on whether you are installing or upgrading, go to the relevant runbook below.

### Installation

Choose [INSTALL.md](./INSTALL.md) if you can say yes to all of the following:
- `integrations.argoCD` is missing or set to `true` in your cluster.yaml
- You have not installed ArgoCD with Makefile or Helm before for this cluster

### Upgrading

Choose [UPGRADE.md](./UPGRADE.md) if you can say yes to all of the following:
- You've uninstalled the Okctl version 
- You've previously installed the Helm version of ArgoCD using Helm or a Makefile

