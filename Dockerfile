# SPDX-FileCopyrightText: © 2025 VEXXHOST, Inc.
# SPDX-License-Identifier: GPL-3.0-or-later

FROM golang:1.27@sha256:1cfcdb11f37fce9429f617100f39e0251748bbaba454bd431275155701765058 AS builder
WORKDIR /app
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:26.04@sha256:da6fc2be547864451aa253836dd926da33623312df4a9a243e35dc877c378a78
RUN <<EOF bash -xe
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
    jq mstflint
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
