#setwd("/Users/bjorn/pepi/guest.bci/bjorn/mixt/experiments/exp_mixt")

library(Hmisc) # for 'capitalize()'

### Get functions etc. 
source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/src/bresat.R", chdir=TRUE)
source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/src/heatmap.R", chdir=TRUE)
source("/Users/bjorn/pepi/guest.bci/bjorn/mixt/experiments/exp_mixt/utils.r")
# library(parallel)

modulesFilename <- "/Users/bjorn/pepi/guest.bci/bjorn/mixt/data/modules-complete.Rdata"
if(file.exists(modulesFilename)){ 
  load(modulesFilename)
} else {
  ### Load datasets 
  load("/Users/bjorn/pepi/guest.bci/bjorn/mixt/data/cc.blood-biopsy-Modules.RData") # modules
  load("/Users/bjorn/pepi/guest.bci/bjorn/mixt/data/CC-Biopsy-Expressions.RData")   # gene expression and others 
  
  
  names(cc.biopsy)<-c("blood", "biopsy")
  names(cc.biopsy.modules)<-c("blood", "biopsy")
  modules <- load.modules(cc.biopsy, cc.biopsy.modules)
  
  
  # Ranksum
  modules$blood$bresat <- lapply(modules$blood$modules[-1], function(mod) {
    sig.ranksum(modules$blood$exprs, ns=mod, full.return=TRUE)
  })
  modules$biopsy$bresat <- lapply(modules$biopsy$modules[-1], function(mod) {
    sig.ranksum(modules$biopsy$exprs, ns=mod, full.return=TRUE)
  })
  
  
  ### roi function
  roi<-NULL
  
  for (tissue in c("blood", "biopsy"))
  {
    module.names <- names(modules[[tissue]]$bresat)
    roi[[tissue]]<- mclapply(module.names, function(module) {
      random.ranks(modules[[tissue]]$bresat[[module]])
    })
    names(roi[[tissue]])<-module.names
  }  
  
  for (tissue in c("blood", "biopsy"))
  {
    module.names <- names(modules[[tissue]]$bresat)
    for (module in module.names){
      modules[[tissue]]$bresat[[module]]$roi<-roi[[tissue]][[module]]
    }
  }
  
  ### define roi categories
  roi.cat<-NULL
  for (tissue in c("blood", "biopsy"))
  {
    module.names <- names(modules[[tissue]]$bresat)
    roi.cat[[tissue]]<- mclapply(module.names, function(module) {
      define.roi.regions(modules[[tissue]]$bresat[[module]],modules[[tissue]]$bresat[[module]]$roi)
    })
    names(roi.cat[[tissue]])<-module.names
  }  
  for (tissue in c("blood", "biopsy"))
  {
    module.names <- names(modules[[tissue]]$bresat)
    for (module in module.names){
      modules[[tissue]]$bresat[[module]]$roi.cat<-roi.cat[[tissue]][[module]]
    }
  }  
  save(modules,file=modulesFilename)
}


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
  if(file.exists(modulesFilename)){
    return (filename)
  } else {
    png(filename)
    plot.new()
    create.modules.heatmap(modules[[tissue]]$bresat[[module]],modules[[tissue]]$clinical,
                           title=capitalize(paste(tissue, module)))
    dev.off()
    return (filename)
  } 
}

getModules <- function(tissue) {
    return (names(modules[[tissue]]$modules))
}

getGenes <- function(tissue, module) {  
    return (c("BRCA1", "BRCA2", "ESR1"))
}

getTissues <- function() {
    return (names(modules))
}

getGeneList <- function(tissue,module){
    genes <- modules[[tissue]]$bresat[[module]]$gene.order
    up.dn <- modules[[tissue]]$bresat[[module]]$up.dn
    res <- matrix(c(genes,up.dn), nrow=length(genes))
    colnames(res) <- c("Gene", "up.dn")
    path <- "tables"
    dir.create(path)
    filename <- paste(path,"/genelist-",tissue,"-",module,".csv",sep="")
    write.table(res, filename, sep=",",row.names=FALSE) 
    
    return(filename) 
} 


