.PHONY: dependency
dependency:
	@go mod download
	@go mod vendor

.PHONY: verify
verify: go-fmt go-vet go-lint test

.PHONY: go-vet
go-vet:
	@go vet -v ./...

.PHONY: go-fmt
go-fmt:
	@git ls-files '*.go' | grep -v 'vendor/' | xargs gofmt -s -w

.PHONY: go-lint
go-lint: install-golint
	@golangci-lint run

.PHONY: install-golint
install-golint:
	GOLINT_CMD=$(shell command -v golangci-lint 2> /dev/null)
ifndef GOLINT_CMD
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

.PHONY: clean-vendor
clean-vendor:
	@find ./vendor -type d | xargs rm -rf

.PHONY: clean-test-cache
clean-test-cache:
	@go clean -testcache ./...

.PHONY: test
test: clean-test-cache
	@go test -v ./... -cover -coverprofile=coverage.out -race -run ./...

.PHONY: code-coverage
code-coverage:
	@go tool cover -html=coverage.out
