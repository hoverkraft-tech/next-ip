# Development

## Environment

Open this repository in Visual Studio Code and select "Reopen in Container". The devcontainer
installs Go tooling, Docker and the GitHub CLI, then runs `make setup`.
The committed feature lockfile makes the feature resolution reproducible.
Runtime and tool versions live in [devcontainer.json](../../.devcontainer/devcontainer.json).

For a host setup, install Docker and the same language tools, then run `make setup`.
Use the supported Go toolchain and golangci-lint declared in the devcontainer. `go.mod` declares the minimum language version.

## Commands

| Command         | Purpose                                           |
| --------------- | ------------------------------------------------- |
| `make setup`    | Install or verify dependencies                    |
| `make lint`     | Check application code and repository files       |
| `make lint-fix` | Apply supported formatting fixes                  |
| `make test`     | Run tests and collect coverage                    |
| `make audit`    | Check production dependencies for vulnerabilities |
| `make build`    | Build locally without publishing                  |
| `make ci`       | Run the complete local validation sequence        |

`make lint-repository` uses the same pinned Hoverkraft Super-Linter image and
validator selection as CI. It checks workflows, action pinning policy, configuration,
Markdown, spelling, shell scripts, Dockerfiles and infrastructure security.
Language lint runs through the project-owned tools in `make lint-code` and the shared
language CI, so dependency resolution and tool configuration stay consistent.

## Continuous integration

Pull requests, merge queues and main-branch pushes run the shared language checks
and repository linter. A weekly run detects newly disclosed dependency vulnerabilities.
Shared CI reports consolidate native reports after the producer jobs finish; producer
failures continue to block validation. The semantic PR workflow and repository linter
share [.commitlintrc.json](../../.commitlintrc.json).

The `ci-github-go` references intentionally use `main` during adoption.
The pinning policy permits that repository only; other external actions use reviewed
SHAs with version comments. Replace those temporary branch references with reviewed
release pins after adoption.

Dependabot groups weekly updates for actions, language dependencies, Docker images
and devcontainer features, with a seven-day cooldown and Conventional Commit titles.
Review the lockfile and run `make ci` when updating dependencies or features.

## Release archives

Go tests run on the supported CI matrix. GoReleaser then packages Linux and macOS
amd64/arm64 and Windows amd64 archives with CGO disabled and records checksums.
Manual or scheduled releases plan a `v`-prefixed tag, validate and package it, then
publish those retained archives without rebuilding during publication.

With GoReleaser v2 installed, validate packaging locally:

```sh
goreleaser release --clean --skip=publish --snapshot
```

## Container policy

The container runs as a non-root user. The scoped Docker health-check exception in
[the scanner configuration](../../.github/linters/trivy-exceptions.yaml) documents
why this application has no HTTP health endpoint.
