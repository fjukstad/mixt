#setwd("/Users/bjorn/pepi/guest.bci/bjorn/mixt/experiments/exp_mixt")

library(Hmisc) # for 'capitalize()'

### Get functions etc. 
source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/src/bresat.R", chdir=TRUE)
source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/src/heatmap.R", chdir=TRUE)
# library(parallel)


### Load datasets 
load("/Users/bjorn/pepi/guest.bci/bjorn/mixt/data/cc.blood-biopsy-Modules.RData") # modules
load("/Users/bjorn/pepi/guest.bci/bjorn/mixt/data/CC-Biopsy-Expressions.RData")   # gene expression and others 

source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/experiments/exp_mixt/utils.r")
names(cc.biopsy)<-c("blood", "biopsy")
names(cc.biopsy.modules)<-c("blood", "biopsy")
modules <- load.modules(cc.biopsy, cc.biopsy.modules)

modules$blood$bresat <- lapply(modules$blood$modules[-1], function(mod) {
  sig.ranksum(modules$blood$exprs, ns=mod, full.return=TRUE)
})
modules$biopsy$bresat <- lapply(modules$biopsy$modules[-1], function(mod) {
  sig.ranksum(modules$biopsy$exprs, ns=mod, full.return=TRUE)
})


### Set Kvik option so that the output is readable in Kvik 
options(width=10000) 

### Where to store images
imgpath <- "images"
dir.create(imgpath,showWarnings = FALSE)


plt <- function() { 
    mat <- rnorm(10)
    filename <- paste(imgpath,"/plot.png",sep="")
    
    png (filename)
    hist(mat)
    dev.off()
    return (filename)
}  

heatmap <- function(tissue,module) { 
  filename <- paste(imgpath, "/heatmap-",tissue,"-",module,".png",sep="")
  png(filename)
  plot.new()
  create.modules.heatmap(modules[[tissue]]$bresat[[module]],modules[[tissue]]$clinical,
                         title=capitalize(paste(tissue, module)))
  dev.off()
  return (filename)
}

getModules <- function(tissue) {
    return (names(modules[[tissue]]$modules))
}

getGenes <- function() { 
    return (c("BRCA1", "BRCA2", "ESR1"))
}

getTissues <- function() {
    return (names(modules))
}

getGeneList <- function(module,tissue){ 
    genes <- matrix(c("BRCA1", "BRCA2", "ESR1",
                      0.8, 0.89, 0.7,
                      0.7, 0.81, 0.61,
                      1, 0.99, 0.97,
                      "up", "up", "up"),
                    nrow=3, byrow=FALSE)
    colnames(genes) <- c("Gene","Correlation", "k", "kin", "up/down")
    
    path <- "tables"
    dir.create(path)
    filename <- paste(path,"/genelist-",tissue,"-",module,".csv",sep="")
    write.table(genes, filename, sep=",",row.names=FALSE) 
    
    return(filename) 
} 


