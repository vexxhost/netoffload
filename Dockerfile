FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:22.04
RUN \
  apt-get update && \
  apt-get install -y mstflint && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
