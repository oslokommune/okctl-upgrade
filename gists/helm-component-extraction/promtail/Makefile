configure:
	echo "No configuration needed"

install:
	@test -f values.yaml || (echo "Missing values.yaml file. Run 'make configure' first" && exit 1)
	yq \
		'.spec.source.helm.values = load_str("values.yaml")' \
		application.yaml | \
		kubectl apply -f -

uninstall:
	kubectl delete -f application.yaml

