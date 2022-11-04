BUF_LINT_INPUT := .
BUF_BREAKING_INPUT := .
BUF_BREAKING_AGAINST_INPUT ?= .git\#branch=main
BUF_FORMAT_INPUT := .
BUF_VERSION ?= v1.9.0
GOLANGCI_SKIP_FILES = pkg/api/spark/v1/*.pb.go
GO_GET_PKGS := $(GO_GET_PKGS) \
	github.com/srikrsna/protoc-gen-gotag@v0.6.2 \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 \
	github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.0 \
	github.com/envoyproxy/protoc-gen-validate@v0.6.2

include make/go/bootstrap.mk
include make/go/buf.mk
include make/go/go.mk
include make/go/dep_protoc.mk
include make/go/dep_protoc_gen_go.mk
include make/go/dep_protoc_gen_go_grpc.mk
include make/go/dep_go_releaser.mk

# HOOKS

bufgeneratedeps:: $(BUF) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)

# if you want to clean a directory before generating proto use rm here
.PHONY: bufgenerateclean
bufgenerateclean::

# Called before linting, testing etc. to make sure all outputs are generated before linting or testing
.PHONY: bufgeneratesteps
bufgeneratesteps::
	buf generate

# Called when make build is run
.PHONY: gorelease
gorelease:: $(GO_RELEASER)
	goreleaser release --rm-dist

# Called when make release is run
.PHONY: gobuild
gobuild:: $(GO_RELEASER)
	goreleaser build --rm-dist
