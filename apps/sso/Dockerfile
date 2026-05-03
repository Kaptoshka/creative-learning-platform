# ─── Stage 1: proto generation ───────────────────────────────────────────────
FROM golang:1.26-alpine as proto-gen

RUN apk add --no-cache \
    git=2.53.0-r0 \
    ca-certificates=20260413-r0 \
    curl=8.18.0-r0

RUN curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64 \
    -o /usr/local/bin/buf \
    && chmod +x /usr/local/bin/buf \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1

WORKDIR /workspace

COPY libs/protos/buf.yaml        libs/protos/buf.yaml
COPY libs/protos/buf.gen.yaml    libs/protos/buf.gen.yaml
COPY libs/protos/proto/          libs/protos/proto/
COPY libs/protos/go.mod          libs/protos/go.mod

RUN buf generate libs/protos --template libs/protos/buf.gen.yaml

# ─── Stage 2: builder ────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache \
    ca-certificates=20260413-r0 \
    tzdata=2026b-r0

WORKDIR /workspace

COPY go.work* ./
COPY libs/protos/go.mod libs/protos/go.mod
COPY --from=proto-gen /workspace/libs/protos/gen/ libs/protos/gen/

COPY apps/sso/go.mod apps/sso/go.mod ./apps/sso/
WORKDIR /workspace/apps/sso
RUN go mod download

WORKDIR /workspace
COPY apps/sso/ apps/sso/

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.28 && \
    wget -qO /bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /workspace/apps/sso
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -trimpath \
    -o /bin/sso ./cmd/sso \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -trimpath \
    -o /bin/migrator ./cmd/migrator

# ─── Stage 3: final ──────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/sso /bin/sso
COPY --from=builder /bin/migrator /bin/migrator
COPY --from=builder /workspace/apps/sso/migrations /migrations
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

EXPOSE 50051
EXPOSE 8080

ENTRYPOINT ["/bin/sso"]
