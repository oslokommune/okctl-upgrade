# ExternalDNS

## Component details

| Component                                                                     | Version  |
| ----------------------------------------------------------------------------- | -------- |
| [source](https://github.com/kubernetes-sigs/external-dns/releases/tag/v0.7.3) | `0.7.3`  |

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

### Uninstall okctl provisioned ExternalDNS

1. Edit the relevant `cluster.yaml` file and set `integrations.externalDNS` to `false`. If all integrations are commented
   out showing default values you need to uncomment the entire `integrations` section to avoid nil references.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.

## Installation

The installation consists of the following parts:

1. Configure neccessary resources
2. Apply the configured resources

### Initial

1. Place this directory (externaldns) in `infrastructure/<cluster name>/helm-components/`
2. Edit the Makefile and configure the user input variables.

### Configure neccessary resources

To configure the neccessary resources, run the following command:

```shell
make configure
```

### Apply the configured resources

To apply the configured resources, run the following command:

```shell
make install
```

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `applications-snapshot.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
		values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.

