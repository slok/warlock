FROM slok/warlock_base:latest
USER root
RUN apk add --no-cache g++

USER warlock
RUN go get github.com/golang/mock/mockgen
