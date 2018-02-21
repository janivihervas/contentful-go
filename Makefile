GO_PACKAGES := `go list ./... | egrep -v "/vendor/"`
GO_FILES := `find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*"`
GO_FILES_NO_TEST := `find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -not -name "*_test.go"`
GO_TOOLS := golang.org/x/tools/cmd/goimports \
            github.com/golang/lint/golint \
            honnef.co/go/tools/staticcheck \
            github.com/fzipp/gocyclo \
            honnef.co/go/tools/cmd/unused \
            github.com/kisielk/errcheck \
            honnef.co/go/tools/cmd/gosimple \
            github.com/mdempsky/unconvert \
            github.com/alexkohler/nakedret \
            mvdan.cc/unparam

.PHONY: install
install:
	dep ensure -v

.PHONY: setup
setup:
	go get -u $(GO_TOOLS)

.PHONY: format
format:
	gofmt -s -w -e $(GO_FILES)
	goimports -w -l -e $(GO_FILES)

.PHONY: lint-sync
lint-sync:
	go vet ./...
	staticcheck $(GO_PACKAGES)
	golint -set_exit_status $(GO_PACKAGES)
	gocyclo -over 12 $(GO_FILES_NO_TEST)
	unused $(GO_PACKAGES)
	interfacer $(GO_PACKAGES)
	errcheck -ignoretests $(GO_PACKAGES)
	gosimple $(GO_PACKAGES)
	unconvert $(GO_PACKAGES)
	nakedret $(GO_PACKAGES)
	unparam $(GO_PACKAGES)

.PHONY: lint
lint:
	make --just-print lint-sync | parallel -k

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: test-ci
test-ci:
	touch coverage.txt
	rm -rf coverage
	mkdir coverage
	parallel -k 'go test -race -coverprofile=coverage/{#}.out -covermode=atomic {}' ::: $(GO_PACKAGES)
	find coverage/* -exec cat {} >> coverage.txt \;
	rm -rf coverage

