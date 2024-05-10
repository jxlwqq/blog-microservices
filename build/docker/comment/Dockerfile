# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.21-alpine as builder
ARG TARGETARCH
ARG TARGETOS
# 启用 Go Modules 功能
RUN go env -w GO111MODULE=on
# go 设置国内源
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go install github.com/grpc-ecosystem/grpc-health-probe@latest
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api/ api/
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o app ./cmd/comment/

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /
COPY configs/ configs/
COPY --from=builder /workspace/app .
COPY --from=builder /go/bin/grpc-health-probe .
EXPOSE 50053
EXPOSE 9053
ENTRYPOINT ["/app"]