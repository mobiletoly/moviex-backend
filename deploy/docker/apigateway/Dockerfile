# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.15 AS builder
#FROM golang:1.17-bullseye as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-w" -o .  ./cmd/apigateway


FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /app/apigateway/
COPY --chown=65534:65534 --from=builder /app/apigateway .
COPY --chown=65534:65534 --from=builder /app/configs/apigateway ./configs/apigateway
USER 65534
#
EXPOSE 8080
ENTRYPOINT [ "./apigateway", "server", "--deployment=k8s" ]
