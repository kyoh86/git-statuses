VERSION ?= `git vertag get`
COMMIT  ?= `git rev-parse HEAD`
DATE    ?= `date --iso-8601`

gen:
	go generate ./...
.PHONY: gen

lint: gen
	golangci-lint run
.PHONY: lint

test: lint
	go test -v --race ./...
.PHONY: test

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...
.PHONY: install

man:
	go run -tags man -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/git-statuses man
.PHONY: man

