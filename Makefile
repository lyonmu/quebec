VERSION =  $(shell git describe --tags --exact-match 2>/dev/null || git branch  --show-current )
REVISION = $(shell git rev-parse HEAD)
BRANCH = $(shell git branch  --show-current)
COMPILE_TIME= $(shell date +"%Y-%m-%d %H:%M:%S")
USER = $(shell  git log -1 --pretty=format:"%an")

# macOS 不支持完全静态链接，只在 Linux 上使用 -static
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	FLAGS = -ldflags "-extldflags '-static' \
        -X 'github.com/prometheus/common/version.Version=${VERSION}' \
		-X 'github.com/prometheus/common/version.Revision=${REVISION}' \
		-X 'github.com/prometheus/common/version.Branch=${BRANCH}' -X \
		'github.com/prometheus/common/version.BuildUser=${USER}' \
		-X 'github.com/prometheus/common/version.BuildDate=${COMPILE_TIME}'"
else
	FLAGS = -ldflags "-X 'github.com/prometheus/common/version.Version=${VERSION}' \
		-X 'github.com/prometheus/common/version.Revision=${REVISION}' \
		-X 'github.com/prometheus/common/version.Branch=${BRANCH}' -X \
		'github.com/prometheus/common/version.BuildUser=${USER}' \
		-X 'github.com/prometheus/common/version.BuildDate=${COMPILE_TIME}'"
endif

.PHONY: ui
ui:
	cd web && bun install && bun run build

.PHONY: ui-dev
ui-dev:
	cd web && bun install && bun run dev

.PHONY: idl
idl:
	cd idl/node && go generate
	cd idl/router && go generate

.PHONY: gateway
gateway:
	make idl
	CGO_ENABLED=0 go mod download && go build ${FLAGS} -o bin/gateway cmd/gateway/gateway.go

.PHONY: core
core:
	make idl
	make ui
	CGO_ENABLED=0 go mod download && go build ${FLAGS} -o bin/core cmd/core/core.go

.PHONY: all
all: gateway core ui

.PHONY: docker-builder
docker-builder:
	docker build -t lyonmu/quebec:builder-bookworm -f Dockerfile_builder . 

.PHONY: docker-all
docker-all:
	docker build -t lyonmu/quebec:core-lts -f Dockerfile_core .
	docker build -t lyonmu/quebec:gateway-lts -f Dockerfile_gateway .

.PHONY: clean
clean:
	rm -rf bin
	rm -rf idl/node/v1/*.pb.go
	rm -rf idl/router/v1/*.pb.go
