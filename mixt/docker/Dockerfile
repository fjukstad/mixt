FROM golang

# update get packages
RUN apt-get update && apt-get install -y \
    build-essential \
    python \
    python-dev \
    python-pip \
    wget \
    git \
    pkg-config \
    tar \
    r-base \
    libtool \ 
    autoconf \
    automake \
    sshfs \ 
    nfs-common \
    vim

RUN go get github.com/PuerkitoBio/goquery && \
    go get github.com/fjukstad/gocache && \
    go get github.com/gorilla/mux && \
    go get github.com/gorilla/securecookie && \
    go get code.google.com/p/gcfg

RUN mkdir -p /go/src/github.com/fjukstad
RUN mkdir -p /go/src/bitbucket.org/vdumeaux/mixt

VOLUME /go/src/github.com/fjukstad/kvik 
VOLUME /go/src/bitbucket.org/vdumeaux/mixt

# Add Kvik and MIxT repos 
#ADD /home/core/kvik/ /go/src/github.com/fjukstad/kvik
#ADD /home/core/mixt/ /go/src/bitbucket.org/vdumeaux/mixt

EXPOSE 8080

WORKDIR /go/src/bitbucket.org/vdumeaux/mixt/mixt/

ENTRYPOINT go run main.go
