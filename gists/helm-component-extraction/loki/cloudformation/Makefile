######################################################################
# User input START ###################################################
######################################################################
# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

######################################################################
# User input END #####################################################
######################################################################

configure:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	cat templates/policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > policy.yaml


install:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -f policy.yaml || (echo "policy.yaml does not exist. Run 'make configure' first"; exit 1)
	aws cloudformation deploy \
			--template-file policy.yaml \
			--stack-name okctl-autoscalerpolicy-$(CLUSTER_NAME) \
			--capabilities CAPABILITY_NAMED_IAM \

	eksctl create iamserviceaccount \
		--cluster $(CLUSTER_NAME) \
		--name autoscaler \
		--namespace kube-system \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-AutoscalerServiceAccountPolicy").Arn') \
		--override-existing-serviceaccounts \
		--approve

