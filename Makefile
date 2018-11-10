SHELL=/bin/bash
MAKEFLAGS += --no-print-directory --output-sync
GO_FILES_NO_TEST := `find . -name "*.go" -not -name "*_test.go"`
GO_TOOLS := golang.org/x/lint/golint \
			golang.org/x/tools/cmd/goimports \
			\
			github.com/alexkohler/nakedret \
			github.com/fzipp/gocyclo \
			github.com/kisielk/errcheck \
			github.com/mdempsky/unconvert \
			\
			gitlab.com/opennota/check/cmd/structcheck \
			gitlab.com/opennota/check/cmd/varcheck \
			\
			honnef.co/go/tools/cmd/gosimple \
			honnef.co/go/tools/cmd/staticcheck \
			honnef.co/go/tools/cmd/unused

.PHONY: all
all: format lint test

.PHONY: install update
install:
	go get ./...
	go get $(GO_TOOLS)
update:
	go get -u ./...
	go mod tidy -v
	go get -u $(GO_TOOLS)

.PHONY: format
format:
	gofmt -s -w -e -l .
	goimports -w -e -l .

.PHONY: vet golint
vet:
	go vet ./...
golint:
	golint -set_exit_status ./...

.PHONY: nakedret gocyclo errcheck unconvert
nakedret:
	nakedret ./...
gocyclo:
	gocyclo -over 13 $(GO_FILES_NO_TEST)
errcheck:
	errcheck -ignoretests ./...
unconvert:
	unconvert ./...

.PHONY: structcheck varcheck
structcheck:
	structcheck ./...
varcheck:
	varcheck ./...

.PHONY: gosimple staticcheck unused
gosimple:
	gosimple ./...
staticcheck:
	staticcheck ./...
unused:
	unused ./...

.PHONY: lint
lint:
	@$(MAKE) -j \
	vet \
	golint \
	\
	nakedret \
	gocyclo \
	errcheck \
	unconvert \
	\
	structcheck \
	varcheck \
	\
	gosimple \
	staticcheck \
	unused


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
