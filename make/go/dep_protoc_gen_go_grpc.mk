# Managed by makego. DO NOT EDIT.

# Must be set
$(call _assert_var,MAKEGO)
$(call _conditional_include,$(MAKEGO)/base.mk)
$(call _assert_var,CACHE_VERSIONS)
$(call _assert_var,CACHE_BIN)

# Settable
# https://github.com/protocolbuffers/protobuf-go/releases 20220831 checked 20221004
# NOTE: This is temporary until the following fix is available in a release:
#   https://github.com/protocolbuffers/protobuf-go/commit/692f4a24f8dc0d375508fc41e657920d411b5b68
PROTOC_GEN_GO_GRPC_VERSION ?= v1.1.0

GO_GET_PKGS := $(GO_GET_PKGS) \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

PROTOC_GEN_GO_GRPC := $(CACHE_VERSIONS)/protoc-gen-go-grpc/$(PROTOC_GEN_GO_GRPC_VERSION)
$(PROTOC_GEN_GO_GRPC):
	@rm -f $(CACHE_BIN)/protoc-gen-go-grpc
	GOBIN=$(CACHE_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@rm -rf $(dir $(PROTOC_GEN_GO_GRPC))
	@mkdir -p $(dir $(PROTOC_GEN_GO_GRPC))
	@touch $(PROTOC_GEN_GO_GRPC)

dockerdeps:: $(PROTOC_GEN_GO_GRPC)
