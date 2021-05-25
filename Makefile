GODOG_VERSION=v0.11.0
GODOG := github.com/cucumber/godog/cmd/godog@$(GODOG_VERSION)
GOIMPORTS_VERSION=v0.1.1
GOIMPORTS := golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
GOLANGCI_VERSION=v1.39.0
GOLANGCI := github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
PROTOC_GEN_GO_VERSION=v1.26.0
PROTOC_GEN_GO := google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
GO_TOOLS := $(GODOG) $(GOLANGCI) $(PROTOC_GEN_GO) $(GOIMPORTS)

BACKEND_GO_FILES := $(shell find . -type f -name '*.go' -not -name '*.pb.go')
BDD_FEATURE_FILES := $(shell find .behaviours -type f -name '*.feature')

GRPC_API_V1_SPEC_SOURCE := boundaries/grpc/v1/proto/service.proto
GRPC_API_V1_SPEC_TARGET := \
	boundaries/grpc/v1/proto/service_grpc.pb.go \
	boundaries/grpc/v1/proto/service.pb.go

BDD_BINARY := cmd/identity-service-bdd-test/godog.test

BINARIES := \
	cmd/identity-service/main

$(GO_TOOLS):
	@echo "[tools] Installing $@"
	@go install $@

.PHONY: tidy
tidy: | $(GOIMPORTS)
	@echo "[go]    Tidying"
	@go mod tidy
	@gofmt -s -w .
	@goimports -w .

gen: $(GRPC_API_V1_SPEC_TARGET)

$(GRPC_API_V1_SPEC_TARGET): $(GRPC_API_V1_SPEC_SOURCE) | $(PROTOC_GEN_GO)
	@echo "[proto] Generating $(@D)"
	@cd $(@D) && \
		protoc \
			--go_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			$(notdir $(GRPC_API_V1_SPEC_SOURCE))

.PHONY: run-service
run-service: gen
	@echo "[go]    Running service"
	@cd cmd/identity-service && go run main.go

.PHONY: run-client
run-client: gen
	@echo "[go]    Running client"
	@cd cmd/identity-service-client && go run main.go

bdd_binary: $(BDD_BINARY)

$(BDD_BINARY): gen $(BACKEND_GO_FILES) | $(GODOG)
	@echo "[godog] Building binary"
	@cd $(@D) && godog build

.PHONY: bdd
bdd: $(BDD_BINARY) $(BDD_FEATURE_FILES)
	@echo "[godog] Executing"
	@$(BDD_BINARY) -f progress -c 100 .behaviours

.PHONY: clean
clean:
	@rm -rf $(GRPC_API_V1_SPEC_TARGET) $(BDD_BINARY) $(BINARIES)

.PHONY: build
build: $(BINARIES)

$(BINARIES): cmd/%/main: gen $(BACKEND_GO_FILES)
	@echo "[go]    Building" $@
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-w -extldflags "-static"' -o ./$@ ./$@.go
