FROM golang:1.21 AS builder

WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/web ./cmd/web && go build -o /go/bin/api ./cmd/server

FROM debian:11-slim

COPY --from=builder /go/bin/web /go/bin/api /go/bin/
