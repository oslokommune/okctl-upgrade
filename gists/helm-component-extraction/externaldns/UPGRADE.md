# Upgrading ExternalDNS

## Upgrade

Before starting, ensure you've configured the user input variables in the `Makefile`.

1. Delete previously generated `deployment.yaml`, `clusterrole.yaml` and `clusterrolebinding.yaml` files
2. Run `make configure` to recreate the updated resources deleted in step 1
3. Run `make install` to update the resources

## How to verify a successful upgrade

1. Run `kubectl get deployment -n kube-system external-dns -o yaml` to verify the image tag is the one you expect
2. Run `kubectl get deployment --namespace=kube-system` and verify the `external-dns` deployment is running and is 
    healthy
3. Apply a new Ingress and verify ExternalDNS creates the expected DNS record in Route53

## Something went wrong

ExternalDNS does it's work in idempotent batches based on applied Ingress', so nothing will break by deleting it. If
something went wrong with the upgrade, you can safely uninstall and reinstall ExternalDNS by following the steps below:

1. Run `make uninstall` to delete associated resources
2. Run `make clean` to delete generated files
3. Follow the steps in [INSTALL.md](INSTALL.md) to reinstall ExternalDNS

