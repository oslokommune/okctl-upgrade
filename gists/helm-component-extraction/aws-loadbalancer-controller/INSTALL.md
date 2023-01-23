# Installing AWS LoadBalancer Controller

## Uninstall okctl provisioned AWS LoadBalancer Controller

1. Edit the relevant `cluster.yaml` file and set `integrations.awsLoadBalancerController` to `false`.
2. Run `okctl apply cluster -f cluster.yaml` to apply the changes.

## Installation

### Prepare Cloudformation stacks

See [cloudformation/README.md](cloudformation/README.md) for details.

### Install the AWS LoadBalancer Controller

1. Edit the Makefile and set the variable(s) listed in the user input section
2. Run `make configure` to generate the values.yaml file
3. Run `make install` to install the component

## FAQ

- **How do I use a different chart version?** To use a different chart version, change the `spec.source.targetRevision` field
    in `application.yaml` to the desired version. Then run `make install` to install the new version. N.B.: The required
    values are not guaranteed to be the same between versions, so pay attention to the changelog in the chart link above.
