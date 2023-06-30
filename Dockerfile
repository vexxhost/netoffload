FROM golang:1.20 AS builder
WORKDIR /app
COPY . /app
RUN \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
    go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:22.04
RUN \
  apt-get update && \
  apt-get install -y mstflint && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
