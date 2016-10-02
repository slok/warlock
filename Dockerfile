FROM golang:1.7-alpine


RUN apk --update add tar git bash wget && rm -rf /var/cache/apk/*

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid warlock
RUN adduser -D -u $uid -G warlock warlock

RUN mkdir -p /go/src/github.com/slok/warlock/
RUN chown -R warlock:warlock /go

WORKDIR /go/src/github.com/slok/warlock/

USER warlock

# Install dependency manager
RUN go get github.com/Masterminds/glide
