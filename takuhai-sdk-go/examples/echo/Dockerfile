FROM golang:latest

WORKDIR /go/src/takuhai-sdk-go-tmp
ADD . /go/src/takuhai-sdk-go-tmp

ENV GOPATH /go
#RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep ensure


CMD ["go", "run", "main.go"]