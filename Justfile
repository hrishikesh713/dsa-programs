# DSA Programs / LeetCode Practice Justfile

# Default recipe - show available commands
default:
    @just --list

# Run a specific problem with optional arguments
# Usage: just run <folder-name> [args...]
# Example: just run two-sum
# Example: just run two-sum --input=test.txt --verbose
run folder *args:
    @echo "Running {{folder}}..."
    go run ./{{folder}}/main.go {{args}}

# Build a specific problem
# Usage: just build <folder-name>
build folder:
    @echo "Building {{folder}}..."
    go build -o ./bin/{{folder}} ./{{folder}}/main.go

# Build all problems
build-all:
    @echo "Building all problems..."
    @for dir in $(find . -maxdepth 1 -type d -not -path '.' -not -path './bin' -not -path './internal' -not -path './.git*'); do \
        if [ -f "$$dir/main.go" ]; then \
            folder=$$(basename $$dir); \
            echo "Building $$folder..."; \
            mkdir -p ./bin; \
            go build -o ./bin/$$folder $$dir/main.go; \
        fi \
    done

# Run tests for a specific problem folder
# Usage: just test <folder-name>
test folder:
    @echo "Testing {{folder}}..."
    go test -v ./{{folder}}/...

# Run tests for a specific problem with coverage
# Usage: just test-coverage <folder-name>
test-coverage folder:
    @echo "Testing {{folder}} with coverage..."
    go test -v -cover -coverprofile=coverage.out ./{{folder}}/...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run all tests across all folders
test-all:
    @echo "Running all tests..."
    go test -v ./...

# Run all tests with coverage
test-all-coverage:
    @echo "Running all tests with coverage..."
    go test -v -cover -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run tests with race detector for a specific folder
# Usage: just test-race <folder-name>
test-race folder:
    @echo "Testing {{folder}} with race detector..."
    go test -race -v ./{{folder}}/...

# Run all tests with race detector
test-all-race:
    @echo "Running all tests with race detector..."
    go test -race -v ./...

# Run benchmarks for a specific problem
# Usage: just bench <folder-name>
bench folder:
    @echo "Running benchmarks for {{folder}}..."
    go test -bench=. -benchmem ./{{folder}}/...

# Run all benchmarks
bench-all:
    @echo "Running all benchmarks..."
    go test -bench=. -benchmem ./...

# Create a new problem folder with boilerplate
# Usage: just new <folder-name> [problem-number]
# Example: just new two-sum 1
# Example: just new valid-parentheses
new folder number="":
    @echo "Creating new problem: {{folder}}..."
    @mkdir -p {{folder}}
    @touch {{folder}}/README.md

# Format all Go code
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Run go vet on all code
vet:
    @echo "Running go vet..."
    go vet ./...

# Run golangci-lint (if installed)
lint:
    @echo "Running linter..."
    @if command -v golangci-lint >/dev/null 2>&1; then \
        golangci-lint run ./...; \
    else \
        echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
    fi

# Clean build artifacts and test cache
clean:
    @echo "Cleaning..."
    rm -rf bin/
    rm -f coverage.out coverage.html
    go clean -testcache

# List all folders
list:
    @echo "Available problems:"
    @for dir in $(find . -maxdepth 1 -type d -not -path '.' -not -path './bin' -not -path './internal' -not -path './.git*' | sort); do \
        if [ -f "$$dir/main.go" ]; then \
            folder=$$(basename $$dir); \
            echo "  - $$folder"; \
        fi \
    done

# Show information about a specific problem
# Usage: just info <folder-name>
info folder:
    @echo "Problem: {{folder}}"
    @if [ -f "./{{folder}}/README.md" ]; then \
        cat ./{{folder}}/README.md; \
    else \
        echo "No README.md found"; \
    fi

# Run go mod tidy
tidy:
    @echo "Running go mod tidy..."
    go mod tidy

# Update dependencies
update:
    @echo "Updating dependencies..."
    go get -u ./...
    go mod tidy

# Check for common issues (fmt, vet, test)
check: fmt vet test-all
    @echo "All checks passed!"

# Show statistics about the repository
stats:
    @echo "Repository Statistics:"
    @echo "====================="
    @echo "Total problems: $$(find . -maxdepth 1 -type d -not -path '.' -not -path './bin' -not -path './internal' -not -path './.git*' | wc -l)"
    @echo "Total Go files: $$(find . -name '*.go' -not -path './bin/*' | wc -l)"
    @echo "Total test files: $$(find . -name '*_test.go' | wc -l)"
    @echo "Lines of code: $$(find . -name '*.go' -not -path './bin/*' -exec cat {} \; | wc -l)"
