.PHONY: $(shell ls -d *)

default:
	@echo "Usage: make [command]"

lint: basic.lint json.lint

%.lint:
	@command -v gometalinter || (go get -u github.com/alecthomas/gometalinter && gometalinter --install)
	@gometalinter ./$*

test: basic.test json.test

%.test:
	@go test -v -race --cover ./$* && echo ""

test-coverage: basic.cov json.cov

%.cov:
	@go test -coverprofile=$*/coverage.out ./$* && go tool cover -html=$*/coverage.out
