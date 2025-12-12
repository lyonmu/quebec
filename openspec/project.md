# Project Context

## Purpose
- Lightweight Envoy control plane to manage gateways without Kubernetes/Istio. Builds two Go binaries (`core`, `gateway`) plus a React/TypeScript admin UI.

## Tech Stack
- Backend: Go 1.24, gRPC, Gin HTTP, Ent ORM, MySQL (go-sql-driver/mysql), Redis.
- Proto/IDL: protoc + `protoc-gen-go` / `protoc-gen-go-grpc` (`idl/*/v1` with `go generate` via `make idl`).
- Frontend: React + TypeScript + Vite, bundled with Bun; UI build under `web/`.

## Project Conventions

### Code Style
- Go: `gofmt` (run before commit), keep generated files via `go generate` in `idl/*/v1`. Log via `global.Logger`.
- Frontend: standard Vite/React TS; prefer functional components; rely on Bun/Vite defaults (Prettier not yet enforced).

### Architecture Patterns
- Services: two executables `cmd/core` (control plane + admin API) and `cmd/gateway` (data plane sidecar/front proxy).
- Persistence: Ent ORM schemas in `cmd/core/internal/ent/schema`; DB is MySQL.
- gRPC/HTTP: IDLs in `idl/`; core registers gRPC services under `cmd/core/internal/service/grpc`; HTTP via Gin under `cmd/core/internal/api/http`.
- IDL generation: `make idl` regenerates Go stubs for node/router services.
- Frontend build: `make ui` (Bun install + build); `core` binary bundles built UI.

### Testing Strategy
- Go: `go test ./...` (currently few/no tests in many packages). Add unit tests for new logic when practical.
- Generated code not hand-edited; validate builds after `make idl`.

### Git Workflow
- Default branch `main` (assumed). Use feature branches + PRs; squash/merge acceptable. Keep commits small, with clear messages; no history rewrites on shared branches.

## Domain Context
- Manages gateway clusters and nodes; nodes register via gRPC streaming (`NodeService.NodeRegister`) and persist to `core_gateway_node` (Ent schema). Admin UI provides system management (users, roles, audit, proxy nodes, etc.).

## Important Constraints
- macOS cannot fully static link; Makefile guards Linux-only `-static`. Set `CGO_ENABLED=0` for releases.
- IDL outputs in `idl/*/v1/*.pb.go` are generated; do not hand-edit.

## External Dependencies
- Envoy data plane consuming control-plane configs (xDS-style).
- MySQL database for core persistence.
- Redis for caching/auth/session helpers.
