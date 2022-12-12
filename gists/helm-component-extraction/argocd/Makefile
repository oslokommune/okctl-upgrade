######################################################################
# User input START ###################################################
######################################################################

# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=julius

######################################################################
# User input END #####################################################
######################################################################

base-data:
	$(eval AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text))
	$(eval AWS_REGION := $(shell aws configure get region))
base-data-test:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(AWS_ACCOUNT_ID)" || (echo "Unable to acquire account id"; exit 1)

######################################################################
# ArgoCD Cognito Client ##############################################
######################################################################

cognito-client-data: base-data
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))

cognito-client-data-test: base-data-test
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)

cfn-argocd-client.yaml: cognito-client-data-test
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	cat templates/cfn-argocd-client.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst > cfn-argocd-client.yaml

configure-cognito-client: cognito-client-data cognito-client-data-test cfn-argocd-client.yaml
	@echo "✅ Configuring Cognito client complete"

install-cognito-client: cognito-client-data cognito-client-data-test
	@test -f "cfn-argocd-client.yaml" || (echo "cfn-argocd-client.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file cfn-argocd-client.yaml \
		--stack-name okctl-identitypoolclient-${CLUSTER_NAME}-argocd
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))
	$(eval CLIENT_SECRET := $(shell aws cognito-idp describe-user-pool-client --user-pool-id=$(USER_POOL_ID) --client-id=$(CLIENT_ID) | jq '.UserPoolClient.ClientSecret'))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/client-secret" --value="$(CLIENT_SECRET)" --type SecureString --overwrite | cat

uninstall-cognito-client: base-data base-data-test
	aws cloudformation delete-stack --stack-name okctl-identitypoolclient-${CLUSTER_NAME}-argocd
	aws ssm delete-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/client-secret" || true

######################################################################
# ArgoCD Certificate     #############################################
######################################################################
certificate-data: base-data
	$(eval HOSTED_ZONE_ID := $(shell aws route53 list-hosted-zones | yq '.HostedZones[] | select(.Name == "${CLUSTER_NAME}.oslo.systems.").Id' | cut -d '/' -f3))

certificate-data-test: base-data-test
	@test -n "$(HOSTED_ZONE_ID)" || (echo "Unable to acquire hosted zone ID"; exit 1)

cfn-argocd-certificate.yaml: certificate-data-test
	cat templates/cfn-argocd-certificate.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		HOSTED_ZONE_ID=$(HOSTED_ZONE_ID) \
		envsubst > cfn-argocd-certificate.yaml

configure-certificate: certificate-data certificate-data-test cfn-argocd-certificate.yaml
	@echo "✅ Configuring certificate complete"

install-certificate: certificate-data certificate-data-test
	@test -f "cfn-argocd-certificate.yaml" || (echo "cfn-argocd-certificate.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file cfn-argocd-certificate.yaml \
	  --stack-name okctl-certificate-${CLUSTER_NAME}-argocd-${CLUSTER_NAME}-oslo-systems

uninstall-certificate: base-data base-data-test
	aws cloudformation delete-stack --stack-name okctl-certificate-${CLUSTER_NAME}-argocd-${CLUSTER_NAME}-oslo-systems

######################################################################
# ArgoCD Helm deployment #############################################
######################################################################
argocd-data: base-data
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	$(eval ISSUER_URL := $(shell aws eks describe-cluster --name $(CLUSTER_NAME) | jq -r '.cluster.identity.oidc.issuer'))
	$(eval CERTIFICATE_ARN := $(shell aws acm list-certificates | yq '.CertificateSummaryList[] | select(.DomainName == "argocd.${CLUSTER_NAME}.oslo.systems").CertificateArn'))
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))

argocd-data-test: base-data-test
	@test -n "$(ISSUER_URL)" || (echo "Unable to acquire issuer URL"; exit 1)
	@test -n "$(CERTIFICATE_ARN)" || (echo "Unable to acquire certificate ARN"; exit 1)
	@test -n "$(CLIENT_ID)" || (echo "Unable to acquire certificate ARN"; exit 1)

namespace.yaml:
	cp templates/namespace.yaml namespace.yaml

values.yaml: argocd-data argocd-data-test
	cat templates/values.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		ISSUER_URL=$(ISSUER_URL) \
		CLIENT_ID=$(CLIENT_ID) \
		CERTIFICATE_ARN=$(CERTIFICATE_ARN) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst '$$AWS_REGION,$$CLUSTER_NAME,$$ISSUER_URL,$$CLIENT_ID,$$CERTIFICATE_ARN,$$USER_POOL_ID' > values.yaml

argocd-secret.yaml: argocd-data argocd-data-test
	$(eval SERVER_SECRET := $(shell uuidgen))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/secret-key" --value="$(SERVER_SECRET)" --type SecureString --overwrite | cat
	cat templates/argocd-secret.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > argocd-secret.yaml

configure-helm-chart: argocd-data namespace.yaml values.yaml argocd-secret.yaml
	echo "✅ Configuring helm chart complete"

install-helm-chart:
	@test -f "namespace.yaml" || (echo "namespace.yaml does not exist. Run 'make configure'"; exit 1)
	@test -f "values.yaml" || (echo "values.yaml is missing. Run 'make configure'"; exit 1)
	@test -f "argocd-secret.yaml" || (echo "argocd-secret.yaml is missing. Run 'make configure'"; exit 1)
	kubectl apply -f namespace.yaml
	kubectl apply -f argocd-secret.yaml
	helm repo add argo https://argoproj.github.io/argo-helm
	helm install --namespace argocd argocd argo/argo-cd --version 3.26.12 --values values.yaml 

uninstall-helm-chart:
	helm uninstall --namespace argocd argocd
	kubectl delete -f argocd-secret.yaml
	kubectl delete -f namespace.yaml

######################################################################
# User facing            #############################################
######################################################################

applications-snapshot:
	kubectl -n argocd get applications -o yaml > applications.yaml

install-applications-snapshot:
	kubectl apply -f applications.yaml
