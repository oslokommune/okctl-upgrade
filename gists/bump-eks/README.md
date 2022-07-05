This guide describes how to upgrade EKS in an Okctl environment. The script supports multiple EKS versions.

**Make sure** you have already upgrade all components described in 
https://github.com/oslokommune/okctl-upgrade/blob/main/gists/bump-eks-to-1-20/README.md.

# Upgrade Okctl environments

Download a version of Okctl runs on the EKS version you need, and run `okctl upgrade`. The [changelog of Okctl](https://github.com/oslokommune/okctl/releases) is best suited to find which version works with which verison of EKS.

This is to ensure that `okctl apply cluster` and other commands support the version of EKS you are running on.

# Download or update tools

Download by using the commands, so we get the correct version expected by this upgrade.

## jq

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

# Prepare applications

## Avoid downtime

To avoid downtime, **make sure** you have completed the steps described in this guide: https://github.com/oslokommune/okctl-upgrade/tree/main/gists/bump-eks-to-1-20#prepare-applications

## If upgrading to EKS 1.22 or later: Update manifests

**Note!** If you do not follow the steps below, your applications **probably will stop working** after upgrading to Kubernetes 1.22.

In Kubernetes 1.22, there are some resources that stop working. The following guide describe these in detail: https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html#update-1.22. For applications made with `okctl apply application`, you don't need to read that guide, but you can just follow the steps below.

### Ingress

Find occurrences of old Ingresses by running

```
grep -A 1 -B 50 -nRsH "apiVersion: networking.k8s.io/v1beta1"
grep -A 1 -B 50 -nRsH "apiVersion: extensions/v1beta1"
```

For every file in the result, edit it and replace the `apiVersion` so it becomes like this:

```yaml
apiVersion: networking.k8s.io/v1
```

# Run the upgrade

## Log in to environment

Log into your kubernetes environment. The default way to do this is:

```sh
export AWS_PROFILE=my-profile
aws sso login
okctl venv -c cluster-dev.yaml
```

## Download upgrade script
Download latest version upgrade script (it may be updated at any time):

```sh
curl --silent --location "https://raw.githubusercontent.com/oslokommune/okctl-upgrade/main/gists/bump-eks/upgrade.sh"
```

## Run the upgrade

It's possible to upgrade only one minor version at the time. So if you are on EKS 1.20 and want to upgrade to EKS 1.22, you must run:

```sh
# Dry run, to see if everything is supposed to work.
./upgrade.sh cluster-dev.yaml eu-west-1 1.21

# Actually run the upgrade
./upgrade.sh cluster-dev.yaml eu-west-1 1.21 --dry-run=false

# Dry run, to see if everything is supposed to work.
./upgrade.sh cluster-dev.yaml eu-west-1 1.22

# Actually run the upgrade
./upgrade.sh cluster-dev.yaml eu-west-1 1.22 --dry-run=false
```

# Something wrong happened

## eksctl delete nodegroup cannot evict pods

Abort/CTRL+C your execution of `eksctl delete nodegroup` if it's running, because we will be running the
command below, which we don't want to run at the same time.

In the following command, replace `/tmp/eks-upgrade/1-21` with `/tmp/eks-upgrade/1-22` or whatever version you're running on. Run:

```shell
/tmp/eks-upgrade/1-21 drain -l 'alpha.eksctl.io/nodegroup-name=ng-generic' --ignore-daemonsets --delete-emptydir-data
```

This should output exactly which pods that cannot be evicted due to its `PodDisruptionBudget`.

## Apps have downtime when draining nodes

* Your app's Deployment must have `replicas: 2`.
* You need a working `PodDisruptionBudget`.

# Resources

- https://eksctl.io/usage/cluster-upgrade/

## Commands

## Set desiredCapacity

If you really need to, you can in nodegroup_config.yaml set desiredCapacity to `1`. Or run:

```
aws autoscaling set-desired-capacity --desired-capacity 1 --auto-scaling-group-name eksctl-my-cluster-nodegroup-ng-generic-1-20-1c-NodeGroup-DFG36JFJY345
```

to have less down time. This is at the cost of having more nodes than needed.
