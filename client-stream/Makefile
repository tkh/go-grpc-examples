.PHONY: go-lint
go-lint:
	golangci-lint run

# Lint proto files using buf.
.PHONY: proto-lint
proto-lint:
	buf lint proto

# Rewrite proto files with autoformat.
.PHONY: proto-fmt
proto-fmt:
	buf format proto -w

# Generate code from proto files.
.PHONY: proto-gen
proto-gen:
	buf generate proto


# Clean any generated files (can be expanded later)
.PHONY: clean
clean:
	@echo "Cleaning..."
