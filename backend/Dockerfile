FROM golang:1.27.1-trixie AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
COPY vendor ./vendor

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /bin/main ./cmd/main.go

FROM alpine:3.24.2

WORKDIR /bin

COPY --from=builder /bin/main .
COPY config.yaml /etc/config.yaml

EXPOSE 8080

ENTRYPOINT [ "./main", "-c", "/etc/config.yaml" ]
