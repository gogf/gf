SHELL := /bin/bash

# execute "go mod tidy" on all folders that have go.mod file
.PHONY: tidy
tidy:
	./.make_tidy.sh

# execute "golangci-lint" to check code style
.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

# make branch to=v2.4.0
.PHONY: branch
branch:
	@set -e; \
	newVersion=$(to); \
	branchName=fix/$$newVersion; \
	echo "Switching to master branch..."; \
	git checkout master; \
	echo "Pulling latest changes from master..."; \
	git pull origin master; \
	echo "Creating and switching to branch $$branchName from master..."; \
	git checkout -b $$branchName; \
	echo "Branch $$branchName created successfully!"

# make version to=v2.4.0
.PHONY: version
version:
	@set -e; \
	newVersion=$(to); \
	./.make_version.sh ./ $$newVersion; \
	echo "make version to=$(to) done"

# make tag to=v2.4.0
.PHONY: tag
tag:
	@set -e; \
	newVersion=$(to); \
	echo "Switching to master branch..."; \
	git checkout master; \
	echo "Pulling latest changes from master..."; \
	git pull origin master; \
	echo "Creating annotated tag $$newVersion..."; \
	git tag -a $$newVersion -m "Release $$newVersion"; \
	echo "Pushing tag $$newVersion..."; \
	git push origin $$newVersion; \
	echo "Tag $$newVersion created and pushed successfully!"

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
