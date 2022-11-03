GO_BINS := $(GO_BINS) cmd/foo
DOCKER_BINS := $(DOCKER_BINS) foo

BUF_LINT_INPUT := .
BUF_BREAKING_INPUT := .
BUF_BREAKING_AGAINST_INPUT ?= .git\#branch=main
BUF_FORMAT_INPUT := .
BUF_VERSION ?= v1.9.0

include make/go/bootstrap.mk
include make/go/buf.mk
include make/go/go.mk
include make/go/dep_protoc_gen_go.mk

bufgeneratedeps:: $(BUF) $(PROTOC_GEN_GO)

# if you want to clean a directory before generating proto use rm here
.PHONY: bufgenerateclean
bufgenerateclean::

# Called before linting, testing etc. to make sure all outputs are generated before linting or testing
.PHONY: bufgeneratesteps
bufgeneratesteps::
	buf generate

# Called when make build is run
.PHONY: gorelease
gorelease::
	goreleaser build --rm-dist

# Called when make release is run
.PHONY: gobuild
gobuild::
	goreleaser release --rm-dist
