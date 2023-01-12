# Prepare applications

Do all the steps prefixed with "Step:".

## Step: Remove potential old Okctl configurations

### Remove dnsPolicy

```sh
grep -A 1 -B 50 -nRsH "dnsPolicy:"
```

should yield no results. If there are, remove the line `dnsPolicy: Default`.

### Use two securitygroups

Applications that use database need two security groups.

Run

```sh
kubectl get sgp -A
```

Every security group should have to security groups, like this:

```
my-app      my-sgp      ["sg-0ab340d9f94c4c0a1","sg-0ed2bd34231484bb3"]
```

and NOT like this:

```
my-app      my-sgp      ["sg-0ab340d9f94c4c0a1"]
```

If the latter is the case, you need to update your security group policies to include both the app security group, and the cluster security group. The first should always be there (configured by Okctl). The last can be found by running:

```sh
aws eks describe-cluster --name my-cluster-dev --query cluster.resourcesVpcConfig.clusterSecurityGroupId
```

## Step: Prepare apps for downtime

Your applications need one of the following configurations.

### Alternative 1: No downtime

* Deployment: Use `RollingUpdate` strategy
* Deployment: Use `replicas: 2` or more
* Create a `PodDisruptionBudget`

Google how to apply using `maxUnavailable=0` and `type: RollingUpdate`.

Example:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
  namespace: hello
spec:
  replicas: 2
  strategy:
    rollingUpdate:
      maxUnavailable: 0
```

### Alternative 2: For stateful applications that must not co-exist

* Deployment: Use `Recreate` strategy
* Deployment: Use `replicas: 1`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
  namespace: hello
spec:
  replicas: 1
  strategy:
    type: Recreate
```

## Step: Add node selectors to pods using PVCs

Before we can bump nodes, we need to make sure that pods that use volumes (via PVCs), spawn on a node in the same AZ as the
volumes. If not the pod will not start, as it cannot find the PV.

To do this, we need to specify which AZ pods in Kubernetes should spawn on. The AZ should be the same as the AZ of the PVC the
application is using.

### List PVCs

To get a list of PVCs, run:

```shell
kubectl get pvc -A

kubectl -n NAMESPACE describe pv PV_ID
# Replace NAMESPACE and PV_ID with values from above command. PV_ID = VOLUME.
```

Look for a label like this `failure-domain.beta.kubernetes.io/zone=eu-west-1c`, `eu-west-1c` is the AZ.

### Update deployments

Now update all your Deployments (or Pods or StatefulSets) that refers to the PVCs of these PVs to use a `nodeSelector` with the same AZ as the PVC.

So for instance, in `deployment.yaml`, you can change from

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
spec:
  template:
    spec:
      containers:
        - name: hello
```

to

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
spec:
  template:
    spec:
      nodeSelector:
        failure-domain.beta.kubernetes.io/zone: eu-west-1c
      containers:
        - name: hello
```

## Step: Add a PodDisruptionBudget for every application

**NOTE** Applications using PVC's cannot use PodDisruptionBudgets - see note in the bottom of this title.

A `PodDisruptionBudget` can be used to make sure for instance 1 pod is always in Running state when draining nodes. For more details, see [documentation](https://kubernetes.io/docs/tasks/run-application/configure-pdb/).

For each application, create a `infrastructure/applications/hello/base/pod-disruption-budget.yaml` with contents:

```yaml
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: hello
  namespace: hello
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: hello
```

**Important!** Your app's deployment must have `replicas: 2` or more. If not, it will be impossible to drain a node, because the pod can never be moved to a new node. Follow steps above to set `replicas` on your deployment.

Update `infrastructure/applications/hello/base/kustomization.yaml` so it includes `pod-disruption-budget.yaml`. For isntance:

```yaml
resources:
- service.yaml
- ingress.yaml
- namespace.yaml
- deployment.yaml
- pod-disruption-budget.yaml
```

### For applications using PVCs: No downtime is impossible

Applications using PVCs cannot use PodDisruptionBudgets, so no downtime will be impossible during node replacements (i.e. an upgrade).

Because our PVC's are using the mode `ReadWriteOnce`, it will become impossible for Kubernetes to one replica of a deployment on a new node, at the same time as it is running a replica on an old node - because `ReadWriteOnce` does not allow it. The
definition of ReadWriteOnce is:

> ReadWriteOnce – the volume can be mounted as read-write by a single node

Hence, we cannot use a PodDisruptionBudget for applications using PVCs. If we did the pod would never be evicted, because the pod/replica trying to spawn on a new node cannot spawn because the volume is being used by a pod/replica on an old node. We have tested and confirmed this behavior.

The solution is:
* Accept some downtime during upgrading / node replacement
* Don't use PVCs, use S3 or something else for your applications
* (Far fetched: Implement support for a PVC driver that supports `ReadWriteMany`)

See discussion here: https://stackoverflow.com/a/62216783/915441

## Step: Apply changes

Run

```shell
git add .
git commit -m "Add node selector to deployments"
git push
```

ArgoCD will then update your apps.

If you don't want to wait, you can run

```shell
CLUSTER_NAME="my-cluster" # See "eksctl get cluster"
kustomize build infrastructure/applications/hello/overlays/$CLUSTER_NAME | kubectl apply -f -
```
