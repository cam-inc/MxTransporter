##
## Build
##
FROM golang:latest as build

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

##
## Deploy
##
FROM alpine:latest

WORKDIR /go/src

COPY --from=build /go/bin/main /go/bin/main
COPY --from=build /go/bin/health /go/bin/health
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

ENTRYPOINT ["/go/bin/main"]