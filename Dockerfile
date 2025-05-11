# syntax=docker/dockerfile:1

# Build full image
FROM golang:1.24.3 AS builder
COPY . /app

WORKDIR /app
RUN go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -o dev \
    .

# Export downsized
FROM scratch as final
WORKDIR /root/
COPY --from=builder /app/dev .

ARG Token
ENTRYPOINT [ "./dev" ]
