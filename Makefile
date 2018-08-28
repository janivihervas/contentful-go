SHELL=/bin/bash
MAKEFLAGS += --no-print-directory --output-sync
GO_FILES_NO_TEST := `find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -not -name "*_test.go"`
GO_TOOLS := golang.org/x/tools/cmd/goimports \
            github.com/golang/lint/golint \
            github.com/fzipp/gocyclo \
            github.com/kisielk/errcheck \
            github.com/alexkohler/nakedret \
            mvdan.cc/interfacer

.PHONY: install
install:
	go get ./...

.PHONY: setup
setup:
	go get -u $(GO_TOOLS)

.PHONY: format
format:
	gofmt -s -w -e ./...
	goimports -w -l -e .

.PHONY: vet golint gocyclo interfacer errcheck nakedret
vet:
	go vet ./...
golint:
	golint -set_exit_status ./...
gocyclo:
	gocyclo -over 12 $(GO_FILES_NO_TEST)
interfacer:
	interfacer ./...
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
	interfacer \
	errcheck \
	nakedret

.PHONY: test
test:
	go test -race -cover -run="^Test.*" ./...

.PHONY: test-all
test-all:
	go test -race -cover ./...

.PHONY: test-ci
test-codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	bash <(curl -s https://codecov.io/bash)

