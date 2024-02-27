PKG := "github.com/cossim/hipush"

dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}


gen: ## generate protobuf file
	@protoc -I=. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG},require_unimplemented_servers=true api/grpc/v1/*.proto
	@go fmt ./...
	@protoc-go-inject-tag -input=api/grpc/v1/*.pb.go