# Managed by makego. DO NOT EDIT.

# Must be set
$(call _assert_var,MAKEGO)
$(call _conditional_include,$(MAKEGO)/base.mk)
$(call _assert_var,CACHE_VERSIONS)
$(call _assert_var,CACHE_BIN)

GO_RELEASER_VERSION ?= latest

GO_RELEASER := $(CACHE_VERSIONS)/goreleaser/$(GO_RELEASER_VERSION)
$(GO_RELEASER):
	@rm -f $(CACHE_BIN)/goreleaser
	GOBIN=$(CACHE_BIN) go install github.com/goreleaser/goreleaser@$(GO_RELEASER_VERSION)
	@rm -rf $(dir $(GO_RELEASER))
	@mkdir -p $(dir $(GO_RELEASER))
	@touch $(GO_RELEASER)

dockerdeps:: $(GO_RELEASER)
