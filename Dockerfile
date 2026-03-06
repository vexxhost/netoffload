# SPDX-FileCopyrightText: © 2025 VEXXHOST, Inc.
# SPDX-License-Identifier: GPL-3.0-or-later

FROM golang:1.25@sha256:779b230b2508037a8095c9e2d223a6405f8426e12233b694dbae50197b9f6d04 AS builder
WORKDIR /app
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:24.04@sha256:d1e2e92c075e5ca139d51a140fff46f84315c0fdce203eab2807c7e495eff4f9
RUN <<EOF bash -xe
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
    jq mstflint
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
