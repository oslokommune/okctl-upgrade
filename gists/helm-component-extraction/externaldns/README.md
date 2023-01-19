# ExternalDNS

## Component details

| Component                                                                     | Version  |
| ----------------------------------------------------------------------------- | -------- |
| [source](https://github.com/kubernetes-sigs/external-dns)                     | `0.13.1` |

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

## Installation or upgrade

Depending on whether you are installing or upgrading, go to the relevant runbook below.

### Installation

Choose [INSTALL.md](./INSTALL.md) if you can say yes to all of the following:
- `integrations.externalDNS` is missing or set to `true` in your cluster.yaml
- You have not installed the ExternalDNS with the Makefile or Helm before

### Upgrading

Choose [UPGRADE.md](./UPGRADE.md) if you can say yes to all of the following:
- You've uninstalled the Okctl version by setting `integrations.externalDNS` in your cluster.yaml to false, then running
    `okctl apply cluster`
- You've previously installed the Helm version of the AWS LoadBalancer Controller

