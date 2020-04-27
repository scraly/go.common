GO   := go

DIRS_TO_CLEAN:=
FILES_TO_CLEAN:=

ifeq ($(origin GO), undefined)
  GO:=$(shell which go)
endif
ifeq ($(GO),)
  $(error Could not find 'go' in path.  Please install go, or if already installed either add it to your path or set GO to point to its directory)
endif


pkgs  = $(shell $(GO) list ./...)
pkgDirs = $(shell $(GO) list -f {{.Dir}} ./... | grep -vE /vendor/)

DIR_OUT:=/tmp

GOLANGCI:=$(shell command -v golangci-lint 2> /dev/null)
GOCOV:=$(shell command -v gocov 2> /dev/null)
WWHRD:=$(shell command -v wwhrd 2> /dev/null)
GOCLOC:=$(shell command -v gocloc 2> /dev/null)
DATABASES := bolt

GO_EXCLUDE := /vendor/|.pb.go|.gen.go
GO_FILES_CMD := find . -name '*.go' | grep -v -E '$(GO_EXCLUDE)'

NAME := datastore-adapter

#-------------------------
# Final targets
#-------------------------
.PHONY: dev

## Execute development pipeline
dev: license generate format lint.full test

#-------------------------
# Download libraries and tools
#-------------------------
.PHONY: get.tools

## Retrieve tools packages
get.tools:
  # Golang DEP
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v github.com/ugorji/go/codec/codecgen
	# License checker
	go get -u -v github.com/frapposelli/wwhrd
 	# linter
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	# test
	go get -u -v github.com/onsi/gomega/...
	go get -u -v github.com/smartystreets/goconvey
	go get -u -v github.com/vektra/mockery/.../
	go get -u -v github.com/jstemmer/go-junit-report
	# proto
	go get -u -v github.com/golang/protobuf/protoc-gen-go

#-------------------------
# Code generation
#-------------------------
.PHONY: generate

## Generate go code
generate:
	@echo "==> generating go code"
	@$(GO) generate $(pkgs)

#-------------------------
# Unit tests
#-------------------------
.PHONY: test test.report test.coverage

## Run tests
test:
	@echo "==> running tests"
	@GOCACHE=off $(GO) test -short -race -cover $(pkgs)

## Run tests report
test.report:
	@echo "==> running tests with coverage report"
	@env GOCACHE=off $(GO) test -short -race -coverprofile=coverage.out -v $(pkgs) > tests.txt

#-------------------------
# Checks
#-------------------------
.PHONY: format license license.csv lint.fast lint.full lint.sonar stats.loc

check: format license lint.full

## Apply code format, import reorganization and code simplification on source code
format:
	@echo "==> formatting code"
	@$(GO) fmt $(pkgs)
	@echo "==> clean imports"
	@goimports -w $(pkgDirs)
	@echo "==> simplify code"
	@gofmt -s -w $(pkgDirs)

## Check external license usage
license:
ifndef WWHRD
	$(error "Please install wwhrd! make get-tools")
endif
	@echo "==> license check"
	wwhrd check

## Export license usage as CSV file
license.csv:
ifndef WWHRD
	$(error "Please install wwhrd! make get-tools")
endif
	@echo "==> license check (csv)"
	wwhrd list 2>&1 | sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g" | cut -c 58- | awk -F"[ =]" '{print $2 ";" $4;}' | sort > licence.csv

## Launch linter
lint.fast:
ifndef GOLANGCI
	$(error "Please install golangci! make get-tools")
endif
	@echo "==> linters (fast)"
	golangci-lint run --fast

## Validate code
lint.full:
ifndef GOLANGCI
	$(error "Please install golangci! make get-tools")
endif
	@echo "==> linters (slow)"
	golangci-lint run

## Execute linter and save result file
lint.report:
ifndef GOLANGCI
	$(error "Please install golangci! make get-tools")
endif
	golangci-lint run ./... > golangci-report.txt

## Output project lines of code table
stats.loc:
ifndef GOCLOC
	$(error "Please install gocloc ! go get -u github.com/hhatto/gocloc/cmd/gocloc ")
endif
	@date
	@gocloc --not-match-d="vendor|node_modules" .

#-------------------------
# Target: depend
#-------------------------
.PHONY: depend vendor.check depend.status depend.update depend.cleanlock depend.update.full

## Run dep Ensure
depend: depend.update

## Test if dependencies are correctly set
depend-test: vendor.check depend.status

vendor.check:
	@echo "Checking that Gopkg.* are in sync with vendor/ submodule:"
	diff Gopkg.toml vendor/
	diff Gopkg.lock vendor/

depend.status:
	@echo "No error means your Gopkg.* are in sync and ok with vendor/"
	dep status -dot > $(DIR_OUT)/dep.dot
	cp Gopkg.* vendor/

depend.update:
	@echo "==> Running dep ensure"
	dep ensure -v
	cp Gopkg.* vendor/

depend.update.full: depend.cleanlock depend.update

## Remove dep lock file
depend.cleanlock:
	-rm Gopkg.lock

depend.vendor:
	@echo "==> Running dep ensure (vendor only)"
	dep ensure -v -vendor-only

#-------------------------
# Target: clean
#-------------------------

## Clean build files
clean: clean.go
	rm -rf $(DIRS_TO_CLEAN)
	rm -f $(FILES_TO_CLEAN)

clean.go: ; $(info cleaning...)
	$(eval GO_CLEAN_FLAGS := -i -r)
	$(GO) clean $(GO_CLEAN_FLAGS)

#-------------------------
# Target: help
#-------------------------

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=20
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
