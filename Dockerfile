# SPDX-FileCopyrightText: © 2025 VEXXHOST, Inc.
# SPDX-License-Identifier: GPL-3.0-or-later

FROM golang:1.27@sha256:2920bb124f1a2e12a1d39f0cab756ae0068441ee1da7655dec7f5e374944095b AS builder
WORKDIR /app
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:26.04@sha256:0e73b71479e7bc5f006bca872e25c3cd680dc052a0116f90d63ae363537db0dd
RUN <<EOF bash -xe
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
    jq mstflint
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
