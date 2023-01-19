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
base-data-test:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(AWS_ACCOUNT_ID)" || (echo "Unable to acquire account id"; exit 1)

######################################################################
# ArgoCD Cognito Client ##############################################
######################################################################

BUILD_DIR=build

$(BUILD_DIR):
	mkdir $(BUILD_DIR)

cognito-client-data: base-data
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))

cognito-client-data-test: base-data-test
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)

$(BUILD_DIR)/cfn-argocd-client.yaml: cognito-client-data-test $(BUILD_DIR)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	cat templates/cfn-argocd-client.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst > $(BUILD_DIR)/cfn-argocd-client.yaml

configure-cognito-client: cognito-client-data cognito-client-data-test $(BUILD_DIR)/cfn-argocd-client.yaml
	@echo "✅ Configuring Cognito client complete"

install-cognito-client: cognito-client-data cognito-client-data-test
	@test -f "$(BUILD_DIR)/cfn-argocd-client.yaml" || (echo "cfn-argocd-client.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file $(BUILD_DIR)/cfn-argocd-client.yaml \
		--stack-name okctl-identitypoolclient-${CLUSTER_NAME}-argocd

install-cognito-parameters: cognito-client-data cognito-client-data-test
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-argocd").ClientId'))
	$(eval CLIENT_SECRET := $(shell aws cognito-idp describe-user-pool-client --user-pool-id=$(USER_POOL_ID) --client-id=$(CLIENT_ID) | jq '.UserPoolClient.ClientSecret'))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/client-secret" --value="$(CLIENT_SECRET)" --type SecureString --overwrite > /dev/null

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

$(BUILD_DIR)/cfn-argocd-certificate.yaml: certificate-data-test $(BUILD_DIR)
	cat templates/cfn-argocd-certificate.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		HOSTED_ZONE_ID=$(HOSTED_ZONE_ID) \
		envsubst > $(BUILD_DIR)/cfn-argocd-certificate.yaml

configure-certificate: certificate-data certificate-data-test $(BUILD_DIR)/cfn-argocd-certificate.yaml
	@echo "✅ Configuring certificate complete"

install-certificate: certificate-data certificate-data-test
	@test -f "$(BUILD_DIR)/cfn-argocd-certificate.yaml" || (echo "cfn-argocd-certificate.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file $(BUILD_DIR)/cfn-argocd-certificate.yaml \
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

$(BUILD_DIR)/namespace.yaml: $(BUILD_DIR)
	cp templates/namespace.yaml $(BUILD_DIR)/namespace.yaml

$(BUILD_DIR)/values.yaml: argocd-data argocd-data-test $(BUILD_DIR)
	cat templates/values.yaml | \
		AWS_REGION=$(AWS_REGION) \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		ISSUER_URL=$(ISSUER_URL) \
		CLIENT_ID=$(CLIENT_ID) \
		CERTIFICATE_ARN=$(CERTIFICATE_ARN) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst '$$AWS_REGION,$$CLUSTER_NAME,$$ISSUER_URL,$$CLIENT_ID,$$CERTIFICATE_ARN,$$USER_POOL_ID' > $(BUILD_DIR)/values.yaml

$(BUILD_DIR)/argocd-secret.yaml: argocd-data argocd-data-test $(BUILD_DIR)
	$(eval SERVER_SECRET := $(shell uuidgen))
	@aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/argocd/secret-key" --value="$(SERVER_SECRET)" --type SecureString --overwrite | cat
	cat templates/argocd-secret.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > $(BUILD_DIR)/argocd-secret.yaml

configure-helm-chart: argocd-data $(BUILD_DIR)/namespace.yaml $(BUILD_DIR)/values.yaml $(BUILD_DIR)/argocd-secret.yaml
	@echo "✅ Configuring helm chart complete"

install-helm-chart:
	@test -f "$(BUILD_DIR)/namespace.yaml" || (echo "namespace.yaml does not exist. Run 'make configure'"; exit 1)
	@test -f "$(BUILD_DIR)/values.yaml" || (echo "values.yaml is missing. Run 'make configure'"; exit 1)
	@test -f "$(BUILD_DIR)/argocd-secret.yaml" || (echo "argocd-secret.yaml is missing. Run 'make configure'"; exit 1)
	kubectl apply -f $(BUILD_DIR)/namespace.yaml
	kubectl apply -f $(BUILD_DIR)/argocd-secret.yaml
	helm upgrade --install \
		--namespace argocd \
		--version 5.17.1 \
		--values $(BUILD_DIR)/values.yaml \
		argocd templates/chart

uninstall-helm-chart:
	helm uninstall --namespace argocd argocd || true
	kubectl delete -f $(BUILD_DIR)/argocd-secret.yaml || true
	kubectl delete -f $(BUILD_DIR)/namespace.yaml || true

######################################################################
# Repository deploy key ##############################################
######################################################################

deploy-key-data: base-data
	$(eval REPOSITORY_NAME := $(shell basename -s .git `git config --get remote.origin.url`))
deploy-key-data-test: base-data-test
	@test -n "$(REPOSITORY_NAME)" || (echo "Unable to acquire repository name"; exit 1)
	@test -n "$(GITHUB_TOKEN)" || (echo "GITHUB_TOKEN not set"; exit 1)

$(BUILD_DIR)/deploy-key: deploy-key-data-test $(BUILD_DIR)
	ssh-keygen -t ed25519 -C "okctl-$(CLUSTER_NAME)-argocd" -f $(BUILD_DIR)/deploy-key -N ""

$(BUILD_DIR)/deploy-key-secret: deploy-key-data-test $(BUILD_DIR)
	cat templates/argocd-repository-secret.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		REPOSITORY_NAME=$(REPOSITORY_NAME) \
		envsubst > $(BUILD_DIR)/deploy-key-secret.yaml

configure-deploy-key: deploy-key-data $(BUILD_DIR)/deploy-key $(BUILD_DIR)/deploy-key-secret
	@echo "✅ Configuring deploy key complete"

install-deploy-key-parameters: deploy-key-data deploy-key-data-test
	aws ssm put-parameter --name="/okctl/$(CLUSTER_NAME)/github/deploykeys/oslokommune/$(REPOSITORY_NAME)/privatekey" --value=file://$(BUILD_DIR)/deploy-key --type SecureString --overwrite > /dev/null
	rm $(BUILD_DIR)/deploy-key

install-deploy-key: deploy-key-data deploy-key-data-test
	kubectl apply -f $(BUILD_DIR)/deploy-key-secret.yaml
	@curl \
		-X POST \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${GITHUB_TOKEN}"\
		-H "X-GitHub-Api-Version: 2022-11-28" \
		"https://api.github.com/repos/oslokommune/$(REPOSITORY_NAME)/keys" \
			-d "{\"title\":\"okctl-iac-$(CLUSTER_NAME)\",\"key\":\"`cat $(BUILD_DIR)/deploy-key.pub`\",\"read_only\":false}"

uninstall-deploy-key: deploy-key-data deploy-key-data-test
	kubectl delete -f $(BUILD_DIR)/deploy-key-secret.yaml || true
	aws ssm delete-parameter --name="/okctl/$(CLUSTER_NAME)/github/deploykeys/oslokommune/$(REPOSITORY_NAME)/privatekey" > /dev/null || true

######################################################################
# Existing ArgoCD Applications #######################################
######################################################################
remove-finalizers:
	kubectl --namespace=argocd get applications.argoproj.io -o yaml | \
		yq '.items[] |= del(.metadata.finalizers)' | \
		kubectl apply -f -

applications-snapshot: $(BUILD_DIR)
	kubectl -n argocd get applications -o yaml \
		| yq '.items[] |= del(.status)' \
		| yq '.items[] |= del(.metadata.creationTimestamp)' \
		| yq '.items[] |= del(.metadata.generation)' \
		| yq '.items[] |= del(.metadata.resourceVersion)' \
		| yq '.items[] |= del(.metadata.uid)' > $(BUILD_DIR)/applications-snapshot.yaml

install-applications-snapshot:
	kubectl apply -f $(BUILD_DIR)/applications-snapshot.yaml
