# Use builds from launchpad
FROM ubuntu:latest

MAINTAINER Bjørn Fjukstad <bjorn@cs.uit.no>

# Install.
RUN \
  apt-get upgrade && \
  apt-get -y dist-upgrade && \
  apt-get update && \
  apt-get install -y  \ 
        r-base \ 
        gdebi-core\
        apache2 \
        apache2-utils \ 
        python3 \
        libapache2-mod-proxy-html \
        libapache2-modsecurity \ 
        libxml2-dev \
        iptables \ 
        git \
        wget 

WORKDIR /tmp
RUN wget https://download2.rstudio.org/rstudio-server-0.99.892-amd64.deb
RUN gdebi -n rstudio-server-0.99.892-amd64.deb

RUN a2enmod headers   

# Libraries etc that we depend on 
ADD deps.r /tmp/deps.r
WORKDIR /tmp/
RUN R --vanilla < deps.r

# Apache ports (without caching)
EXPOSE 80
EXPOSE 443 

# directories etc. 
RUN mkdir -p /home/rstudio
RUN mkdir -p /home/data
RUN mkdir -p /usr/local/lib/R/site-library 

# fix permissions to the shared r pkg location so that we can read/write
# to it from the docker host 
RUN chmod -R g+w /usr/local/lib/R/site-library
RUN chmod -R g+w /home/rstudio

# Set up passwords and that
RUN htpasswd -b -c /etc/http.passwd biopsy@mcgill van-mi-ka-al

# different config etc etc
ADD rstudio.conf /etc/apache2/sites-enabled/rstudio.conf 
ADD rstudio.conf /etc/apache2/sites-available/rstudio.conf 

# install path for r packages 
ADD rsession.conf /etc/rstudio/rsession.conf
ADD rserver.conf /etc/rstudio/rserver.conf

# add empty apache html page
ADD index.html /var/www/html/index.html

# fix group permissions so that the rstudio user is in the Hallett group and can
# read/write to the data dir
RUN groupadd -g 1001 Hallett
RUN useradd -g 1001 -u 1000 -p $(openssl passwd -1 rstudio) rstudio

# Git RAW
ADD https://raw.githubusercontent.com/fjukstad/git-raw/master/git-raw /bin/git-raw
RUN chmod +x /bin/git-raw

ADD git-raw-fix /bin/git-raw-fix
RUN chmod +x /bin/git-raw-fix
RUN chmod 755 /bin/git-raw 

# Define default command.
#CMD rstudio-server restart 
#CMD service apache2 restart && \
#    /usr/lib/rstudio-server/bin/rserver
RUN a2enmod headers \  
    	&& a2enmod proxy \
	&& a2enmod proxy_http 
