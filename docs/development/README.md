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
version comments. Dependabot maintains these references.

Dependabot groups weekly updates for actions, language dependencies, Docker images
and devcontainer features, with a seven-day cooldown and Conventional Commit titles.
Renovate maintains runtime versions and embedded tool versions through the scheduled
[workflow](../../.github/workflows/__renovate.yml) and
[configuration](../../.github/renovate.json). It uses the organization CI bot's
`CI_BOT_APP_CLIENT_ID` variable and `CI_BOT_APP_PRIVATE_KEY` secret. Dependabot
continues to maintain package dependencies, actions, Docker images and features.
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

The container runs as a non-root user and declares `HEALTHCHECK NONE`. This
one-shot CLI exits after printing its result and has no health endpoint.
