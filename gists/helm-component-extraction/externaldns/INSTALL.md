# Installing External DNS

## Uninstall okctl provisioned ExternalDNS

1. Edit the relevant `cluster.yaml` file and set `integrations.externalDNS` to `false`. If all integrations are commented
   out showing default values you need to uncomment the entire `integrations` section to avoid nil references.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.
3. Run `kubectl --namespace=kube-system delete deployments.apps external-dns`. This will delete the deployment and the
   associated pods.

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

- **How do I use a different version?** To use a different version, change the `spec.template.spec.containers[].image` field
    in `deployment.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
    values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.

