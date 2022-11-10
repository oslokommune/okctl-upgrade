######################################################################
# User input START ###################################################
######################################################################
# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME=

######################################################################
# User input END #####################################################
######################################################################

configure:
	cat templates/role.yaml | \
		OIDC_URI=$$(aws eks describe-cluster --name $(CLUSTER_NAME) --query cluster.identity.oidc.issuer --output text | sed -e "s/^https:\/\///") \
		AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --query Account --output text) \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > role.yaml


install:
	aws cloudformation deploy \
			--template-file policy.yaml \
			--stack-name okctl-autoscalerpolicy-$(CLUSTER_NAME) \
			--capabilities CAPABILITY_NAMED_IAM \

	POLICY_ARN=
	eksctl create iamserviceaccount \
		--cluster $(CLUSTER_NAME) \
		--name autoscaler \
		--namespace kube-system \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-AutoscalerServiceAccountPolicy").Arn') \
		--override-existing-serviceaccounts \
		--approve

