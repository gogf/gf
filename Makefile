SHELL := /bin/bash

# commit changes with AI-generated commit message
.PHONY: up
up:
	@if git diff --quiet HEAD && git diff --cached --quiet && [ -z "$$(git ls-files --others --exclude-standard)" ]; then \
		echo "No changes to commit"; \
		exit 0; \
	fi
	@git add -A
	@echo "Analyzing changes and generating commit message via AI..."
	@set -e; \
	MSG=$$(git diff --cached --stat && echo "---" && git diff --cached | head -2000 | \
		claude -p "Analyze the git diff above and generate a concise commit message (single line, max 72 chars, lowercase, no quotes). Output only the commit message itself, nothing else." \
		--model haiku) || { echo "Error: Claude command failed"; exit 1; }; \
	COMMIT_MSG=$$(echo "$$MSG" | tail -1); \
	if [ -z "$$COMMIT_MSG" ]; then \
		echo "Error: Failed to generate commit message"; \
		exit 1; \
	fi; \
	echo "Commit: $$COMMIT_MSG"; \
	git commit -m "$$COMMIT_MSG" && \
	git push origin

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
.PHONY: branch
branch:
	@set -e; \
	newVersion=$(to); \
	if [ -z "$$newVersion" ]; then \
		echo "Error: 'to' variable is required. Usage: make branch to=vX.Y.Z"; \
		exit 1; \
	fi; \
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

# manage docker services for local development
# usage: make docker or make docker cmd=start svc=mysql
.PHONY: docker
docker:
	@if [ -z "$(cmd)" ]; then \
		./.github/workflows/scripts/docker-services.sh; \
	else \
		./.github/workflows/scripts/docker-services.sh $(cmd) $(svc) $(extra); \
	fi
