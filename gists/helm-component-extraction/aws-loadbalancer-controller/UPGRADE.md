# Upgrading AWS LoadBalancer Controller

## Upgrade

Before starting, ensure you've configured the user input variables in the `Makefile`.

1. Run `make configure` to configure the AWS LoadBalancer Controller values
2. Run `make install` to update the AWS LoadBalancer Controller ArgoCD application manifest with a new version and new
    values
