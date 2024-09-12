ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.54.0

ROOT_DIR = $(shell pwd)
PROTO_DIR = $(ROOT_DIR)/protos
COVERAGE_DIR = $(ROOT_DIR)/coverage
COV_UNIT_DIR = $(COVERAGE_DIR)/unit
COV_INTEGRATION_DIR = $(COVERAGE_DIR)/integration

$(GOLANGCI_LINT): $(GOBIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)

.PHONY: protos
protos:
	protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative $(PROTO_DIR)/*.proto

.PHONY: test
test: integration unit

.PHONY: integration
integration: $(GOLEAK)
	mkdir -p $(COV_INTEGRATION_DIR) || true
	go test -tags=integration -timeout 30m -cover ./... -args -test.gocoverdir=$(COV_INTEGRATION_DIR) 

.PHONY: unit
unit:
	mkdir -p $(COV_UNIT_DIR) || true
	go test -tags=unit -cover ./... -args -test.gocoverdir=$(COV_UNIT_DIR)

.PHONY: coverage
coverage: test
	go tool covdata textfmt -i="$(COV_INTEGRATION_DIR),$(COV_UNIT_DIR)" -o=$(COVERAGE_DIR)/tmp.cov
	go tool cover -func=$(COVERAGE_DIR)/tmp.cov
	grep -Ev '(testutils\.go|.*\.pb\.go)' $(COVERAGE_DIR)/tmp.cov > $(COVERAGE_DIR)/total.cov
	

PHONY: view-coverage
view-coverage: $(COVERAGE_DIR)/total.cov
	go tool cover -html=$(COVERAGE_DIR)/total.cov

PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run