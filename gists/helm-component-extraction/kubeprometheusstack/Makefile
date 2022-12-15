##############################
## User input START         ##
##############################
# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

##############################
## User input END           ##
##############################

base-data:
	$(eval AWS_REGION := $(shell aws configure get region))

base-data-test:
	@test -n "$(AWS_PROFILE)" || (echo "Missing AWS_PROFILE environment variable. export AWS_PROFILE with the relevant profile" && exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire AWS_REGION" && exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "Missing CLUSTER_NAME value. Edit the Makefile and configure the user input" && exit 1)

######################################################################
# Cognito Client #####################################################
######################################################################
cognito-client-data: base-data
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))

cognito-client-data-test: base-data-test
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)

cfn-cognito-client.yaml: base-data cognito-client-data cognito-client-data-test
	cat templates/cfn-cognito-client.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		USER_POOL_ID=$(USER_POOL_ID) \
		envsubst > cfn-cognito-client.yaml

configure-cognito-client: cfn-cognito-client.yaml
	@echo "✅ Configuring Cognito Client complete"

install-cognito-client:
	@test -f "cfn-cognito-client.yaml" || (echo "cfn-cognito-client.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file cfn-cognito-client.yaml \
		--stack-name okctl-identitypoolclient-${CLUSTER_NAME}-grafana

######################################################################
# Certificate ########################################################
######################################################################

certificate-data: base-data
	$(eval HOSTED_ZONE_ID := $(shell aws route53 list-hosted-zones | yq '.HostedZones[] | select(.Name == "${CLUSTER_NAME}.oslo.systems.").Id' | cut -d '/' -f3))

certificate-data-test: base-data-test
	@test -n "$(HOSTED_ZONE_ID)" || (echo "Unable to acquire hosted zone ID"; exit 1)

cfn-certificate.yaml: certificate-data certificate-data-test
	cat templates/cfn-certificate.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		HOSTED_ZONE_ID=$(HOSTED_ZONE_ID) \
		envsubst > cfn-certificate.yaml

configure-certificate: cfn-certificate.yaml
	@echo "✅ Configuring certificate complete"

install-certificate: certificate-data certificate-data-test
	@test -f "cfn-certificate.yaml" || (echo "cfn-certificate.yaml does not exist. Run 'make configure'"; exit 1)
	aws cloudformation deploy \
		--template-file cfn-certificate.yaml \
	  --stack-name okctl-certificate-${CLUSTER_NAME}-grafana-${CLUSTER_NAME}-oslo-systems

######################################################################
# Div ################################################################
######################################################################

cfn-cloudwatch-datasource.yaml: base-data-test
	cat templates/cfn-cloudwatch-datasource.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > cfn-cloudwatch-datasource.yaml

######################################################################
# Helm chart #########################################################
######################################################################

chart-data: base-data
	$(eval CERTIFICATE_ARN := $(shell aws acm list-certificates | yq '.CertificateSummaryList[] | select(.DomainName == "grafana.${CLUSTER_NAME}.oslo.systems").CertificateArn'))
	$(eval USER_POOL_ID := $(shell aws cognito-idp list-user-pools --max-results 10 | yq '.UserPools[] | select(.Name == "okctl-$(CLUSTER_NAME)-userpool").Id'))
	$(eval CLIENT_ID := $(shell aws cognito-idp list-user-pool-clients --user-pool-id $(USER_POOL_ID) | yq '.UserPoolClients[] | select(.ClientName == "okctl-$(CLUSTER_NAME)-grafana").ClientId'))

chart-data-test: base-data-test
	@test -n "$(CERTIFICATE_ARN)" || (echo "Unable to acquire certificate ARN"; exit 1)
	@test -n "$(USER_POOL_ID)" || (echo "Unable to acquire user pool ID"; exit 1)
	@test -n "$(CLIENT_ID)" || (echo "Unable to acquire Cognito client ID"; exit 1)

values.yaml: chart-data chart-data-test
	cat templates/values.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		AWS_REGION=$$(aws configure get region) \
		COGNITO_CLIENT_ID=$(CLIENT_ID) \
		CERTIFICATE_ARN=$(CERTIFICATE_ARN) \
		envsubst > values.yaml

configure-helm-chart: values.yaml
	@echo "✅ Configuring Helm chart complete"

install-helm-chart:
	@test -f values.yaml || (echo "Missing values.yaml file. Run 'make configure' first" && exit 1)
	yq \
		'.spec.source.helm.values = load_str("values.yaml")' \
		application.yaml | \
		kubectl apply -f -

uninstall:
	kubectl delete -f application.yaml
