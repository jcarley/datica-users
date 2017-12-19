# Metadata about this makefile and position
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(dir $(realpath $(MKFILE_PATH)))
CURRENT_DIR := $(CURRENT_DIR:/=)

# Get the project metadata
GOVERSION := 1.9.2-stretch
VERSION := 0.0.1
PROJECT := github.com/jcarley/datica-users
OWNER := $(dir $(PROJECT))
OWNER := $(notdir $(OWNER:/=))
NAME := $(notdir $(PROJECT))
EXTERNAL_TOOLS =

# Current system information (this is the invoking system)
ME_OS = $(shell go env GOOS)
ME_ARCH = $(shell go env GOARCH)

# Default os-arch combination to build
XC_OS ?= darwin freebsd linux netbsd openbsd solaris windows
XC_ARCH ?= 386 amd64 arm
XC_EXCLUDE ?= darwin/arm solaris/386 solaris/arm windows/arm

# GPG Signing key (blank by default, means no GPG signing)
GPG_KEY ?=

# List of tests to run
TEST ?= ./...

# List all our actual files, excluding vendor
GOFILES = $(shell go list $(TEST) | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

# bootstrap installs the necessary go tools for development or build
bootstrap:
	@echo "==> Bootstrapping ${PROJECT}..."
	@for t in ${EXTERNAL_TOOLS}; do \
		echo "--> Installing $$t" ; \
		go get -u "$$t"; \
	done

# bin builds the project by invoking the compile script inside of a Docker
# container. Invokers can override the target OS or architecture using
# environment variables.
bin:
	@echo "==> Building ${PROJECT}..."
	@docker run \
		--interactive \
		--tty \
		--rm \
		--dns=8.8.8.8 \
		--env="VERSION=${VERSION}" \
		--env="PROJECT=${PROJECT}" \
		--env="OWNER=${OWNER}" \
		--env="NAME=${NAME}" \
		--env="GOMAXPROCS=${GOMAXPROCS}" \
		--env="GOTAGS=${GOTAGS}" \
		--env="XC_OS=${XC_OS}" \
		--env="XC_ARCH=${XC_ARCH}" \
		--env="XC_EXCLUDE=${XC_EXCLUDE}" \
		--env="DIST=${DIST}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		"golang:${GOVERSION}" /usr/bin/env sh -c "scripts/compile.sh"

# deps gets all the dependencies for this repository and vendors them.
deps:
	@echo "==> Updating dependencies for ${CURRENT_DIR}..."
	@docker run \
		--interactive \
		--tty \
		--rm \
		--dns=8.8.8.8 \
		--env="GOMAXPROCS=${GOMAXPROCS}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		"golang:${GOVERSION}" /usr/bin/env sh -c "scripts/deps.sh"

shell:
	@echo "==> Starting up a shell for ${CURRENT_DIR}..."
	@docker-compose run \
		--rm \
		-e "GOMAXPROCS=${GOMAXPROCS}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		shell /bin/bash

database:
	@echo "==> Migrating database for ${CURRENT_DIR}..."
	@docker-compose run \
		--rm \
		-T \
		-e "GOMAXPROCS=${GOMAXPROCS}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		shell ./scripts/build_database.sh

package:
	@echo "==> Packaging into container for ${CURRENT_DIR}..."
	@docker-compose build package
	@echo "==> Cleaning up dangling images ..."
	@docker images --quiet --filter=dangling=true | xargs docker rmi

run:
	@echo "==> Running ${PROJECT} ..."
	@docker-compose run --rm -p 3000:3000 package /datica-users server

# test runs the test suite
test:
	@echo "==> Testing ${PROJECT}..."
	@docker run \
		--interactive \
		--tty \
		--rm \
		--dns=8.8.8.8 \
		--env="GOMAXPROCS=${GOMAXPROCS}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}:delegated" \
		"golang:${GOVERSION}" go test -cover -timeout=60s -parallel=10 -tags="${GOTAGS}" ${GOFILES} ${TESTARGS}

convey:
	@echo "==> Testing ${PROJECT} with goconvey..."
	@docker run \
		--interactive \
		--tty \
		--rm \
		-p 0.0.0.0:8080:8080 \
		--dns=8.8.8.8 \
		--env="GOMAXPROCS=${GOMAXPROCS}" \
		--workdir="/go/src/${PROJECT}" \
		--volume="${CURRENT_DIR}:/go/src/${PROJECT}" \
		"golang:${GOVERSION}" /usr/bin/env sh -c "scripts/testgoconvey.sh"

db:
	@echo "==> Starting up database..."
	@docker-compose up -d db

.PHONY: package database shell bin bin-local bootstrap deps dev dist docker docker-push generate test test-race convey

