# SPDX-FileCopyrightText: © 2025 VEXXHOST, Inc.
# SPDX-License-Identifier: GPL-3.0-or-later

FROM golang:1.26@sha256:672c944cfd07599b5f1bb179869270cdc94c6254e41db4d431259693cb286b82 AS builder
WORKDIR /app
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:26.04@sha256:c098726dd9b9733f3e885f8de0ef375206a311b0a317eadc40d347b1115fb944
RUN <<EOF bash -xe
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
    jq mstflint
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
