GIT_COMMIT=$(shell git rev-list -1 HEAD)
LDFLAGS=-ldflags "-X CourseWork/request.GitCommit=${GIT_COMMIT}"

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: build
build: lint
build: test
	go build ${LDFLAGS}
