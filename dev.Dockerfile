FROM golang:alpine

# Set go compiler to use modules
ENV GO111MODULE=on

RUN apk update -qq && apk add git

WORKDIR /go/src/github.com/nedson202/auth-manager

RUN go get github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT cd cmd/server && CompileDaemon -build="go build -o bin/auth-manager-rpc ." -command="./bin/auth-manager-rpc" -color -graceful-kill
