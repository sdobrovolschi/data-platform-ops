FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o db-cleaner ./cmd/db-cleaner

FROM ghcr.io/sdobrovolschi/containers/golang:1.26.0-wolfi
COPY --from=builder /app/db-cleaner .
ENTRYPOINT ["/db-cleaner"]
