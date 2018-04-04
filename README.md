# compare

[![GoDoc](https://godoc.org/github.com/haimberger/compare?status.svg)](https://godoc.org/github.com/haimberger/compare)

Package `compare` provides customizable functionality for comparing values.

## Development

Before committing any changes, make sure to run `make precommit`. It does the following:

1. Make sure that `Gopkg.lock` and the `vendor` directory are up to date. For this, you'll need to install [dep](https://github.com/golang/dep) if you haven't already.
1. Concurrently run a bunch of linters including [go vet](https://golang.org/cmd/vet/) and [megacheck](https://github.com/dominikh/go-tools/tree/master/cmd/megacheck).
1. Run all tests.
1. Save test coverage information to a file, then open a browser window showing the covered (green), uncovered (red), and uninstrumented (grey) source. You can find more information under "Viewing the results" [here](https://blog.golang.org/cover).
