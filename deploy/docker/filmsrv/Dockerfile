# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.15 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-w" -o .  ./cmd/filmsrv


FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /app/filmsrv/
COPY --chown=65534:65534 --from=builder /app/filmsrv .
COPY --chown=65534:65534 --from=builder /app/configs/filmsrv ./configs/filmsrv
USER 65534
#
EXPOSE 8081
ENTRYPOINT [ "./filmsrv", "server", "--port=8081" ]
