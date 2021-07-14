SHELL := /bin/bash
VERSION := $(shell cat .version)
SPEC_VERSION ?= $(shell cat .version)
NAME ?= metal-basecamp
BUILD_DIR ?= $(PWD)/dist/rpmbuild
SPEC_NAME ?= metal-basecamp
SPEC_FILE ?= ${SPEC_NAME}.spec
SOURCE_NAME ?= ${SPEC_NAME}-${SPEC_VERSION}
SOURCE_PATH := ${BUILD_DIR}/SOURCES/${SOURCE_NAME}.tar.bz2
BUILD_METADATA ?= "1~development~$(shell git rev-parse --short HEAD)"

.PHONY: \
	help \
	run \
	help \
	clean \
	clean-artifacts \
	clean-releases \
	tools \
	test \
	vet \
	lint \
	fmt \
	env \
	build \
	doc \
	version

all: fmt lint vet build
rpm: rpm_package_source rpm_build_source rpm_build

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    run                Run basecamp.'
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binaries, artifacts and releases.'
	@echo '    clean-artifacts    Remove build artifacts only.'
	@echo '    clean-releases     Remove releases only.'
	@echo '    tools              Install tools needed by the project.'
	@echo '    test               Run unit tests.'
	@echo '    vet                Run go vet.'
	@echo '    lint               Run golint.'
	@echo '    fmt                Run go fmt.'
	@echo '    env                Display Go environment.'
	@echo '    build              Build project for current platform.'
	@echo '    doc                Start Go documentation server on port 8080.'
	@echo '    version            Display Go version.'
	@echo ''
	@echo 'Targets run by default are: fmt, lint, vet, and build.'
	@echo ''

prepare:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/SPECS $(BUILD_DIR)/SOURCES
	cp $(SPEC_FILE) $(BUILD_DIR)/SPECS/

print-%:
	@echo $* = $($*)

clean: clean-artifacts clean-releases
	go clean -i ./...
	rm -vf \
	  $(CURDIR)/coverage.* \

clean-artifacts:
	rm -Rf artifacts/*

clean-releases:
	rm -Rf releases/*

clean-all: clean clean-artifacts

# Run tests
test: tools build
	go test ./cmd/... ./internal/...  -v -cover -coverprofile cover.out -covermode count

tools:
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/AlekSi/gocov-xml
	go get -u golang.org/x/lint/golint
	go get -u github.com/t-yuki/gocover-cobertura
	go get -u github.com/jstemmer/go-junit-report

vet:
	go vet -v ./...

lint: tools
	golint -set_exit_status  ./...

fmt:
	go fmt ./...

env:
	@go env

# Run against the configured Kubernetes cluster in ~/.kube/configs
run: build
	go run ./cmd/main.go(TARGET) .

build: fmt vet
	go build -o bin/basecamp ./cmd/main.go

doc:
	godoc -http=:8080 -index

version:
	@go version

rpm_package_source:
	tar --transform 'flags=r;s,^,/$(SOURCE_NAME)/,' --exclude .git --exclude dist -cvjf $(SOURCE_PATH) .

rpm_build_source:
	BUILD_METADATA=$(BUILD_METADATA) rpmbuild --nodeps -ts $(SOURCE_PATH) --define "_topdir $(BUILD_DIR)"

rpm_build:
	BUILD_METADATA=$(BUILD_METADATA) rpmbuild --nodeps -ba $(SPEC_FILE) --define "_topdir $(BUILD_DIR)"

# image:
# 	./runPostBuild.sh
