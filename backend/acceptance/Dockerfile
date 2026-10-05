FROM golang:1.27.1-trixie

WORKDIR /app

ENV CGO_ENABLED=0
ENV ACC_ENABLE=1

COPY go.mod .
COPY go.sum .
COPY vendor ./vendor

COPY migrations ./migrations
COPY internal ./internal

ENTRYPOINT ["go", "test", "./..."]
