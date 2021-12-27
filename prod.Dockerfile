FROM golang:alpine as builder

# Set go compiler to use modules
ENV GO111MODULE=on

RUN apk --no-cache add gcc g++ make ca-certificates

RUN mkdir -p /go/src/github.com/nedson202/auth-manager
WORKDIR /go/src/github.com/nedson202/auth-manager

COPY . .

RUN apk update -qq && apk add git

RUN CGO_ENABLED=0 GOOS=linux cd cmd/server && go build -installsuffix cgo -o bin/auth-manager .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/nedson202/auth-manager/cmd/server/bin/auth-manager .

CMD ["./auth-manager"]
