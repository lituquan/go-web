FROM golang

MAINTAINER Razil "zc6496359"

WORKDIR $GOPATH/src/godocker

ADD hello.go $GOPATH/src/godocker

RUN go build hello.go

EXPOSE 8080

ENTRYPOINT ["./hello"]