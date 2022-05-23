This guide describes how to upgrade EKS from 1.19 to 1.20 in an EKS cluster.

**Note!** If you have not persisted your Loki logs with a PVC or in S3, they will disappear!

# Update tools

* Download the latest version of [eksctl](https://github.com/weaveworks/eksctl/releases).
* Download kubectl CLI version 1.20

# Bump EKS control plane

Bump your EKS control plane, by running.

```shell
okctl venv ...

eksctl get cluster

# Replace my-cluster with the name of cluster you want to upgrade from above command.
eksctl upgrade cluster --name my-cluster --version 1.20
eksctl upgrade cluster --name my-cluster --version 1.20 --approve 
```

# Bump EC2 nodes in your cluster

## Add node selectors to pods using PVCs

Before we can bump nodes, we need to make sure that pods that use volumes (via PVCs), spawn on a node in the same AZ as the volumes. If not the pod will not start, as it cannot find the PV.

To do this, we need to specify which AZ pods in Kubernetes should spawn on. The AZ should be the same as the AZ of the PVC the application is using.


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
  name: helloworld
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
  name: helloworld
spec:
  template:
    spec:
      nodeSelector:
        failure-domain.beta.kubernetes.io/zone: eu-west-1c
      containers:
        - name: hello
```

Deploy these changes (with `kubectl apply ...`, or `git` commit and push if you use ArgoCD).

## Spin up new nodes

We're using 3 nodes to ensure we have 1 node for every AZ. We need one in every AZ to ensure that any applications using PVCs can
be placed on a node in the same AZ as the PVC. For instance: If some-app use a PVC in AZ B, we need to have a node in AZ B as
well. If it was possible to specify   

```shell
CLUSTER_NAME="my-cluster"
REGION="eu-west-1"

cat <<EOF >nodegroup_config.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: $CLUSTER_NAME
  region: $REGION
nodeGroups:
EOF

for AZ_ID in a b c
do
  AZ="${REGION}${AZ_ID}"

  cat <<EOF >>nodegroup_config.yaml
  - name: "ng-generic-1-20-1${AZ_ID}"
    availabilityZones: ["$AZ"]
    instanceType: "m5.large"
    # desiredCapacity: 1 # TODO
    minSize: 0
    maxSize: 10
    labels:
      pool: ng-generic-$AZ
    tags:
      k8s.io/cluster-autoscaler/enabled: "true"
      k8s.io/cluster-autoscaler/$CLUSTER_NAME: owned
    privateNetworking: true
EOF
done

```

Replace
* `CLUSTER_NAME` with the name from `eksctl get cluster`
* `REGION` with your region.

Now, create a new nodegroup:

```shell
eksctl create nodegroup --config-file=nodegroup_config.yaml
```

# Update EKS add-on: vpc-cni

The recommended vpc-cni addon version for all Kubernetes versions is `1.11.0-eksbuild.1`
([source](https://docs.aws.amazon.com/eks/latest/userguide/managing-vpc-cni.html)).

Get the IAM role the VPC-CNI addon uses:

```shell
eksctl get addon --cluster ykctl --name vpc-cni -o json
```

See field "IAMRole", it should be something like

```
arn:aws:iam::123456789012:role/eksctl-mycluster-addon-vpc-cni-Role1-DMGPR03HYLWR
```

```shell
CLUSTER_NAME="mycluster"
ROLE_ARN="arn:aws:iam::123456789012:role/eksctl-mycluster-addon-vpc-cni-Role1-DMGPR03HYLWR"

# Update default addons
eksctl utils update-kube-proxy --cluster=$CLUSTER_NAME --approve
eksctl utils update-aws-node --cluster=$CLUSTER_NAME --approve
eksctl utils update-coredns --cluster=$CLUSTER_NAME --approve

# Update vpc-cni addon
eksctl update addon \
  --cluster $CLUSTER_NAME \
  --name vpc-cni \
  --version 1.7.10-eksbuild.1 \
  --service-account-role-arn $ROLE_ARN

eksctl update addon \
  --cluster $CLUSTER_NAME \
  --name vpc-cni \
  --version 1.8.0-eksbuild.1 \ #todo
  --service-account-role-arn $ROLE_ARN
  
eksctl update addon \
  --cluster $CLUSTER_NAME \
  --name vpc-cni \
  --version 1.9.3-eksbuild.1 \
  --service-account-role-arn $ROLE_ARN

eksctl update addon \
  --cluster $CLUSTER_NAME \
  --name vpc-cni \
  --version 1.10.3-eksbuild.1 \
  --service-account-role-arn $ROLE_ARN

eksctl update addon \
  --cluster $CLUSTER_NAME \
  --name vpc-cni \
  --version 1.11.0-eksbuild.1 \
  --service-account-role-arn $ROLE_ARN

```


## Drain old node(s)

(Draining also sets a taint on the nodes, i.e. prohibits new pods to be scheduled on them. So there is no need to taint nodes before draining them.)

To see which nodes are going to be drained, run:

```shell
kubectl drain -l 'alpha.eksctl.io/nodegroup-name=ng-generic' --ignore-daemonsets --delete-local-data --dry-run=client
```

Verify that the list of nodes above are indeed the nodes you want to drain.

Now actually drain nodes:

```shell
kubectl drain -l 'alpha.eksctl.io/nodegroup-name=ng-generic' --ignore-daemonsets --delete-local-data
```

## Delete the old nodegroup

Verify that no pods is running on the nodes in the old nodegroup:

```shell
kubectl get pod -o wide
```

Then delete the nodegroup:

```shell
eksctl delete nodegroup --cluster ykctl ng-generic
```
