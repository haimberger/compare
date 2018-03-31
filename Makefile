.PHONY: $(shell ls -d *)

default:
	@echo "Usage: make [command]"

precommit: dep lint test test-coverage

dep:
	@echo Syncing dependencies...
	@dep ensure
	@echo ""

lint:
	@echo Running linter...
	@command -v gometalinter > /dev/null || (go get -u github.com/alecthomas/gometalinter && gometalinter --install)
	@gometalinter .
	@echo ""

test:
	@echo Running tests...
	@go test -v -race --cover .
	@echo ""

test-coverage:
	@echo Calculating test coverage...
	@go test -coverprofile=coverage.out . && go tool cover -html=coverage.out
	@echo ""
