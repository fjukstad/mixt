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

# zeromq
RUN wget http://download.zeromq.org/zeromq-4.0.4.tar.gz && \
    tar -xvf zeromq-4.0.4.tar.gz && \
    cd zeromq-* && \
    ./configure && \
    make install && \
    ldconfig && \
    cd .. && \
    rm -rf zeromq*

# zmq bindings for python 
RUN pip install pyzmq

# Add Kvik and MIxT repos 
ADD kvik github.com/fjukstad/kvik
ADD mixt bitbucket.org/vdumeaux/mixt

