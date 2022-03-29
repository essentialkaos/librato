################################################################################

# This Makefile generated by GoMakeGen 1.5.1 using next command:
# gomakegen --mod .
#
# More info: https://kaos.sh/gomakegen

################################################################################

export GO111MODULE=on

.DEFAULT_GOAL := help
.PHONY = fmt vet deps mod-init mod-update mod-vendor help

################################################################################

deps: mod-update ## Download dependencies

mod-init: ## Initialize new module
	go mod init
	go mod tidy

mod-update: ## Download modules to local cache
	go mod download

mod-vendor: ## Make vendored copy of dependencies
	go mod vendor

fmt: ## Format source code with gofmt
	find . -name "*.go" -exec gofmt -s -w {} \;

vet: ## Runs go vet over sources
	go vet -composites=false -printfuncs=LPrintf,TLPrintf,TPrintf,log.Debug,log.Info,log.Warn,log.Error,log.Critical,log.Print ./...

help: ## Show this info
	@echo -e '\n\033[1mSupported targets:\033[0m\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-12s\033[0m %s\n", $$1, $$2}'
	@echo -e ''
	@echo -e '\033[90mGenerated by GoMakeGen 1.5.1\033[0m\n'

################################################################################
