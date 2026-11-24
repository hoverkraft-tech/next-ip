# Development

Use Go 1.26 or 1.27 and golangci-lint v2.15.0.

```sh
gofmt -l .
go mod tidy -diff
golangci-lint run
go test -race -count=1 -timeout=5m -covermode=atomic -coverprofile=coverage.out ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.10.0 ./...
```

CI calls `ci-github-go` at `main` while the shared workflow is being validated
in consumer repositories. Tests run on both supported Go versions; lint, audit,
build and coverage reporting use Go 1.27. Dependency changes are reviewed on PRs.

After the Go checks pass, GoReleaser packages Linux and macOS amd64/arm64 and
Windows amd64 archives, with CGO disabled, and retains their checksums.
To validate snapshot packaging locally with GoReleaser v2:

```sh
goreleaser release --clean --skip=publish --snapshot
```

The scheduled or manually dispatched release workflow plans a `v`-prefixed tag,
validates and packages that version, then publishes those retained archives.
It does not rebuild during publication.
