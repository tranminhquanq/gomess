.PHONY: all build deps dev-deps image migrate test vet sec format unused
CHECK_FILES?=./...

FLAGS=-ldflags "-X github.com/tranminhquanq/gomess/internal/utils.Version=`git describe --tags`" -buildvcs=false
ifdef RELEASE_VERSION
	FLAGS=-ldflags "-X github.com/tranminhquanq/gomess/internal/utils.Version=v$(RELEASE_VERSION)" -buildvcs=false
endif

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: vet sec static build ## Run the tests and build the binary.

build: deps ## Build the binary.
	CGO_ENABLED=0 go build $(FLAGS)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(FLAGS) -o auth-arm64

deps: ## Install dependencies.
	@go mod download
	@go mod verify

format: ## Format the code.
	gofmt -s -w .