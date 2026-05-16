FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY nextip.go ./
COPY cmd/get-next-ip ./cmd/get-next-ip
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/get-next-ip ./cmd/get-next-ip

FROM scratch
COPY --from=builder /out/get-next-ip /get-next-ip
ENTRYPOINT ["/get-next-ip"]
