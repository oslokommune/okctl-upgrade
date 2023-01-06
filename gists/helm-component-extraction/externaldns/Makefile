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
# ExternalDNS ########################################################
######################################################################
external-dns-data: base-data
	$(eval HOSTED_ZONE_ID := $(shell aws route53 list-hosted-zones | yq '.HostedZones[] | select(.Name == "${CLUSTER_NAME}.oslo.systems.").Id' | cut -d '/' -f3))

external-dns-data-test: base-data-test
	@test -n "$(HOSTED_ZONE_ID)" || (echo "Unable to acquire hosted zone ID"; exit 1)

cfn-service_account_policy.yaml: base-data-test
	cat templates/cfn-service_account_policy.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		envsubst > cfn-service_account_policy.yaml

install-service-account-policy: base-data-test
	aws cloudformation deploy \
		--stack-name okctl-externaldns-$(CLUSTER_NAME) \
		--template-file cfn-service_account_policy.yaml \
		--capabilities CAPABILITY_NAMED_IAM

uninstall-service-account-policy: base-data base-data-test
	aws cloudformation delete-stack \
		--stack-name okctl-externaldns-$(CLUSTER_NAME)

uninstall-service-user-role: base-data base-data-test
	$(eval ROLE_ARN := $(shell eksctl get iamserviceaccount --cluster julius --namespace=kube-system -o json | jq -r '.[] | select(.metadata.name == "external-dns").status.roleARN'))
	$(eval ROLE_NAME := $(shell echo $(ROLE_ARN) | cut -d '/' -f2))
	$(eval POLICY_ARN := $(shell aws iam list-attached-role-policies --role-name $(ROLE_NAME) | jq -r '.AttachedPolicies[].PolicyArn'))
	aws iam detach-role-policy \
		--role-name $(ROLE_NAME) \
		--policy-arn $(POLICY_ARN)
	aws iam delete-role --role-name $(ROLE_NAME)

install-service-user: base-data-test
	eksctl create iamserviceaccount \
		--cluster $(CLUSTER_NAME) \
		--name "external-dns" \
		--namespace "kube-system" \
		--attach-policy-arn $$(aws iam list-policies | jq -r '.Policies[] | select(.PolicyName == "okctl-$(CLUSTER_NAME)-ExternalDNSServiceAccountPolicy").Arn') \
		--override-existing-serviceaccounts \
		--approve

uninstall-service-user: base-data-test uninstall-service-user-role uninstall-service-account-policy
	kubectl delete --namespace=kube-system serviceaccount external-dns

clusterrole.yaml:
	cat templates/clusterrole.yaml > clusterrole.yaml

clusterrolebinding.yaml:
	cat templates/clusterrolebinding.yaml > clusterrolebinding.yaml

deployment.yaml: external-dns-data-test
	cat templates/deployment.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		HOSTED_ZONE_ID=$(HOSTED_ZONE_ID) \
		envsubst > deployment.yaml

install-kubernetes-manifests: clusterrole.yaml clusterrolebinding.yaml deployment.yaml
	kubectl apply -f clusterrole.yaml
	kubectl apply -f clusterrolebinding.yaml
	kubectl apply -f deployment.yaml

uninstall-kubernetes-manifests:
	kubectl delete --ignore-not-found -f clusterrole.yaml
	kubectl delete --ignore-not-found -f clusterrolebinding.yaml
	kubectl delete --ignore-not-found -f deployment.yaml

######################################################################
# User facing ########################################################
######################################################################

configure: external-dns-data deployment.yaml cfn-service_account_policy.yaml clusterrole.yaml clusterrolebinding.yaml
	@echo "✅ Configuring ExternalDNS complete"

install: external-dns-data install-service-account-policy install-service-user install-kubernetes-manifests
	@echo "✅ Installing ExternalDNS complete"

uninstall: base-data uninstall-kubernetes-manifests uninstall-service-user
	@echo "✅ Uninstalling ExternalDNS complete"
