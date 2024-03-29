.PHONY: test
test:
	go test -race -v -coverprofile=cover.out  ./...

.PHONY: test-ui
test-ui: test
	go tool cover -html=cover.out -o cover.html
	open cover.html

.PHONY: fmt
fmt:
	gofmt -w $(shell find . -name "*.go")

.PHONY: lint
lint:
	golangci-lint run
