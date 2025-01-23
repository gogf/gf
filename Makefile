SHELL := /bin/bash

# execute "go mod tidy" on all folders that have go.mod file
.PHONY: tidy
tidy:
	$(eval files=$(shell find . -name go.mod))
	@set -e; \
	for file in ${files}; do \
		goModPath=$$(dirname $$file); \
		if ! echo $$goModPath | grep -q "testdata"; then \
			cd $$goModPath; \
			go mod tidy; \
			cd -; \
		fi \
	done

# execute "golangci-lint" to check code style
.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

# make version to=v2.4.0
.PHONY: version
version:
	@set -e; \
	newVersion=$(to); \
	./.set_version.sh ./ $$newVersion; \
	echo "make version to=$(to) done"
