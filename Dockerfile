# Build the manager binary
FROM golang:1.16 AS builder

WORKDIR /workspace
ENV http_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128
ENV https_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128
ENV no_proxy=localhost,127.0.0.1,gitlab.snapp.ir,mirror-teh-1.snappcloud.io,mirror-fra-1.snappcloud.io
ENV GOPRIVATE="gitlab.snapp.ir/snappcloud"
ENV personal_access_token="dummy"
RUN git config \
  --global \
  url."https://oauth2:${personal_access_token}@gitlab.snapp.ir".insteadOf \
  "https://gitlab.snapp.ir"
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o server cmd/server/main.go


# Use alpine as a small base image to package the server binary
FROM docker.io/library/alpine:latest
WORKDIR /
COPY --from=builder /workspace/server .
USER 65532:65532

ENTRYPOINT ["/server"]
