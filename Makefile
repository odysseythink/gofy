SHELL := /bin/bash

ROOT      := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND   := $(ROOT)/backend
FRONTEND  := $(ROOT)/frontend
DEPLOY    := $(ROOT)/deploy
LINK_DIST := $(BACKEND)/main/link/frontend/dist

SERVICES  := link admin app tools plugins sandbox datasets

GO       ?= go
PNPM     ?= pnpm
GOFLAGS  ?=

.PHONY: all build frontend backend services deploy-dir clean clean-frontend clean-backend help $(SERVICES)

all: build

build: frontend backend ## Build frontend + backend (default)

## ---------- Frontend ----------

frontend: ## Build the Next.js frontend and copy dist into link service
	@echo "=== Building frontend ==="
	cd $(FRONTEND) && $(PNPM) install --frozen-lockfile && $(PNPM) build
	@echo "=== Copying frontend dist into link service ==="
	rm -rf $(LINK_DIST)
	mkdir -p $(LINK_DIST)
	cp -r $(FRONTEND)/dist/. $(LINK_DIST)/

## ---------- Backend ----------

backend: deploy-dir services ## Build all Go service binaries

services: $(SERVICES)

deploy-dir:
	@mkdir -p $(DEPLOY)

# link has extra frontend.go; others glob main/<svc>/*.go
link: deploy-dir
	@echo "=== Building link ==="
	cd $(BACKEND) && $(GO) build $(GOFLAGS) -o $(DEPLOY)/link ./main/link

$(filter-out link,$(SERVICES)): deploy-dir
	@echo "=== Building $@ ==="
	cd $(BACKEND) && $(GO) build $(GOFLAGS) -o $(DEPLOY)/$@ ./main/$@

## ---------- Clean ----------

clean: clean-backend clean-frontend ## Remove all build artifacts

clean-backend:
	@echo "=== Cleaning backend binaries ==="
	rm -f $(foreach s,$(SERVICES),$(DEPLOY)/$(s))

clean-frontend:
	@echo "=== Cleaning frontend artifacts ==="
	rm -rf $(FRONTEND)/dist $(LINK_DIST)

## ---------- Help ----------

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make              # build frontend + backend"
	@echo "  make frontend     # frontend only"
	@echo "  make backend      # all Go services"
	@echo "  make app          # single service"
	@echo "  make clean"
