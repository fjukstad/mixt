library(Hmisc) # for 'capitalize()'

# Where data and scripts are stored 
datadir <- "/go/src/bitbucket.org/vdumeaux/mixt/data"
scriptdir <- "/go/src/bitbucket.org/vdumeaux/mixt/src"

# Get helper scripts 
source(paste0(scriptdir, "/bresat.R"), chdir=TRUE)
source(paste0(scriptdir, "/heatmap.R"), chdir=TRUE)
source(paste0(scriptdir, "/mixt-utils.r"), chdir=TRUE)
source(paste0(scriptdir, "/modules.R"), chdir=TRUE)

# Get datafiles 
rawModulesFilename <- paste0(datadir, "/cc.blood-biopsy-Modules.RData")
exprsFilename <-  paste0(datadir, "/CC-Biopsy-Expressions.RData")
modulesFilename <-paste0(datadir, "/modules-complete-pepi.RData")

modules <- loadModulesAndROI(rawModulesFilename,exprsFilename,modulesFilename)

### Set Kvik option so that the output is readable in Kvik 
options(width=10000) 

### Where to store images
imgpath <- "images"
dir.create(imgpath,showWarnings = FALSE)
### Directory to store tables (output as csv files)
tablePath <- "tables"
dir.create(tablePath,showWarnings = FALSE)

### Generate heatmap plot for the given tissue and module. If the heatmap
### already exists, it finds the appropriate png file where it is supposed
### to store a new one, it returns this file. This heat map function
### generates both a png and a pdf. All plots are stored in the path
### given by 'imgpath'
heatmap <- function(tissue,module,imgpath="images") { 
  pngFilename <- paste(imgpath, "/heatmap-",tissue,"-",module,".png",sep="")
  if(file.exists(pngFilename)){
    return (pngFilename)
  } else {
    png(pngFilename)
    plot.new()
    create.modules.heatmap(modules[[tissue]]$bresat[[module]],modules[[tissue]]$clinical,
                           title=capitalize(paste(tissue, module)))
    dev.off()
    
    pdfFilename <- paste(imgpath, "/heatmap-",tissue,"-",module,".pdf",sep="")
    pdf(pdfFilename)
    plot.new()
    create.modules.heatmap(modules[[tissue]]$bresat[[module]],modules[[tissue]]$clinical,
                           title=capitalize(paste(tissue, module)))
    dev.off()
    
    return (pngFilename)
  } 
}

### Returns a list of modules found for the given tissue
getModules <- function(tissue) {
    return (names(modules[[tissue]]$modules))
}

### Returns the location of a csv files containing a list of all genes found in 
### all modules across all tissues. 
getAllGenes <- function(tablePath="tables"){
  filename = paste(tablePath,"/genes.csv",sep="")
  if(!file.exists(filename)){
    getAllGenesAndModules()
  }
  
  genesAndModules = read.csv(filename)
  g = genesAndModules$gene
  genes = matrix(g)
  colnames(genes) = c("gene")
  geneFilename = paste(tablePath,"/all-genes.csv",sep="")
  write.table(genes, geneFilename, sep=",",row.names=FALSE) 
  return(geneFilename)
}

### Get all modules a specific gene is found in. 
getAllModules <- function(gene) {
  filename = paste(tablePath,"/genes.csv",sep="")
  genesAndModules = read.csv(filename)
  id = match(gene,genesAndModules$gene)
  g = genesAndModules[id,]
  d = c(lapply(g,as.character))
  return(c(d$blood, d$biopsy))
}

### Retrieves all genes and the modules they participate in. 
### Writes it to a file so that we can read it later. 
getAllGenesAndModules <- function() {
  filename = paste(tablePath,"/genes.csv",sep="")
  if(file.exists(filename)){
    return (filename)
  } 
  res <- NULL
  tissues <- c("blood", "biopsy")
  for (tissue in tissues){
    for(module in names(modules[[tissue]]$modules)) {
      if(module == "grey"){
        next 
      }
      gs <- modules[[tissue]]$bresat[[module]]$gene.order
      for(gene in gs){
        if(length(res[[gene]])==0) {
          res[[gene]] = list()
          res[[gene]][["blood"]] = NA
          res[[gene]][["biopsy"]] = NA
        }
        res[[gene]][[tissue]] = c(module)
      }
    }
  }
  genes = matrix(unlist(res), nrow=length(names(res)))
  genes = cbind(names(res), genes)
  colnames(genes) <-  c("gene",tissues)
  write.table(genes, filename, sep=",",row.names=FALSE) 
  return (filename)
}

### Get available tissues
getTissues <- function() {
    return (names(modules))
}

### Get a list of genes for a specific module and tissue. Results are 
### written to a csv file and its location is returned. 
getGeneList <- function(tissue,module){
  filename <- paste(tablePath,"/genelist-",tissue,"-",module,".csv",sep="")
  if(file.exists(filename)){
    return (filename)
  } 
  
  genes <- modules[[tissue]]$bresat[[module]]$gene.order
  up.dn <- modules[[tissue]]$bresat[[module]]$up.dn
  res <- matrix(c(genes,up.dn), nrow=length(genes))
  colnames(res) <- c("Gene", "up.dn")
  write.table(res, filename, sep=",",row.names=FALSE) 
    
  return(filename) 
} 


