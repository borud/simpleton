GOOS ?= darwin
GOARCH ?= amd64
CGO_ENABLED ?= 1
CGO_CFLAGS ?=
CGO_LDFLAGS ?=
BUILD_TAGS ?=
APP_NAME ?= simpleton
VERSION ?=
BIN_EXT ?=

GO := GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) CGO_CFLAGS=$(CGO_CFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) GO111MODULE=auto go
PACKAGES = $(shell $(GO) list ./... | grep -v '/vendor/')
PROTOBUFS = $(shell find . -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq | grep -v /vendor/)
TARGET_PACKAGES = $(shell find . -name 'main.go' -print0 | xargs -0 -n1 dirname | sort | uniq | grep -v /vendor/)

ifeq ($(GOOS),windows)
  BIN_EXT = .exe
endif

ifeq ($(VERSION),)
  VERSION = latest
endif

.DEFAULT_GOAL := build

.PHONY: help
help:
	@echo "   GOOS        = $(GOOS)"
	@echo "   GOARCH      = $(GOARCH)"
	@echo "   CGO_ENABLED = $(CGO_ENABLED)"
	@echo "   CGO_CFLAGS  = $(CGO_CFLAGS)"
	@echo "   CGO_LDFLAGS = $(CGO_LDFLAGS)"
	@echo "   BUILD_TAGS  = $(BUILD_TAGS)"
	@echo "   VERSION     = $(VERSION)"


.PHONY: protoc
protoc:
	@for proto_dir in $(PROTOBUFS); do echo $$proto_dir; protoc --proto_path=. --proto_path=$$proto_dir --go_out=plugins=grpc:$(GOPATH)/src $$proto_dir/*.proto || exit 1; done

.PHONY: format
format:
	@$(GO) fmt $(PACKAGES)

.PHONY: test
test:
	@$(GO) test -v -tags="$(BUILD_TAGS)" $(PACKAGES)

.PHONY: build
build:
	@for target_pkg in $(TARGET_PACKAGES); do $(GO) build -tags="$(BUILD_TAGS)" $(LDFLAGS) -o ./bin/`basename $$target_pkg`$(BIN_EXT) $$target_pkg || exit 1; done

.PHONY: install
install:
	@for target_pkg in $(TARGET_PACKAGES); do $(GO) install -tags="$(BUILD_TAGS)" $(LDFLAGS) $$target_pkg || exit 1; done

.PHONY: dist
dist: build
	@mkdir -p ./dist/$(GOOS)-$(GOARCH)/bin
	@(cd ./dist/$(GOOS)-$(GOARCH); tar cfz ../$(APP_NAME)-${VERSION}.$(GOOS)-$(GOARCH).tar.gz .)

.PHONY: git-tag
git-tag:
ifeq ($(VERSION),$(filter $(VERSION),latest master ""))
	@echo "please specify VERSION"
else
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
endif

.PHONY: clean
clean:
	@rm -rf ./bin
	@rm -rf ./data
	@rm -rf ./dist
