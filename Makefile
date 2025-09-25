# Generate protobuf files
.PHONY: config
config:
	@echo "Generating config protobuf files..."
	@buf generate proto/config
	@echo "Config protobuf files generated successfully"

# Generate all protobuf files
.PHONY: proto
proto:
	@echo "Generating all protobuf files..."
	@buf generate
	@echo "All protobuf files generated successfully"

# Install buf if not present
.PHONY: install-buf
install-buf:
	@which buf > /dev/null || (echo "Installing buf..." && \
		go install github.com/bufbuild/buf/cmd/buf@latest)

# Clean generated files
.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated protobuf files..."
	@find config -name "*.pb.go" -delete
	@echo "Cleaned generated protobuf files"

# make api (keeping original comment)
.PHONY: api
api: config