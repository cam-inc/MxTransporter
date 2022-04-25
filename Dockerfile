# ------------------------
# First stage: build
# ------------------------
FROM golang:latest as builder

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

# ------------------------
# Complete stage
# ------------------------
FROM golang:latest
COPY --from=builder /go /go

ENTRYPOINT ["/go/bin/main"]

