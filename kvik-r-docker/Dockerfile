FROM fjukstad/kvik-r

RUN apt-get install -y default-jre \
    default-jdk \
    libmariadb-client-lgpl-dev
    

# Install R package dependencies 
RUN R -e 'install.packages(c("ggplot2", "Hmisc", "Rcpp", "roxygen2", "jsonlite", "igraph", "dplyr", "parallel", "colorspace", "ic10", "igraph", "network","GGally","sna"), repos="http://cran.rstudio.com")'
RUN R -e 'source("https://bioconductor.org/biocLite.R"); biocLite(c("genefu", "WGCNA", "impute", "preprocessCore", "GO.db"), dependencies=TRUE)'


COPY mixt_1.0_R_x86_64-pc-linux-gnu.tar.gz /mixt_1.0_R_x86_64-pc-linux-gnu.tar.gz 
RUN R CMD INSTALL /mixt_1.0_R_x86_64-pc-linux-gnu.tar.gz 
RUN rm /mixt_1.0_R_x86_64-pc-linux-gnu.tar.gz

WORKDIR /
COPY BCIutils_0.99.0.tar.gz BCIutils_0.99.0.tar.gz
RUN R CMD INSTALL BCIutils_0.99.0.tar.gz
RUN rm BCIutils_0.99.0.tar.gz

COPY huc huc
RUN R CMD INSTALL huc
RUN rm -rf huc

WORKDIR /go/src/github.com/fjukstad/kvik/r/examples
