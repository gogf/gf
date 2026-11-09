SHELL := /bin/bash

# execute "go mod tidy" on all folders that have go.mod file
.PHONY: tidy
tidy:
	./.make_tidy.sh

# execute "golangci-lint" to check code style
# go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

# make branch to=v2.4.0
# When 'to' is omitted, the next patch version derived from the latest git tag is used.
.PHONY: branch
branch:
	@set -e; \
	echo "Switching to master branch..."; \
	git checkout master; \
	echo "Pulling latest changes from master..."; \
	git pull origin master; \
	newVersion=$(to); \
	if [ -z "$$newVersion" ]; then \
		latestTag=$$(git tag --list 'v*' --sort=-v:refname | head -n 1); \
		if [ -z "$$latestTag" ]; then \
			echo "Error: no version tag found. Usage: make branch to=vX.Y.Z"; \
			exit 1; \
		fi; \
		baseVersion=$${latestTag%%-*}; \
		major=$$(echo "$$baseVersion" | cut -d. -f1); \
		minor=$$(echo "$$baseVersion" | cut -d. -f2); \
		patch=$$(echo "$$baseVersion" | cut -d. -f3); \
		if [ -z "$$patch" ]; then \
			echo "Error: cannot parse version from tag $$latestTag. Usage: make branch to=vX.Y.Z"; \
			exit 1; \
		fi; \
		newVersion=$${major}.$${minor}.$$((patch + 1)); \
		echo "Auto detected next version: $$newVersion (latest tag: $$latestTag)"; \
	fi; \
	branchName=fix/$$newVersion; \
	echo "Creating and switching to branch $$branchName from master..."; \
	git checkout -b $$branchName; \
	echo "Branch $$branchName created successfully!"; \
	echo "Next step: make version to=$$newVersion"

# make version to=v2.4.0
.PHONY: version
version:
	@set -e; \
	newVersion=$(to); \
	$(MAKE) -C cmd/gf pack; \
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

# manage docker services for local development
# usage: make docker or make docker cmd=start svc=mysql
.PHONY: docker
docker:
	@if [ -z "$(cmd)" ]; then \
		./.github/workflows/scripts/docker-services.sh; \
	else \
		./.github/workflows/scripts/docker-services.sh $(cmd) $(svc) $(extra); \
	fi

# make agents agent=claude [action=unlink] [force=1]
.PHONY: agents
agents:
	@cd cmd/gf && go run . agents $(if $(agent),--agent=$(agent)) $(if $(action),--action=$(action)) $(if $(force),--force)
