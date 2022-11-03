# Managed by makego. DO NOT EDIT.

# Must be set
$(call _assert_var,MAKEGO)
$(call _conditional_include,$(MAKEGO)/base.mk)
$(call _assert_var,CACHE_VERSIONS)
$(call _assert_var,CACHE_BIN)

GO_COVER_VIEW_VERSION ?= latest
GO_JUNIT_REPORT_VERSION ?= latest

GO_COVER_VIEW := $(CACHE_VERSIONS)/go-cover-view/$(GO_COVER_VIEW_VERSION)
GO_JUNIT_REPORT := $(CACHE_VERSIONS)/go-junit-report/$(GO_JUNIT_REPORT_VERSION)
$(GO_COVER_VIEW):
	@rm -f $(CACHE_BIN)/go-cover-view
	@rm -f $(CACHE_BIN)/go-junit-report
	GOBIN=$(CACHE_BIN) go install github.com/johejo/go-cover-view@$(GO_COVER_VIEW_VERSION)
	GOBIN=$(CACHE_BIN) go install github.com/jstemmer/go-junit-report/v2@$(GO_JUNIT_REPORT_VERSION)
	@rm -rf $(dir $(GO_COVER_VIEW))
	@rm -rf $(dir $(GO_JUNIT_REPORT))
	@mkdir -p $(dir $(GO_COVER_VIEW))
	@mkdir -p $(dir $(GO_JUNIT_REPORT))
	@touch $(GO_COVER_VIEW)
	@touch $(GO_JUNIT_REPORT)

dockerdeps:: $(GO_COVER_VIEW) $(GO_JUNIT_REPORT)
