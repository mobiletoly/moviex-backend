# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.15 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-w" -o .  ./cmd/usersrv


FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /app/usersrv/
COPY --chown=65534:65534 --from=builder /app/usersrv .
COPY --chown=65534:65534 --from=builder /app/configs/usersrv ./configs/usersrv
USER 65534
#
EXPOSE 8082
ENTRYPOINT [ "./usersrv", "server", "--port=8082" ]
