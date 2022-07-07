# EKS upgrade

This document describes how to use the [upgrade script](upgrade.sh) in this repository in order to upgrade your EKS cluster.

The upgrade script basically just does what is described in step 1-3 in the official guide, https://eksctl.io/usage/cluster-upgrade/, which is:


> 1. upgrade control plane version with `eksctl upgrade cluster`
> 2. replace each of the nodegroups by creating a new one and deleting the old one
> 3. update default add-ons:
>     - `kube-proxy`
>     - `aws-node`
>     - `coredns`

## Prerequisites

**Make sure** you have already followed this guide previsouly to get to EKS 1.20:
https://github.com/oslokommune/okctl-upgrade/blob/main/gists/bump-eks-to-1-20/README.md. The exception to this if you have used Okctl to create a 1.20 cluster, as then it does not need upgrading.

## Tips

* If the upgrade script breaks, or you want to customize it in any way, it's not that hard to just edit the script to your needs.
* All steps in the scripts are written to be idempotent. This means if the script breaks, or if you edit it and want to re-run it, you can, and it should still work.

# Step 1: Upgrade Okctl environments

Download a version of Okctl runs on the EKS version you need, and run `okctl upgrade`. The [changelog of Okctl](https://github.com/oslokommune/okctl/releases) is best suited to find which version works with which verison of EKS.

This is to ensure that `okctl apply cluster` and other commands support the version of EKS you are running on.

(ToDo: Update this guide with versions so users don't have to read the changelog.)

# Step 2: Download or update tools

The upgrade script expects the following tools to exist on your machine, so make sure to install these.

## aws CLI

Follow instructions in https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

## jq

jq is a tool for parsing JSON.

### Linux / apt

```sh
sudo apt-get install jq
```

### macOS / Linux via Homebrew

```sh
brew install jq
```

### Other

See: https://stedolan.github.io/jq/download/.

## yq

yq is a tool for parsing YAML.

### Linux / snap

```sh
snap install yq
```

###  macOS / Linux via Homebrew

```sh
brew install yq
```

### Other

See https://github.com/mikefarah/yq.

## watch

### Linux

No need to install, this ususally comes preinstalled in most distributions.

### macOS

```
brew install watch
```

# Step 3: Prepare applications

## Avoid downtime

To avoid downtime, **make sure** you have completed the steps described in this guide: https://github.com/oslokommune/okctl-upgrade/tree/main/gists/bump-eks-to-1-20#prepare-applications

## If upgrading to EKS 1.22 or later: Update manifests

**Note!** If you do not follow the steps below, your applications **probably will stop working** after upgrading to Kubernetes 1.22.

Kubernetes 1.22 stops supporting some resources. The following guide describe these in detail: https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html#update-1.22. For applications made with `okctl apply application`, you don't need to read that guide, just follow the steps below.

### Update Ingress resources

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

For every file in the result, edit it and replace the `apiVersion` so it becomes like this:

```yaml
apiVersion: networking.k8s.io/v1
```

# Step 4: Monitor everything while upgrading

It's nice to see that stuff changes while upgrading, so while we run the upgrade script in the next step, we want to monitor pods and nodes. We'll monitor these things in a separate terminal.

* Open a new terminal window. Log in to your kubernetes cluster. The default Okctl way is to run `okctl venv` with your usual arguments.

```sh
okctl venv ... 
```

Then start monitoring pods:

```sh
watch -n 1 kubectl get pod --all-namespaces -o wide
```

* Open a new terminal tab. Log in to your kubernetes cluster like above.

Then start monitoring nodes:

```sh
watch -n 4 kubectl get node -o wide
```

* Open a new terminal tab. Log in to your kubernetes cluster like above.

In the following command, replace my-cluster-dev.yaml with your Okctl cluster manifest.

```sh
CLUSTER_NAME=$(yq e '.metadata.name' "my-cluster-dev.yaml")
```

Then start monitoring node groups:

```sh
watch -n 15 eksctl get nodegroup --cluster $CLUSTER_NAME
```

# Step 5: Run the upgrade

## Log in to environment

Log into the correct AWS environment. The default way to do this is:

```sh
export AWS_PROFILE=my-profile
aws sso login
```

## Download upgrade script

Download latest version upgrade script (it may be updated at any time):

```sh
curl --silent --location "https://raw.githubusercontent.com/oslokommune/okctl-upgrade/main/gists/bump-eks/upgrade.sh"
```

## Run the upgrade

### Usage

```sh
USAGE:
upgrade.sh <cluster-manifest file> <aws-region> <EKS target version> [dry-run={false|true}] | tee logfile.txt

cluster-manifest file      The Okctl cluster manifest
aws-region                 AWS region
EKS target version         Example: 1.21
dry-run                 Default true. Set to false to actually run upgrade.
```

### Tips

* You can upgrade only one minor version at the time. So if you are on EKS 1.20 and want to upgrade to EKS 1.22, you must first upgrade to 1.21, then to 1.22.

* :info: The `tee` thing in the following commands is there to create a nice upgrade log. You do not have to, but we recommend storing this (in git or somewhere else), because
  * It gives a pretty nice and accurate way of telling what you have done with your cluster, which can be useful for future reference.
  * It helps immensely for debugging in case something wrong happens.

### Example, upgrading EKS 1.20 to 1.22

```sh
# Dry run the upgrade, hoping to catch any errors before actually upgrading
./upgrade.sh cluster-dev.yaml eu-west-1 1.21 | tee "eks-upgrade-log-$(date +"%Y-%m-%dx%H-%M-%S").log"

# Actually run the upgrade
./upgrade.sh cluster-dev.yaml eu-west-1 1.21 dry-run=false | tee "eks-upgrade-log-$(date +"%Y-%m-%dx%H-%M-%S").log"

# Store the logs
git add *.log
git commit -m "Upgraded to EKS 1.22"

# Dry run the upgrade, hoping to catch any errors before actually upgrading
./upgrade.sh cluster-dev.yaml eu-west-1 1.22 | tee "eks-upgrade-log-$(date +"%Y-%m-%dx%H-%M-%S").log"

# Actually run the upgrade
./upgrade.sh cluster-dev.yaml eu-west-1 1.22 dry-run=false | tee "eks-upgrade-log-$(date +"%Y-%m-%dx%H-%M-%S").log"

# Store the log
git add *.log
git commit -m "Upgraded to EKS 1.22"
```

# Something wrong happened

## eksctl delete nodegroup cannot evict pods

Abort/CTRL+C your execution of `eksctl delete nodegroup` if it's running, because we will be running the
command below, which we don't want to run at the same time.

In the following command, replace `/tmp/eks-upgrade/1-21` with `/tmp/eks-upgrade/1-22` or whatever version you're running on. Run:

```shell
/tmp/eks-upgrade/1-21/kubectl drain -l 'alpha.eksctl.io/nodegroup-name=ng-generic' --ignore-daemonsets --delete-emptydir-data
```

This should output exactly which pods that cannot be evicted due to its `PodDisruptionBudget` (or for other reasons?).

## Apps have downtime when draining nodes

* Your app's Deployment must have `replicas: 2`.
* You need a working `PodDisruptionBudget`.

How to setup these correctly is described in https://github.com/oslokommune/okctl-upgrade/blob/main/gists/bump-eks-to-1-20/README.md.

# Resources

- https://eksctl.io/usage/cluster-upgrade/

## Commands

## Set desiredCapacity

If you really need to, you can in nodegroup_config.yaml set desiredCapacity to `1`. Or run:

```
aws autoscaling set-desired-capacity --desired-capacity 1 --auto-scaling-group-name eksctl-my-cluster-nodegroup-ng-generic-1-20-1c-NodeGroup-DFG36JFJY345
```

to have less down time. This is at the cost of having more nodes than needed.
