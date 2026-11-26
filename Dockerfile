FROM golang:1.27-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS builder
WORKDIR /src
COPY go.mod ./
COPY nextip.go ./
COPY cmd/next-ip ./cmd/next-ip
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/next-ip ./cmd/next-ip

FROM scratch
COPY --from=builder /out/next-ip /next-ip
USER 65532:65532
# A one-shot CLI exits after producing its output; it has no health endpoint.
HEALTHCHECK NONE
ENTRYPOINT ["/next-ip"]
