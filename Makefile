SHELL := /bin/bash

.PHONY: tidy
tidy:
	$(eval files=$(shell find . -name go.mod))
	@set -e; \
	for file in ${files}; do \
		goModPath=$$(dirname $$file); \
		cd $$goModPath; \
		go mod tidy; \
		cd -; \
	done

.PHONY: lint
lint:
	golangci-lint run

gftidy:
	$(eval files=$(shell find . -name go.mod))
	@set -e; \
	# GITHUB_REF_NAME=v2.4.0; \
	if [[ $$GITHUB_REF_NAME =~ "v" ]]; then \
		latestVersion=$$GITHUB_REF_NAME; \
	else \
		latestVersion=latest; \
	fi; \
	for file in ${files}; do \
		goModPath=$$(dirname $$file); \
		if [[ $$goModPath =~ "./contrib" || $$goModPath =~ "./cmd/gf" || $$goModPath =~ "./example" ]]; then \
			echo ""; \
			echo "processing dir: $$goModPath"; \
			# Do not modify the order of any of the following sentences \
			cd $$goModPath; \
			go mod tidy; \
			go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@$$latestVersion{{end}}" -m all | grep "^github.com/gogf/gf/contrib" | xargs -L1 go get -v; \
			go get -v github.com/gogf/gf/v2@$$latestVersion; \
			go mod tidy; \
			cd -; \
		fi \
	done