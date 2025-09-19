# Copyright (c) 2025 tafaust
# SPDX-License-Identifier: MIT

# Build the provider
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the provider
RUN CGO_ENABLED=0 go build -o terraform-provider-peekaping

# Create final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the provider binary
COPY --from=builder /workspace/terraform-provider-peekaping .

# Set the entrypoint
ENTRYPOINT ["./terraform-provider-peekaping"]
