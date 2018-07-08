NAME := golab
PKG := github.com/dimitrovvlado/$(NAME)

# Set any default go build tags
BUILDTAGS :=

# Populate version
VERSION := $(shell cat VERSION)
GITCOMMIT := $(shell git rev-parse --short HEAD)
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

# Set our default go compiler
GO := go
BINDIR := $(CURDIR)/bin

.PHONY: build
build: $(NAME) ## Builds a dynamic executable or package

$(NAME): $(wildcard *.go) $(wildcard */*.go) VERSION
	@echo "+ $@"
	$(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) .

.PHONY: install
install: ## Installs the executable or package
	@echo "+ $@"
	GOBIN=$(BINDIR) $(GO) install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} .