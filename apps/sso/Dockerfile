# ─── Stage 1: proto generation ───────────────────────────────────────────────
FROM golang:1.26-alpine as proto-gen

RUN apk add --no-cache git ca-certificates curl

RUN curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64 \
    -o /usr/local/bin/buf \
    && chmod +x /usr/local/bin/buf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /workspace

COPY libs/protos/buf.yaml        libs/protos/buf.yaml
COPY libs/protos/buf.gen.yaml    libs/protos/buf.gen.yaml
COPY libs/protos/proto/          libs/protos/proto/
COPY libs/protos/go.mod          libs/protos/go.mod

RUN buf generate libs/protos --template libs/protos/buf.gen.yaml

# ─── Stage 2: builder ────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /workspace

COPY go.work go.work.sum ./

COPY libs/protos/go.mod libs/protos/go.mod
COPY --from=proto-gen /workspace/libs/protos/gen/ libs/protos/gen/

COPY apps/sso/go.mod apps/sso/go.mod ./apps/sso/
RUN cd apps/assignments && go mod download

COPY apps/sso/ apps/sso/

RUN cd apps/sso && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" -trimpath \
    -o /bin/sso ./cmd/sso \
    && go build \
    -ldflags="-s -w" -trimpath \
    -o /bin/migrator ./cmd/migrator

# ─── Stage 3: final ──────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/sso /bin/sso
COPY --from=builder /bin/migrator /bin/migrator
COPY --from=builder /workspace/apps/sso/migrations /migrations

EXPOSE 50052

ENTRYPOINT ["/bin/sso"]
