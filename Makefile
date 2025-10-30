VERSION =  $(shell git describe --tags --exact-match 2>/dev/null || git branch  --show-current )
REVISION = $(shell git rev-parse HEAD)
BRANCH = $(shell git branch  --show-current)
COMPILE_TIME= $(shell date +"%Y-%m-%d %H:%M:%S")
USER = $(shell  git log -1 --pretty=format:"%an")

FLAGS = -ldflags "-extldflags '-static' \
        -X 'github.com/prometheus/common/version.Version=${VERSION}' \
		-X 'github.com/prometheus/common/version.Revision=${REVISION}' \
		-X 'github.com/prometheus/common/version.Branch=${BRANCH}' -X \
		'github.com/prometheus/common/version.BuildUser=${USER}' \
		-X 'github.com/prometheus/common/version.BuildDate=${COMPILE_TIME}'"

.PHONY: build-gateway
build-gateway:
	CGO_ENABLED=0 go mod download && go build ${FLAGS} -o bin/gateway cmd/gateway/gateway.go

.PHONY: clean
clean:
	rm -rf bin
