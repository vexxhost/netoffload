# SPDX-FileCopyrightText: © 2025 VEXXHOST, Inc.
# SPDX-License-Identifier: GPL-3.0-or-later

FROM golang:1.27@sha256:c292963ad0e5f51e843c726edc204a4e4c26a859ba1e4d736d91028aa3633048 AS builder
WORKDIR /app
COPY . /app
RUN go build -v -o offloadctl cmd/offloadctl/main.go

FROM ubuntu:26.04@sha256:6a06111c3b0a8600feadb270aec771dc239bf407b25625f93f44ab00e11f4d3c
RUN <<EOF bash -xe
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
    jq mstflint
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY --from=builder /app/offloadctl /usr/local/bin/offloadctl
ENTRYPOINT ["/usr/local/bin/offloadctl"]
