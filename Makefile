.PHONY: help test vet tidy

.DEFAULT_GOAL := help

help:
	@echo "operatorconfig — XDG config discovery, Viper load, go-keyring secrets"
	@echo ""
	@echo "  make test   go test ./..."
	@echo "  make vet    go vet ./..."
	@echo "  make tidy   go mod tidy"

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy
