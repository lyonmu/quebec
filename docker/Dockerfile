FROM golang:alpine AS builder

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update && apk upgrade && apk add --no-cache build-base git

WORKDIR /build


FROM alpine:latest AS base

FROM base

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update && apk upgrade && apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone && update-ca-certificates \
    && echo TZ='Asia/Shanghai'; export TZ >> /etc/profile && source /etc/profile