# Promtail

## Chart details

| Chart                                                           | Version  | App version |
| --------------------------------------------------------------- | -------- | ----------- |
| [source](https://artifacthub.io/packages/helm/grafana/promtail) | `v6.6.1` | `v2.6.1`    |

## Prerequisites

### Log into AWS

1. Export the `AWS_PROFILE` variable with the relevant profile name.

```bash
export AWS_PROFILE=your-profile-name
```

2. Log in to the AWS account with `aws sso login`

### Log into okctl

1. Run `okctl venv -c <path to relevant cluster.yaml>` to log into the okctl cluster.

## Uninstall okctl provisioned Promtail

1. Edit the relevant `cluster.yaml` file and set `integrations.promtail` to `false`. If all integrations are commented
   out showing default values you need to uncomment the entire `integrations` section to avoid nil references.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.

## Installation

### Install Promtail

1. Place this directory (promtail) in `infrastructure/<cluster name>/helm-components/`
2. **Optional**; edit the `values.yaml` to suit your setup.
3. Run `make install` to install Promtail.

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `application.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
		values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.
