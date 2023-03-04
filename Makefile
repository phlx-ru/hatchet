.PHONY: help
# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: generate
# run golang generate
generate:
	go generate ./...

.PHONY: vendor
# make ./vendor folder with dependencies
vendor:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: test
# makes go test ./...
test:
	@go test -race -parallel 10 ./...

.PHONY: lint
# runs linter from golangci-lint docker image
lint:
	@docker run --rm -v $$(pwd):/app \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v $$(go env GOCACHE):/cache/go \
		-v $$(go env GOPATH)/pkg:/go/pkg \
		-w /app golangci/golangci-lint:latest-alpine \
		golangci-lint run --verbose --timeout 5m

.PHONY: lintfix
# runs linter from golangci-lint docker image with --fix flag
lintfix:
	@docker run --rm -v $$(pwd):/app \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v $$(go env GOCACHE):/cache/go \
		-v $$(go env GOPATH)/pkg:/go/pkg \
		-w /app golangci/golangci-lint:latest-alpine \
		golangci-lint run --verbose --timeout 5m --fix

.PHONY: changelog
# run changelog generator (use `yarn global add changelog.md` before run this command)
changelog:
	@changelog
