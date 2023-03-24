OUT_DIR := ./out
GO_FILES := $(shell find . -type f \( -iname '*.go' \))

.PHONY: build
build:
	mkdir -p $(OUT_DIR)
	go build -ldflags="-w -s" -o $(OUT_DIR)/cheergo ./cmd/cheergo

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: test
test:
	stdout=$$(gofumpt -l . 2>&1); if [ "$$stdout" ]; then exit 1; fi
	go vet ./...
	gocyclo -over 10 $(GO_FILES)
	staticcheck ./...
	errcheck ./...
	go test -v -cover ./...
	@printf '\n%s\n' "> Test successful"

.PHONY: setup
setup:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/kisielk/errcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: fmt
fmt:
	gofumpt -l -w .
