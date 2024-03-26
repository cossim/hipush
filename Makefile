PKG := "github.com/cossim/hipush"
IMG ?= hub.hitosea.com/cossim/hipush:latest
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
#PLATFORM ?= linux/amd64  # 默认平台参数为 linux/amd64

.PHONY: dep
dep: ## Get the dependencies
	@go mod tidy

.PHONY: lint
lint: ## Lint Golang files
	@golint -set_exit_status ./...

.PHONY: vet
vet: ## Run go vet
	@go vet ./...

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: test
test: fmt vet ## Run unittests
	@go test -short ./...

install: ## Install dependencies and protoc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: gen
gen: ## generate protobuf file
	#protoc -I=. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG},require_unimplemented_servers=true api/pb/v1/*.proto
	#protoc -I api/pb/v1 api/pb/v1/*.proto --gofast_out=api/pb/v1 --go-grpc_out=require_unimplemented_servers=false:api/pb/v1
	protoc -I api/pb/v1 api/pb/v1/*.proto --go_out=api/pb/v1 --go-grpc_out=require_unimplemented_servers=false:api/pb/v1
	protoc-go-inject-tag -input=api/pb/v1/*.pb.go

.PHONY: docker-build
# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
docker-build: dep test ## Build docker image with the manager.
	docker build -t "${IMG}" .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}
