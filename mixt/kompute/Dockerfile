FROM rocker/rstudio

MAINTAINER Bjørn Fjukstad <bjorn@cs.uit.no>

# From the rstudio docker file at 
# https://registry.hub.docker.com/u/opencpu/rstudio/dockerfile/ 
RUN \
  useradd -ms /bin/bash builder && \
  apt-get update && \
  apt-get -y dist-upgrade && \
  apt-get install -y wget make devscripts apache2-prefork-dev apache2-mpm-prefork libapreq2-dev r-base r-base-dev libapparmor-dev libcurl4-openssl-dev libprotobuf-dev libprotoc-dev xvfb xauth xfonts-base curl libssl-dev pkg-config libxml2-dev

USER builder

RUN \
  cd ~ && \
  wget https://github.com/jeffreyhorner/rapache/archive/v1.2.6.tar.gz && \
  tar xzf v1.2.6.tar.gz && \
  cd rapache-1.2.6 && \
  dpkg-buildpackage

RUN \
  cd ~ && \
  wget https://github.com/jeroenooms/opencpu-server/archive/v1.4.6.tar.gz && \
  tar xzf v1.4.6.tar.gz && \
  cd opencpu-server-1.4.6 && \
  dpkg-buildpackage

USER root

RUN \
  cd /home/builder && \
  dpkg -i libapache2-mod-r-base_*.deb && \
  dpkg -i opencpu-lib_*.deb && \
  dpkg -i opencpu-server_*.deb

RUN \
  apt-get purge -y devscripts apache2-prefork-dev libapreq2-dev \ 
  libapparmor-dev libprotobuf-dev libprotoc-dev xvfb xauth pkg-config && \
  apt-get autoremove --purge -y && \
  deluser builder --remove-home

# Libraries etc that we depend on 
ADD deps.r /tmp/deps.r
WORKDIR /tmp/
RUN R --vanilla < deps.r

# Apache ports (without caching)
EXPOSE 80
EXPOSE 443
EXPOSE 8004

# Add r server config file to only accept connections on localhost
# ADD rserver.conf /etc/rstudio/rserver.conf 

VOLUME /home/rstudio                    # home directory of rstudio user
VOLUME /home/data                       # mcgill data dir 
VOLUME /usr/local/lib/R/site-library    # r pkg install location

# Set up passwords and that
RUN htpasswd -b -c /etc/http.passwd biopsy@mcgill van-mi-ka-al

# different config etc etc
ADD rstudio.conf /etc/apache2/sites-enabled/rstudio.conf 
ADD rstudio.conf /etc/apache2/sites-available/rstudio.conf 
ADD opencpu.conf /etc/apache2/sites-enabled/opencpu.conf 
ADD opencpu.conf /etc/apache2/sites-available/opencpu.conf 

# install path for r packages 
ADD rsession.conf /etc/rstudio/rsession.conf

# add empty apache html page
ADD index.html /var/www/html/index.html

# fix permissions to the shared r pkg location so that we can read/write
# to it from the docker host 
RUN chmod -R g+w /usr/local/lib/R/site-library

# fix group permissions so that the rstudio user is in the Hallett group and can
# read/write to the data dir
RUN groupadd -g 1001 Hallett
RUN usermod -g 1001 rstudio

# Create users that can use rstudio etc
# RUN useradd -s /bin/bash -d /home/bjorn -u 501 -g 20 bjorn
# RUN echo 'bjorn:test' | chpasswd
# RUN addgroup bjorn staff
# RUN mkdir /home/bjorn && chown bjorn:staff /home/bjorn

# Define default command.
CMD service opencpu restart && rstudio-server restart && \
    tail -F /var/log/opencpu/apache_access.log
