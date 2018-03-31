.PHONY: $(shell ls -d *)

default:
	@echo "Usage: make [command]"

lint:
	@command -v gometalinter || (go get -u github.com/alecthomas/gometalinter && gometalinter --install)
	@gometalinter .

test:
	@go test -v -race --cover . && echo ""

test-coverage:
	@go test -coverprofile=coverage.out . && go tool cover -html=coverage.out
