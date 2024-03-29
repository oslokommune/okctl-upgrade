######################################################################
# User input START ###################################################
######################################################################

# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

######################################################################
# User input END #####################################################
######################################################################

SERVICEUSER_NAME=external-dns-ext

BUILD_DIR=build
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

base-data:
	$(eval AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text))
	$(eval AWS_REGION := $(shell aws configure get region))
base-data-test:
	@test -n "$(CLUSTER_NAME)" || (echo "CLUSTER_NAME is not set. See user input section in Makefile"; exit 1)
	@test -n "$(AWS_PROFILE)" || (echo "AWS_PROFILE is not set. Export the AWS_PROFILE environment with the relevant profile"; exit 1)
	@test -n "$(AWS_REGION)" || (echo "Unable to acquire region"; exit 1)
	@test -n "$(AWS_ACCOUNT_ID)" || (echo "Unable to acquire account id"; exit 1)

######################################################################
# ExternalDNS ########################################################
######################################################################
external-dns-data: base-data
	$(eval HOSTED_ZONE_ID := $(shell aws route53 list-hosted-zones | yq '.HostedZones[] | select(.Name == "${CLUSTER_NAME}.oslo.systems.").Id' | cut -d '/' -f3))

external-dns-data-test: base-data-test
	@test -n "$(HOSTED_ZONE_ID)" || (echo "Unable to acquire hosted zone ID"; exit 1)

$(BUILD_DIR)/cfn-service_account_policy.yaml: base-data-test $(BUILD_DIR)
	cat templates/cfn-service_account_policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > $(BUILD_DIR)/cfn-service_account_policy.yaml

install-service-account-policy: base-data-test
	aws cloudformation deploy \
		--stack-name okctl-externaldns-$(CLUSTER_NAME) \
		--template-file $(BUILD_DIR)/cfn-service_account_policy.yaml \
		--capabilities CAPABILITY_NAMED_IAM

uninstall-service-account-policy: base-data base-data-test
	aws cloudformation delete-stack \
		--stack-name okctl-externaldns-$(CLUSTER_NAME)

uninstall-service-user-role: base-data base-data-test
	$(eval ROLE_ARN := $(shell eksctl get iamserviceaccount --cluster julius --namespace=kube-system -o json | jq -r '.[] | select(.metadata.name == "$(SERVICEUSER_NAME)").status.roleARN'))
	$(eval ROLE_NAME := $(shell echo $(ROLE_ARN) | cut -d '/' -f2))
	$(eval POLICY_ARN := $(shell aws iam list-attached-role-policies --role-name $(ROLE_NAME) | jq -r '.AttachedPolicies[].PolicyArn'))
	aws iam detach-role-policy \
		--role-name $(ROLE_NAME) \
		--policy-arn $(POLICY_ARN) || true
	aws iam delete-role --role-name $(ROLE_NAME) || true

install-service-user: base-data-test
	eksctl create iamserviceaccount \
		--cluster $(CLUSTER_NAME) \
		--name "$(SERVICEUSER_NAME)" \
		--namespace "kube-system" \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-ExternalDNSServiceAccountPolicy").Arn') \
		--override-existing-serviceaccounts \
		--approve

uninstall-service-user: base-data-test uninstall-service-user-role uninstall-service-account-policy
	kubectl delete --namespace=kube-system serviceaccount external-dns

$(BUILD_DIR)/clusterrole.yaml: $(BUILD_DIR)
	cat templates/clusterrole.yaml > $(BUILD_DIR)/clusterrole.yaml

$(BUILD_DIR)/clusterrolebinding.yaml: $(BUILD_DIR)
	cat templates/clusterrolebinding.yaml | \
		SERVICEUSER_NAME=$(SERVICEUSER_NAME) \
		envsubst > $(BUILD_DIR)/clusterrolebinding.yaml

$(BUILD_DIR)/deployment.yaml: external-dns-data-test $(BUILD_DIR)
	cat templates/deployment.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		HOSTED_ZONE_ID=$(HOSTED_ZONE_ID) \
		SERVICEUSER_NAME=$(SERVICEUSER_NAME) \
		AWS_REGION=$(AWS_REGION) \
		envsubst > $(BUILD_DIR)/deployment.yaml

install-kubernetes-manifests:
	test -f $(BUILD_DIR)/clusterrole.yaml || (echo "Missing $(BUILD_DIR)/clusterrole.yaml. Did you run make configure?"; exit 1)
	test -f $(BUILD_DIR)/clusterrolebinding.yaml || (echo "Missing $(BUILD_DIR)/clusterrolebinding.yaml. Did you run make configure?"; exit 1)
	test -f $(BUILD_DIR)/deployment.yaml || (echo "Missing $(BUILD_DIR)/deployment.yaml. Did you run make configure?"; exit 1)
	kubectl apply -f $(BUILD_DIR)/clusterrole.yaml
	kubectl apply -f $(BUILD_DIR)/clusterrolebinding.yaml
	kubectl apply -f $(BUILD_DIR)/deployment.yaml

uninstall-kubernetes-manifests:
	kubectl delete --ignore-not-found -f $(BUILD_DIR)/clusterrole.yaml
	kubectl delete --ignore-not-found -f $(BUILD_DIR)/clusterrolebinding.yaml
	kubectl delete --ignore-not-found -f $(BUILD_DIR)/deployment.yaml

######################################################################
# User facing ########################################################
######################################################################

configure: external-dns-data $(BUILD_DIR)/deployment.yaml $(BUILD_DIR)/cfn-service_account_policy.yaml $(BUILD_DIR)/clusterrole.yaml $(BUILD_DIR)/clusterrolebinding.yaml
	@echo "✅ Configuring ExternalDNS complete"

install: external-dns-data install-service-account-policy install-service-user install-kubernetes-manifests
	@echo "✅ Installing ExternalDNS complete"

uninstall: base-data uninstall-kubernetes-manifests uninstall-service-user
	@echo "✅ Uninstalling ExternalDNS complete"

clean:
	rm -rf $(BUILD_DIR)
