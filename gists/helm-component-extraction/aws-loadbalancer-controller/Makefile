##############################
## User input START         ##
##############################
# Use `aws eks list-clusters` to get the cluster name
# Example: CLUSTER_NAME=kjoremiljo
CLUSTER_NAME?=

##############################
## User input END           ##
##############################

configure:
	@test -n "$(AWS_PROFILE)" || (echo "Missing AWS_PROFILE environment variable. export AWS_PROFILE with the relevant profile" && exit 1)
	@test -n "$(CLUSTER_NAME)" || (echo "Missing CLUSTER_NAME value. Edit the Makefile and configure the user input" && exit 1)
	cat templates/values.yaml | \
		CLUSTER_NAME=$(CLUSTER_NAME) \
		AWS_REGION=$$(aws configure get region) \
		VPC_ID=$$(aws eks describe-cluster --name $(CLUSTER_NAME) --query cluster.resourcesVpcConfig.vpcId --output text) \
		envsubst > values.yaml

install:
	@test -f values.yaml || (echo "Missing values.yaml file. Run 'make configure' first" && exit 1)
	yq \
		'.spec.source.helm.values = load_str("values.yaml")' \
		application.yaml | \
		kubectl apply -f -

uninstall:
	kubectl delete -f application.yaml
