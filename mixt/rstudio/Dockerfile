FROM rocker/rstudio

RUN apt-get update && \
	apt-get uptrade && \
	apt-get install -y \
		build-essential


ADD deps.r /tmp/deps.r
WORKDIR /tmp/
RUN R --vanilla < deps.r

