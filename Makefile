PKG := "github.com/cossim/hipush"
IMG ?= hub.hitosea.com/cossim/hipush:latest
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: dep
dep: ## Get the dependencies
	@go mod tidy

.PHONY: lint
lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: vet
vet: ## Run go vet
	@go vet ${PKG_LIST}

.PHONY: test
test: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: gen
gen: ## generate protobuf file
	@protoc -I=. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG},require_unimplemented_servers=true api/grpc/v1/*.proto
	@go fmt ./...
	@protoc-go-inject-tag -input=api/grpc/v1/*.pb.go

.PHONY: docker-build
# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
docker-build: dep test## Build docker image with the manager.
	docker build -t "${IMG}" .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}