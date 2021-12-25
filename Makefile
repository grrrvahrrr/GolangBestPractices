.PHONY: lint
lint:
	/home/deus/.go/bin/golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: build
build: lint
build: test
	go build
