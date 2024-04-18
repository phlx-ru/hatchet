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

.PHONY: init
# init env
init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
	@GOGC=95 golangci-lint run --verbose --timeout 5m

.PHONY: lintfix
# runs linter from golangci-lint docker image with --fix flag
lintfix:
	@GOGC=95 golangci-lint run --fix --verbose --timeout 5m

.PHONY: changelog
# run changelog generator (use `yarn global add changelog.md` before run this command)
changelog:
	@changelog
