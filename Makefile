SHELL=/bin/bash
MAKEFLAGS += --no-print-directory --output-sync
GO_FILES_NO_TEST := `find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -not -name "*_test.go"`
GO_TOOLS := golang.org/x/tools/cmd/goimports \
            github.com/golang/lint/golint \
            github.com/fzipp/gocyclo \
            github.com/kisielk/errcheck \
            github.com/alexkohler/nakedret

.PHONY: install
install:
	go get ./...

.PHONY: setup
setup:
	go get -u $(GO_TOOLS)

.PHONY: format
format:
	gofmt -s -w -e -l .
	goimports -w -e -l .

.PHONY: vet golint gocyclo errcheck nakedret
vet:
	go vet ./...
golint:
	golint -set_exit_status ./...
gocyclo:
	gocyclo -over 12 $(GO_FILES_NO_TEST)
errcheck:
	errcheck -ignoretests ./...
nakedret:
	nakedret ./...
.PHONY: lint
lint:
	@$(MAKE) -j \
	vet \
	golint \
	gocyclo \
	errcheck \
	nakedret

.PHONY: test
test:
	go test -race -cover -run="^Test.*" ./...

.PHONY: test-all
test-all:
	go test -race -cover ./...

.PHONY: test-codecov
test-codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	bash <(curl -s https://codecov.io/bash)

