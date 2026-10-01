.PHONY: help test vet tidy hooks-install

.DEFAULT_GOAL := help

help:
	@echo "operatorconfig — XDG config discovery, Viper load, go-keyring secrets"
	@echo ""
	@echo "  make test   go test ./..."
	@echo "  make vet    go vet ./..."
	@echo "  make tidy   go mod tidy"
	@echo "  make hooks-install  Install Lefthook git hooks (once per clone)"

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

hooks-install:
	@command -v lefthook >/dev/null 2>&1 || { \
		if command -v brew >/dev/null 2>&1; then brew install lefthook; \
		else go install github.com/evilmartians/lefthook@latest; fi; }
	@command -v lefthook >/dev/null 2>&1 || { echo "lefthook not on PATH; add $$(go env GOPATH)/bin"; exit 1; }
	lefthook install
