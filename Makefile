APPNAME=$(shell basename $(shell go list))
VERSION?=snapshot
COMMIT=$(shell git rev-parse --verify HEAD)
DATE?=$(shell date +%FT%T%z)
RELEASE?=0

GO_LDFLAGS+=-X main.appName=$(APPNAME)
GO_LDFLAGS+=-X main.buildVersion=$(VERSION)
GO_LDFLAGS+=-X main.buildCommit=$(COMMIT)
GO_LDFLAGS+=-X main.buildDate=$(DATE)
ifeq ($(RELEASE), 1)
	# Strip debug information from the binary
	GO_LDFLAGS+=-s -w
endif
GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS)"

LEVEL=debug

REFLEX=$(GOPATH)/bin/reflex
$(REFLEX):
	go install github.com/cespare/reflex@latest

GOLANGCILINTVERSION:=v1.44.2
GOLANGCILINT=$(GOPATH)/bin/golangci-lint
$(GOLANGCILINT):
	curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCILINTVERSION)

VENOMVERSION:=v1.0.1
VENOM=$(GOPATH)/bin/venom
$(VENOM):
	go install github.com/ovh/venom/cmd/venom@$(VENOMVERSION)

CADDYVERSION:=v2.4.6
CADDY=$(GOPATH)/bin/caddy
$(CADDY):
	go install github.com/caddyserver/caddy@$(CADDYVERSION)

.PHONY: default
default: start

.PHONY: start
start: $(REFLEX)
	$(REFLEX) --start-service \
		--decoration='none' \
		--regex='\.go$$' \
		--inverse-regex='^vendor|node_modules|.cache/' \
		-- go run $(GO_LDFLAGS) main.go --log-level=$(LEVEL)