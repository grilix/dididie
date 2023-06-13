FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/web ./cmd/web && go build -o /go/bin/api ./cmd/server

FROM debian:12-slim

COPY --from=builder /go/bin/web /go/bin/api /go/bin/
