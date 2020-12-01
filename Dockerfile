# syntax = docker/dockerfile:experimental

FROM golang:1.13

WORKDIR /app

COPY ./ ./

RUN make fmt

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

RUN chmod +x $GOPATH/bin/terraform-provider-prismacloud