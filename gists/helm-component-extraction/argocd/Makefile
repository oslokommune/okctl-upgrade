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
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq eval '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq eval '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))
	$(eval CERTIFICATE_ARN := $(aws acm list-certificates | yq eval '.CertificateSummaryList[] | select(.DomainName == "argocd.$(CLUSTER_NAME).oslo.systems").CertificateArn'))

test-data:
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	@test -n "$(CLIENT_ID)" || (echo "Unable to acquire client ID"; exit 1)
	@test -n "$(CERTIFICATE_ARN)" || (echo "Unable to acquire certificate ARN"; exit 1)

values.yaml: test-data
	cat templates/values.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		ISSUER_URL=$(ISSUER_URL) \
		CLIENT_ID=$(CLIENT_ID) \
		CERTIFICATE_ARN=$(CERTIFICATE_ARN) \
		envsubst > values.yaml

configure: data test-data values.yaml
	echo "✅ Configuration complete"

install:
	@test -f values.yaml || (echo "Missing values.yaml file. Run 'make configure' first" && exit 1)
	yq \
		'.spec.source.helm.values = load_str("values.yaml")' \
		application.yaml | \
		kubectl apply -f -

uninstall:
	kubectl delete -f application.yaml

