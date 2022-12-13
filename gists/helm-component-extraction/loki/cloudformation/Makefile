######################################################################
# User input START ###################################################
######################################################################
# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

######################################################################
# User input END #####################################################
######################################################################

data:
	$(eval AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text))
	$(eval AWS_REGION := $(shell aws configure get region))

test-data:
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(AWS_ACCOUNT_ID)" || (echo "Unable to acquire account id"; exit 1)

test-install:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -f dynamodb-policy.yaml || (echo "dynamodb-policy.yaml does not exist. Run 'make configure' first"; exit 1)
	@test -f s3bucket-policy.yaml || (echo "s3bucket-policy.yaml does not exist. Run 'make configure' first"; exit 1)
	@test -f s3bucket.yaml || (echo "s3bucket.yaml does not exist. Run 'make configure' first"; exit 1)

dynamodb-policy.yaml: test-data
	cat templates/dynamodb-policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		AWS_REGION=$(AWS_REGION) \
		AWS_ACCOUNT_ID=$(AWS_ACCOUNT_ID) \
		envsubst > dynamodb-policy.yaml

s3bucket-policy.yaml: test-data
	cat templates/s3bucket-policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > s3bucket-policy.yaml

s3bucket.yaml: test-data
	cat templates/s3bucket.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > s3bucket.yaml

configure: data dynamodb-policy.yaml s3bucket-policy.yaml s3bucket.yaml
	@echo "Configuration complete"

install: test-install
	aws cloudformation deploy \
			--template-file dynamodb-policy.yaml \
			--stack-name okctl-dynamodbpolicy-$(CLUSTER_NAME)-loki \
			--capabilities CAPABILITY_NAMED_IAM \
	
	aws cloudformation deploy \
			--template-file s3bucket-policy.yaml \
			--stack-name okctl-s3bucketpolicy-$(CLUSTER_NAME)-loki \
			--capabilities CAPABILITY_NAMED_IAM \
	
	aws cloudformation deploy \
			--template-file s3bucket.yaml \
			--stack-name okctl-s3bucket-$(CLUSTER_NAME)-okctl-$(CLUSTER_NAME)-loki \
			--capabilities CAPABILITY_NAMED_IAM \

	eksctl create iamserviceaccount \
		--cluster $(CLUSTER_NAME) \
		--name loki \
		--namespace monitoring \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-LokiDynamoDBServiceAccountPolicy").Arn') \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-LokiS3ServiceAccountPolicy").Arn') \
		--override-existing-serviceaccounts \
		--approve

