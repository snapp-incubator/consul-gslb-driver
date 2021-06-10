# Build the manager binary
FROM golang:1.16 AS builder

WORKDIR /workspace
ENV http_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128
ENV https_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o manager main.go


# Use alpine as a small base image to package the manager binary
FROM docker.io/library/alpine:latest
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
