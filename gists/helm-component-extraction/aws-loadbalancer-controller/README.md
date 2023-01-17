# AWS LoadBalancer Controller

## Chart details

| Chart                                                                                 | Version  | App version |
| ----------------------------------------------------------------------------------    | -------- | ----------- |
| [source](https://artifacthub.io/packages/helm/aws/aws-load-balancer-controller/1.1.3) | `v1.1.3` | `v2.1.1`    |

## Prerequisites

### Log into AWS

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
    ```
2. Log in to the AWS account with `aws sso login`

### Log into okctl
 
1. Run `okctl venv -c <path to relevant cluster.yaml>` to log into the okctl cluster.

## Installation or upgrade

Depending on whether you are installing or upgrading, go to the relevant runbook below.

### Installation

Choose [INSTALL.md](./INSTALL.md) if you can say yes to all of the following:
- `integrations.awsLoadBalancerController` is missing or set to `false` in your cluster.yaml
- You have not installed the AWS LoadBalancer Controller with the Makefile or Helm before

### Upgrading

Choose [UPGRADE.md](./UPGRADE.md) if you can say yes to all of the following:
- You've uninstalled the Okctl version 
- You've installed the Helm version of the AWS LoadBalancer Controller

