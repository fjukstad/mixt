add <- function(a,b){
    return(a+b)
}


plt <- function() { 
    mat <- rnorm(10)
    filename <- "images/plot.png"
    png (filename)
    hist(mat)
    dev.off()
    return (filename)
}  

getModules <- function(tissue) {
    return (c("red","green","blue","purple"))
}

getGenes <- function() { 
    return (c("BRCA1", "BRCA2", "ESR1"))
}

getTissues <- function() {
    return (c("Blood", "Biopsy"))
}
