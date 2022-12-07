######################################################################
# User input START ###################################################
######################################################################

# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

######################################################################
# User input END #####################################################
######################################################################

base-data:
	$(eval AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text))
	$(eval AWS_REGION := $(shell aws configure get region))
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
base-data-test:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)

data:
	$(eval AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text))
	$(eval AWS_REGION := $(shell aws configure get region))
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))

test-data:
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	@test -n "$(CLIENT_ID)" || (echo "Unable to acquire client ID"; exit 1)

values.yaml: test-data
	$(eval CERTIFICATE_ARN := $(shell aws acm list-certificates | yq '.CertificateSummaryList[] | select(.DomainName == "argocd.${CLUSTER_NAME}.oslo.systems").CertificateArn'))
	@test -n "$(CERTIFICATE_ARN)" || (echo "Unable to acquire certificate ARN"; exit 1)
	cat templates/values.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		ISSUER_URL=$(ISSUER_URL) \
		CLIENT_ID=$(CLIENT_ID) \
		CERTIFICATE_ARN=$(CERTIFICATE_ARN) \
		envsubst > values.yaml

argocd-secret.yaml:
	$(eval SERVER_SECRET := $(shell uuidgen))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/secret-key" --value="$(SERVER_SECRET)" --type SecureString --overwrite
	cat templates/argocd-secret.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > argocd-secret.yaml

cfn-argocd-client.yaml:
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	cat templates/cfn-argocd-client.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst > cfn-argocd-client.yaml

cognito-client: cfn-argocd-client.yaml base-data
	aws cloudformation deploy \
		--template-file cfn-argocd-client.yaml \
		--stack-name okctl-identitypoolclient-${CLUSTER_NAME}-argocd
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))
	$(eval CLIENT_SECRET := $(shell aws cognito-idp describe-user-pool-client --user-pool-id=$(USER_POOL_ID) --client-id=$(CLIENT_ID) | jq '.UserPoolClient.ClientSecret'))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/client-secret" --value="$(CLIENT_SECRET)" --type SecureString --overwrite

namespace.yaml:
	cp templates/namespace.yaml namespace.yaml

configure: data values.yaml namespace.yaml argocd-secret.yaml
	echo "✅ Configuration complete"

bootstrap:
	kubectl apply -f namespace.yaml
	helm repo add argo https://argoproj.github.io/argo-helm
	helm install --namespace argocd argocd argo/argo-cd --version 3.26.12 --values values.yaml 

install:
	@test -f values.yaml || (echo "Missing values.yaml file. Run 'make configure' first" && exit 1)
	yq \
		'.spec.source.helm.values = load_str("values.yaml")' \
		application.yaml | \
		kubectl apply -f -

uninstall:
	kubectl delete -f application.yaml || true
	helm delete argocd --namespace argocd || true
	kubectl delete -f argocd-secret.yaml || true
	aws ssm delete-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/secret-key" || true
	aws ssm delete-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/client-secret" || true
	aws cloudformation delete-stack --stack-name okctl-identitypoolclient-${CLUSTER_NAME}-argocd || true
