# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.17-alpine as builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api/ api/
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o app ./cmd/user/

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /
COPY configs/ configs/
COPY --from=builder /workspace/app .
EXPOSE 50051
ENTRYPOINT ["/app"]