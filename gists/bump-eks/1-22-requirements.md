This document describes preparations you must do before upgrading EKS to 1.22.  

## Preparations before upgrading to EKS 1.22

If you do not follow the below steps, your application **probably will stop working**.

AWS describes the necessary changes we need to take into account in detail: https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html#update-1.22. However, we have attempted to extract everything that is necessary, so you shouldn't need to read that guide.

### Update Ingress manifests

First `cd` into the directory where you store your Kubernetes manifests/YAMLs. The default way in Okctl is to put these in your IAC repository.

```yaml
cd your-okctl-iac-repository
```

Find occurrences of old Ingress resources by running

```
grep -nRsH "apiVersion: networking.k8s.io/v1beta1"
grep -nRsH "apiVersion: extensions/v1beta1"
```

Example output:

```
$ grep -nRsH "apiVersion: networking.k8s.io/v1beta1"
--
infrastructure/applications/okctl-reference-app/base/ingress.yaml:1:apiVersion: networking.k8s.io/v1beta1
infrastructure/applications/okctl-reference-app/base/ingress.yaml-2-kind: Ingress
--
infrastructure/applications/hello/base/ingress.yaml:1:apiVersion: networking.k8s.io/v1beta1
infrastructure/applications/hello/base/ingress.yaml-2-kind: Ingress
```

#### Update the YAML

* For every file in the result, edit it and replace the `apiVersion` so it becomes like this:

```yaml
apiVersion: networking.k8s.io/v1
```

* Also, you need to change the Ingress YAML somewhat:

Old ingress:

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
# ... (some more YAML)
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        backend:
          serviceName: echo
          servicePort: 80
```

You must rewrite this to:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
# ... (some more YAML)
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: test
                port:
                  number: 80
```

Some useful links regarding this change:
* [Example](https://docs.konghq.com/kubernetes-ingress-controller/latest/concepts/ingress-versions/)
* [Kubernetes Ingress documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/#the-ingress-resource)

### Make sure Okctl version is x.x.x or later (ToDo: update this when Okctl supports 1.22)

Okctl doesn't support EKS 1.22 yet. We will update this guide when it is.

You can still upgrade to EKS 1.22, but you cannot create a cluster from scratch with version 1.21. You
will have to create it with version 1.20 and upgade to 1.21.

We need to bump the AWS load balancer controller to 2.4.1 or later. `okctl upgrade` handles this for us.

From update-1.22 documentation:

> - If you currently have the AWS Load Balancer Controller deployed to your cluster, you must update it to version `2.4.1` before updating your cluster to Kubernetes version `1.22`.

```
okctl venv -a aws-profile -c my-cluster.yaml

# TODO
kubectl -n kube-system get pod -o=jsonpath='{$.spec.template.spec.containers[:0].image}'
```
