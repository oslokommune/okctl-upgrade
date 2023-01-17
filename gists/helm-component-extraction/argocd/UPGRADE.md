# Upgrade ArgoCD

## From 2.1.7 to 2.5.6

### Back up existing applications

```shell
make applications-snapshot
```

### Upgrade ArgoCD

```shell
# Prepare upgraded values
make configure-helm-chart
# Install upgraded chart with values
make install-helm-chart
```

### Restore applications

```shell
make install-applications-snapshot
```

## Something went wrong

```shell
# Ensure ArgoCD doesn't uninstall applications
make remove-finalizers
# Uninstall existing helm chart
make uninstall-helm-chart
# Install upgraded helm chart
make install-helm-chart
```

## FAQ
- **How do I check the current version of ArgoCD?** Run `kubectl describe --namespace=argocd argocd deployment argocd-server`
    and look for the `Image:` field.
