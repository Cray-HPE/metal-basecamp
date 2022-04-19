#
# MIT License
#
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#
# Build...
FROM        arti.dev.cray.com/baseos-docker-master-local/golang:alpine as builder
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
FROM        arti.dev.cray.com/baseos-docker-master-local/alpine:3
WORKDIR     /app
COPY        configs configs/
COPY        static/ static/
COPY        --from=builder /workspace/basecamp .
ENTRYPOINT  ["/app/basecamp"]
