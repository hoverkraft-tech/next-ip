FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY nextip.go ./
COPY cmd/next-ip ./cmd/next-ip
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/next-ip ./cmd/next-ip

FROM scratch
COPY --from=builder /out/next-ip /next-ip
ENTRYPOINT ["/next-ip"]
