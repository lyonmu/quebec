VERSION =  $(shell git describe --tags --exact-match 2>/dev/null || git branch  --show-current )
REVISION = $(shell git rev-parse HEAD)
BRANCH = $(shell git branch  --show-current)
COMPILE_TIME= $(shell date +"%Y-%m-%d %H:%M:%S")
USER = $(shell  git log -1 --pretty=format:"%an")
FLAGS = -ldflags "-s -w \
	-X 'github.com/prometheus/common/version.Version=${VERSION}' \
	-X 'github.com/prometheus/common/version.Revision=${REVISION}' \
	-X 'github.com/prometheus/common/version.Branch=${BRANCH}' \
	-X 'github.com/prometheus/common/version.BuildUser=${USER}' \
	-X 'github.com/prometheus/common/version.BuildDate=${COMPILE_TIME}'"


.PHONY: ui
ui:
	cd web && bun install && bun run build

.PHONY: ui-dev
ui-dev:
	cd web && bun install && bun run dev

.PHONY: idl
idl:
	go mod download
	cd idl/node && go generate
	cd idl/router && go generate

.PHONY: gateway
gateway:
	make idl
	CGO_ENABLED=0 go build -tags netgo,osusergo ${FLAGS} -o bin/gateway cmd/gateway/gateway.go

.PHONY: core
core:
	make idl
	make ui
	CGO_ENABLED=0 go build -tags netgo,osusergo ${FLAGS} -o bin/core cmd/core/core.go

.PHONY: swag
swag:
	cd cmd/core && swag init -g core.go -o ./internal/docs --parseDependency --parseInternal

.PHONY: ent
ent:
	cd cmd/core/internal/ent \
	&& go run -mod=mod entgo.io/ent/cmd/ent --feature sql/upsert,sql/execquery,sql/modifier generate ./schema \
	&& go run -mod=mod entgo.io/ent/cmd/ent describe ./schema

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
