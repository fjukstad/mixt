add <- function(a,b){
    return(a+b)
}


plt <- function() { 
    mat <- rnorm(1000)
    filename <- "images/plot.png"
    png (filename)
    hist(mat)
    dev.off()
    return (filename)
}  
