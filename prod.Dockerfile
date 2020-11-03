FROM golang:alpine AS build-env

# Set go compiler to use modules
ENV GO111MODULE=on

RUN apk update -qq && apk add git

WORKDIR /go/src/github.com/nedson202/auth-manager

RUN go get github.com/githubnemo/CompileDaemon

COPY . .

RUN go build -o /bin/auth-manager

ENTRYPOINT /bin/auth-manager
