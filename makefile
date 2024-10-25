ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.58.0

MOCKGEN ?= $(GOBIN)/mockgen
MOCKGEN_VERSION ?= v0.3.0

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

.PHONY: get-mockgen
get-mockgen: $(MOCKGEN) ## Download mockgen locally if necessary.
$(MOCKGEN): $(GOBIN)
	go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

.PHONY: mocks
mocks: $(MOCKGEN) $(ROOT_DIR)/client_mock.go $(ROOT_DIR)/connection_provider_mock.go $(PROTO_DIR)/auth_grpc_mock.pb.go $(PROTO_DIR)/index_grpc_mock.pb.go $(PROTO_DIR)/transact_grpc_mock.pb.go $(PROTO_DIR)/types_mock.pb.go $(PROTO_DIR)/user-admin_grpc_mock.pb.go $(PROTO_DIR)/vector-db_grpc_mock.pb.go $(PROTO_DIR)/transact_service_vector_search_client_mock.pb.go

$(ROOT_DIR)/client_mock.go: $(ROOT_DIR)/client.go
	$(MOCKGEN) --source $< --destination $@ --package avs

$(ROOT_DIR)/connection_provider_mock.go: $(ROOT_DIR)/connection_provider.go
	$(MOCKGEN) --source $< --destination $@ --package avs

$(PROTO_DIR)/auth_grpc_mock.pb.go: $(PROTO_DIR)/auth_grpc.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/index_grpc_mock.pb.go: $(PROTO_DIR)/index_grpc.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/transact_grpc_mock.pb.go: $(PROTO_DIR)/transact_grpc.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/transact_service_vector_search_client_mock.pb.go: $(PROTO_DIR)/transact_service_vector_search_client.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/types_mock.pb.go: $(PROTO_DIR)/types.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/user-admin_grpc_mock.pb.go: $(PROTO_DIR)/user-admin_grpc.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

$(PROTO_DIR)/vector-db_grpc_mock.pb.go: $(PROTO_DIR)/vector-db_grpc.pb.go
	$(MOCKGEN) --source $< --destination $@ --package protos

.PHONY: test
test: unit integration 

.PHONY: integration
integration:
	mkdir -p $(COV_INTEGRATION_DIR) || true
	go test -tags=integration -timeout 30m -cover ./... -args -test.gocoverdir=$(COV_INTEGRATION_DIR) 

.PHONY: unit
unit: mocks
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
lint: $(GOLANGCI_LINT) mocks
	$(GOLANGCI_LINT) run