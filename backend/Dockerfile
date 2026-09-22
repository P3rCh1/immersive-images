FROM golang:1.27.1-trixie AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /bin/main ./cmd/main.go

FROM alpine:3.24.2

WORKDIR /bin

COPY --from=builder /bin/main .
COPY config.yaml /etc/config.yaml

EXPOSE 8080

ENTRYPOINT [ "./main", "-c", "/etc/config.yaml" ]
