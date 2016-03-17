packages <- c("ggplot2","Hmisc, Rcpp, roxygen2") 
install.packages(packages, repos='http://cran.us.r-project.org')

source("http://bioconductor.org/biocLite.R")
pkgs <- c("Biobase" ,"DBI" ,"RSQLite" ,"AnnotationDbi" ,"GO.db" ,"RColorBrewer" ,"latticeExtra" ,"colorspace" ,"munsell", "plyr" ,"scales" ,"reshape" ,"preprocessCore" ,"WGCNA" ,"illuminaHumanv3.db" ,"illuminaHumanv4.db" ,"animation" ,"limma")
biocLite(pkgs, ask=FALSE) 
