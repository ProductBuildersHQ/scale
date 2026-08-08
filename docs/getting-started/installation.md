# Installation

SCALE ships as a Go library and a `scale` CLI.

## CLI

```bash
go install github.com/ProductBuildersHQ/scale/cmd/scale@latest
```

This installs the `scale` binary to `$(go env GOPATH)/bin`. Make sure that
directory is on your `PATH`.

Verify:

```bash
scale --help
```

## Library

```bash
go get github.com/ProductBuildersHQ/scale@latest
```

```go
import scale "github.com/ProductBuildersHQ/scale"
```

The embedded reference catalog is available without any files on disk:

```go
import "github.com/ProductBuildersHQ/scale/catalog"

f, err := catalog.Default()
```

## From source

```bash
git clone https://github.com/ProductBuildersHQ/scale
cd scale
go build ./...
go test ./...
```

## Requirements

- Go 1.26 or newer (see `go.mod` for the exact minimum).
- No runtime dependencies for report generation — the HTML report is
  self-contained with no external assets.
