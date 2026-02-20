.PHONY: gen help clean install-deps check-deps

# Default target
help:
	@echo "Available commands:"
	@echo "  make install-deps - Install protoc plugins"
	@echo "  make gen          - Compile all proto files"
	@echo "  make clean        - Clean generated pb.go files"

# Install dependencies
install-deps:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/syncore/protoc-go-inject-tag@latest
	@echo "Plugins installed successfully!"

# Check dependencies
check-deps:
	@which protoc > /dev/null || (echo "Error: protoc not found, please install Protocol Buffers" && exit 1)
	@which protoc-gen-go > /dev/null || (echo "Error: protoc-gen-go not found, run 'make install-deps'" && exit 1)
	@which protoc-gen-go-grpc > /dev/null || (echo "Error: protoc-gen-go-grpc not found, run 'make install-deps'" && exit 1)
	@which protoc-go-inject-tag > /dev/null || (echo "Error: protoc-go-inject-tag not found, run 'make install-deps'" && exit 1)

# Compile all proto files
gen: check-deps
	@echo "Compiling proto files..."
	@$(MAKE) gen-resource
	@$(MAKE) gen-secret
	@$(MAKE) gen-tags
	@echo "Proto compilation completed!"

# Compile resource module
# 1. protoc 生成 .pb.go
# 2. protoc-go-inject-tag 读取 // @gotags: 注释注入 struct tag
gen-resource:
	@if [ -f apps/resource/service.proto ]; then \
		echo "Compiling apps/resource/service.proto..."; \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			apps/resource/service.proto; \
		protoc-go-inject-tag -input=apps/resource/service.pb.go; \
	else \
		echo "apps/resource/service.proto not found"; \
	fi

# Compile secret module
gen-secret:
	@if [ -f apps/secret/service.proto ]; then \
		echo "Compiling apps/secret/service.proto..."; \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			apps/secret/service.proto; \
		protoc-go-inject-tag -input=apps/secret/service.pb.go; \
	else \
		echo "apps/secret/service.proto not found, skipping"; \
	fi

# Compile tags module
gen-tags:
	@if [ -f apps/tags/service.proto ]; then \
		echo "Compiling apps/tags/service.proto..."; \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			apps/tags/service.proto; \
		protoc-go-inject-tag -input=apps/tags/service.pb.go; \
	else \
		echo "apps/tags/service.proto not found, skipping"; \
	fi

# Clean generated files
clean:
	@echo "Cleaning generated pb.go files..."
	@find apps -name "*.pb.go" -type f -delete
	@echo "Clean completed!"