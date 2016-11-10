FROM golang

RUN go get github.com/fjukstad/mixt
EXPOSE 80
WORKDIR $GOPATH/src/github.com/fjukstad/mixt
ENTRYPOINT mixt -port=:80
