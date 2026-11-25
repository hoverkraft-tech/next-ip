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
tool configuration as CI, without a repository validator allowlist. It checks
workflows, action pinning policy, configuration,
Markdown, spelling, shell scripts, Dockerfiles and infrastructure security.
Language lint runs through the project-owned tools in `make lint-code` and the shared
language CI, so dependency resolution and tool configuration stay consistent.
Super-Linter checks Go recursively through `GO_MODULES`; its standalone-file mode
cannot combine files from multiple packages and is disabled. `GOTOOLCHAIN=local`
keeps its bundled compiler compatible with its golangci-lint binary. The dedicated
Go CI still uses the project toolchain.

## Continuous integration

Pull requests, merge queues and main-branch pushes run the shared language checks
and repository linter. A weekly run detects newly disclosed dependency vulnerabilities.
Shared CI reports consolidate native reports after the producer jobs finish; producer
failures continue to block validation. The semantic PR workflow and repository linter
share [.commitlintrc.json](../../.commitlintrc.json).

External actions and reusable workflows use reviewed release SHAs with matching
version comments. Dependabot maintains these references. The pending shared workflow
capabilities use immutable candidate commits with linked TODOs; replace those pins
with the first containing release after the supporting PRs merge.

Dependabot is the primary updater for Go modules, GitHub Actions, Docker images
and devcontainer features, with weekly groups and a seven-day cooldown.
Renovate is limited to versions outside those supported manifests: Go and
golangci-lint feature options, and the govulncheck version embedded in the Makefile.
Renovate's Go module manager is disabled; Dependabot owns `go.mod`.
The scheduled [Renovate workflow](../../.github/workflows/__renovate.yml) uses
`CI_BOT_APP_CLIENT_ID` and `CI_BOT_APP_PRIVATE_KEY`; its
[configuration](../../.github/renovate.json) explicitly limits enabled managers.
Review the lockfile and run `make ci` when updating dependencies or features.

## Release archives

Go tests run on the supported CI matrix. The shared Go `build-release.yml`
workflow uses the project GoReleaser configuration to package Linux and macOS
amd64/arm64 and Windows amd64 archives with CGO disabled and records checksums.
Manual or scheduled releases plan a `v`-prefixed tag, validate and package it, then
publish those retained archives through the shared Go `release.yml` workflow.
Publication verifies the artifact belongs to the same run, commit and planned tag,
then attaches it without rebuilding. `CHANGELOG_SUMMARY` contains the optional
complete summary configuration; no provider-specific key is assembled in this repository.

With GoReleaser v2 installed, validate packaging locally:

```sh
goreleaser release --clean --skip=publish --snapshot
```

## Container policy

The container runs as a non-root user and declares `HEALTHCHECK NONE`. This
one-shot CLI exits after printing its result and has no health endpoint.
