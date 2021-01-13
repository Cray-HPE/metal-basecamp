# Build...
FROM        dtr.dev.cray.com/craypc/golang:1.15 as builder
# Copy the Go Modules manifests and all third-party libraries that are unlikely to change frequently
WORKDIR     /workspace
COPY        go.mod go.mod
COPY        go.sum go.sum
# Copy the go source...
COPY        configs configs/
COPY        cmd/ cmd/
COPY        internal/ internal/
RUN         CGO_ENABLED=0 \
            GOOS=linux \
            GOARCH=amd64 \
            GO111MODULE=on \
            go build -a -o basecamp ./cmd/main.go
# Run...
FROM        dtr.dev.cray.com/baseos/alpine:3.12.0
WORKDIR     /app
COPY        configs configs/
COPY        static/ static/
COPY        --from=builder /workspace/basecamp .
ENTRYPOINT  ["/app/basecamp"]
