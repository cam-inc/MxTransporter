FROM golang:latest

LABEL org.opencontainers.image.source="https://github.com/cam-inc/MxTransporter"

WORKDIR /go/src

COPY . ./
RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

RUN go build -o /go/bin/main -ldflags '-s -w' ./cmd/main.go
RUN go install ./cmd/main.go

RUN go build -o /go/bin/health -ldflags '-s -w' ./cmd/health.go
RUN go install ./cmd/health.go

ENTRYPOINT ["/go/bin/main"]