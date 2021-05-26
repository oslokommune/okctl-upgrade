SHELL         = bash
.PHONY: release

help: ## Print this menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

release: ## Make a release. This target assumes that running `git tag` returns the version to release.
	# Because goreleaser expects this tag, we'll read it here as well.
	$(eval VERSION := $(shell git tag))

	cd ${VERSION} && \
	VERSION=${VERSION} \
		goreleaser release \
			--config ../.goreleaser.yaml \
			--rm-dist

release-test: ## Test making a release. Example usage: make release-test VERSION=0.0.50
	# Because goreleaser expects this tag, we'll read it here as well.
	cd ${VERSION} && \
		goreleaser release \
			--config ../.goreleaser.yaml \
			--rm-dist --skip-publish --snapshot
