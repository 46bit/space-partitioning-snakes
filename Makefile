.DEFAULT_GOAL := help

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## run the testsuite
	ginkgo -r

.PHONY: regenerate
regenerate: ## regenerate all automatically-generated code
	go generate ./...
