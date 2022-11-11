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
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	AWS_REGION=$$(aws configure get region)
	AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --query Account --output text)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(AWS_ACCOUNT_ID)" || (echo "Unable to acquire account id"; exit 1)

	cat templates/dynamodb-policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		AWS_REGION=$${AWS_REGION} \
		AWS_ACCOUNT_ID=$${AWS_ACCOUNT_ID} \
		envsubst > dynamodb-policy.yaml


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

