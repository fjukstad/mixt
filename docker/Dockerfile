FROM golang

# update get packages
RUN apt-get update && apt-get install -y \
    build-essential \
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

RUN go get github.com/fjukstad/mixt-blood-tumor 
EXPOSE 80
WORKDIR /
ENTRYPOINT mixt-blood-tumor -port=:80
