FROM golang:latest

WORKDIR /go/src

#COPY application config cs-exporter ./
COPY . ./
RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN ls
RUN pwd
RUN go build -o /go/bin/main -ldflags '-s -w' ./cmd/main.go
RUN go install ./cmd/main.go

RUN go build -o /go/bin/health -ldflags '-s -w' ./cmd/health.go
RUN go install ./cmd/health.go

ENTRYPOINT ["/go/bin/main"]