FROM golang:1.20.5 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

RUN mkdir -p /gout_web

WORKDIR /gout_web
COPY . .
RUN go mod tidy

