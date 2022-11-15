# Loki

## Chart details

| Chart                                                                                | Version  | App version |
| ----------------------------------------------------------------------------------   | -------- | ----------- |
| [source](https://artifacthub.io/packages/helm/grafana/loki/2.3.0)                    | `v2.3.0` | `v2.1.0`   |

## Prerequisites

### Log into AWS

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
		```
2. Log in to the AWS account with `aws sso login`

### Log into okctl
 
1. Run `okctl venv -c <path to relevant cluster.yaml>` to log into the okctl cluster.

## Uninstall okctl provisioned component

1. Run `helm --namespace monitoring uninstall loki` to uninstall the component. Your newer logs will still be picked up by
    promtail, and existing logs is stored in an S3 bucket.

## Installation

### Install the component

1. Edit the Makefile and set the variable(s) listed in the user input section
2. Run `make values` to generate the values.yaml file
3. Run `make install` to install the component

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `application.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
		values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.
- **There is no instructions regarding the CloudFormation directory, what gives?** The CloudFormation directory is only 
    relevant when installing the component from scratch. If the component is already installed through okctl, the resources
		are already created and the CloudFormation directory is redundant.
