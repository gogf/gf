SHELL := /bin/bash

# execute "go mod tidy" on all folders that have go.mod file
.PHONY: tidy
tidy:
	./.make_tidy.sh

# execute "golangci-lint" to check code style
.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

# make version to=v2.4.0
.PHONY: version
version:
	@set -e; \
	newVersion=$(to); \
	./.make_version.sh ./ $$newVersion; \
	echo "make version to=$(to) done"


# update submodules
.PHONY: subup
subup:
	@set -e; \
	echo "Updating submodules..."; \
	git submodule init;\
	git submodule update;

# update and commit submodules
.PHONY: subsync
subsync: subup
	@set -e; \
	echo "";\
	cd examples; \
	echo "Checking for changes..."; \
	if git diff-index --quiet HEAD --; then \
		echo "No changes to commit"; \
	else \
		echo "Found changes, committing..."; \
		git add -A; \
		git commit -m "examples update"; \
		git push origin; \
	fi; \
	cd ..;
