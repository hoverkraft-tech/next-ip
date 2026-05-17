# next-ip

Small Go binary that returns the next IP address(es) from a CIDR input.

## Usage

```bash
go run ./cmd/next-ip --help

go run ./cmd/next-ip 192.168.100.12/24
# 192.168.100.13

go run ./cmd/next-ip --count 3 192.168.100.13/24
# 192.168.100.14
# 192.168.100.15
# 192.168.100.16

go run ./cmd/next-ip --step 2 192.168.100.100/24
# 192.168.100.102

go run ./cmd/next-ip --count 3 --step 3 192.168.100.102/24
# 192.168.100.105
# 192.168.100.108
# 192.168.100.111
```

Flags:
- `--count`, `-c`: number of next IP addresses to output
- `--step`, `-s`: increment step between emitted addresses

## Docker image

```bash
docker build -t next-ip .
docker run --rm next-ip 192.168.100.12/24
```

The final image is based on `scratch` and runs a static binary.

## CI (Dagger)

The CI workflow (`/.github/workflows/ci.yml`) runs Dagger pipelines for:

- lint (`go fmt` check + `go vet`)
- tests (`go test ./...`)
- multi-arch builds (`linux/amd64` and `linux/arm64`)

## Releases

Tagged releases (`v*`) are published with GoReleaser via `/.github/workflows/release.yml`.
