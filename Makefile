APP="xdpdropper"

# Obtain an absolute path to the directory of the Makefile.
# Assume the Makefile is in the root of the repository.
REPODIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
UID := $(shell id -u)
GID := $(shell id -g)

# The development version of clang is distributed as the 'clang' binary,
# while stable/released versions have a version number attached.
# Pin the default clang to a stable version.
CLANG ?= clang-14
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)

# Assume docker as container engine
CONTAINER_ENGINE = docker
# Use cilium ebpf-builder built from upstream
IMAGE := quay.io/cilium/ebpf-builder
VERSION := 1648566014

.PHONY: all container-all container-shell generate build lint test

.DEFAULT_TARGET = container-all

# Build all ELF binaries using a containerized LLVM toolchain.
container-all:
	${CONTAINER_ENGINE} run --rm --user "${UID}:${GID}" \
		-v "${REPODIR}":/ebpf -w /ebpf --env MAKEFLAGS \
		--env CFLAGS="-fdebug-prefix-map=/ebpf=." \
		--env HOME="/tmp" \
		"${IMAGE}:${VERSION}" \
		make generate

# (debug) Drop the user into a shell inside the container as root.
container-shell:
	${CONTAINER_ENGINE} run --rm -ti \
		-v "${REPODIR}":/ebpf -w /ebpf \
		"${IMAGE}:${VERSION}"

# $BPF_CLANG is used in go:generate invocations.
generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate:
	cd pkg/xdp && go generate ./...

build:
	docker build -t ${APP}:build . -f ${REPODIR}/Dockerfile

lint:
	golangci-lint run

test:
	go test -cover ./...
