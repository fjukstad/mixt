var log = d3.scaleLog()

var xnames;
var ynames = []

var max,
    min;

var cellHeight,
    cellWidth,
    height,
    width,
    xScale,
    yScale,
    xAxis,
    yAxis,
    color,
    svg,
    rows,
    row,
    xg,
    yg,
    legend;

var retry = true; 

cellHeight = 15;
cellWidth = cellHeight;
margin = 125;

function loadRanksumHeatmap(tissueA,tissueB,cohort){
    var url = "/tissues/" + tissueA + "/" + tissueB + "/ranksum/"+cohort
    heatmap(url, tissueA, tissueB) 
}

function loadEigengeneHeatmap(tissueA, tissueB, cohort="all") {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/eigengene/"+cohort
    heatmap(url, tissueA, tissueB)
}

function loadGeneOverlapHeatmap(tissueA, tissueB, cohort="all") {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/overlap/"+cohort
    heatmap(url, tissueA, tissueB)
}


function loadROIHeatmap(tissueA, tissueB, cohort="all") {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/roi/"+cohort
    heatmap(url, tissueA, tissueB)
}


function loadPatientRankHeatmap(tissueA, tissueB, cohort="all") {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/patientrank/"+cohort
    heatmap(url, tissueA, tissueB)
}

function loadClinicalEigengeneHeatmap(tissue, cohort="all") {
    var url = "/clinical-comparison/" + tissue + "/eigengene/"+cohort
    heatmap(url, tissue, "")
}

function loadClinicalROIHeatmap(tissue, cohort="all") {
    var url = "/clinical-comparison/" + tissue + "/roi/"+cohort
    heatmap(url, tissue, "")
}

function loadClinicalRanksumHeatmap(tissue, cohort="all") {
    var url = "/clinical-comparison/" + tissue + "/ranksum/"+cohort
    heatmap(url, tissue, "")
}

function heatmap(url, tissueA, tissueB) {
    loadStart() ;

    svg = d3.select("#heatmap-" + tissueA + " svg").attr("id", tissueA)
    //legend = d3.select("#legend-" + tissueA + " svg").attr("id", tissueA)
    
    xlabel = tissueB
    
    if(tissueB == ""){
        ylabel = "Clinical" 
    } else { 
        ylabel = tissueA
    }
    

    xnames = []
    ynames = []

    min = 100
    max = 0

    d3.csv(url, function(d) {
            if (tissueB == "") {
                ynames.push(d.Clinical)
                delete d.Clinical
            } else {
                ynames.push(d.module)
                delete d.module
            }

            xnames = Object.keys(d)
            var data = Object.keys(d).map(function(key) {
                var num = -log(parseFloat(d[key]))
                if (num > 10 || isNaN(num)) {
                    num = 10;
                }
                return num
            });

            localmax = d3.max(data)
            if (localmax > max) {
                max = localmax
            }


            localmin = d3.min(data)
            if (localmin < min) {
                min = localmin;
            }

            return {
                xnames: xnames,
                data: data
            };
        },

        function(error, csvRows) {
            if (csvRows.length < 1) {
                if(retry) { 
                    retry =false; 
                    heatmap(url, tissueA, tissueB)
                }
                
            }

            $("#heatmap-" + tissueA + " svg").html("")
            $("#legend-" + tissueA + " svg").html("")

            xnames = strip(xnames, "ME")
            ynames = strip(ynames, "ME")

            cellMargin = cellHeight + 2
            width = cellMargin * xnames.length,
                height = cellMargin * ynames.length;


            xScale = d3.scaleOrdinal()
                .domain(xnames)
                .range([0, width - cellWidth]);

            yScale = d3.scaleOrdinal()
                .domain(ynames)
                .range([0, height - cellHeight]);
            
            xAxis= d3.axisBottom(xScale)
                .ticks(xnames.length)
                .tickSizeInner(2)
                .tickSizeOuter(0)

            yAxis = d3.axisLeft(yScale)
                .ticks(ynames.length);


            color = d3.scaleThreshold()
                .domain([0,1.30103,2,3,4])
                .range(["#D4D4D4", "#D4D4D4", "#C7A3AA", "#B87A86", "#8E063B", "#8E063B"])

            svg.attr("width", width + margin * 2)
                .attr("height", height + margin)

            rows = svg.selectAll("g")
                .data(csvRows)
                .enter()
                .append("g")
                .attr("transform", function(d, i) {
                    return "translate(" + margin + "," + i * cellMargin + ")"
                })
                .attr("class", "column")
                .attr("id", function(d,i){
                    return ynames[i]
                });

            row = rows.selectAll("svg#" + tissueA)
                .data(function(d, i) {
                    res = []
                    for (j in d.data) {
                        a = {
                            data: d.data[j],
                            index: i
                        }
                        res.push(a)
                    }
                    return res;
                })
                .enter()
                .append("g")
                .append("a")
                .attr("xlink:href", function(d, i) {
                    xname = xnames[i];
                    yname = ynames[d.index];
                    if (tissueB == "") {
                        return "/modules/" + tissueA + "/" + xname + "/cohort/all"
                    } else {
                        cohort= $("#cohort-select").val();
                        return "/compare/" + tissueA + "/" + tissueB + "/" + yname + "/" + xname + "/cohort/"+cohort
                    }
                })
                .append("rect")
                .attr("transform", function(d, i) {
                    return "translate(" + i * cellMargin + ",0)"
                })
                .attr("height", cellHeight)
                .attr("width", cellWidth)
                .style("fill", function(d) {
                    var val = d.data
                    if (isNaN(val)) {
                        return "#d3d3d3"
                    }
                    return color(val)
                })
                .on("mouseover", function(d, i) {
                    d3.select(this)
                        .style("stroke", "black")
                        .style("stroke-width", 1)

                    xname = xnames[i];
                    yname = ynames[d.index];


                    xg.selectAll(".tick")
                        .each(function(d, i) {
                            if (d == xname) {
                                d3.select(this)
                                    .selectAll('text')
                                    .style("font-weight", "bold")
                            }
                        })
                    yg.selectAll(".tick")
                        .each(function(d, i) {
                            if (d == yname) {
                                d3.select(this)
                                    .selectAll('text')
                                    .style("font-weight", "bold")
                            }
                        })


                })
                .on("mouseout", function(d, i) {
                    d3.select(this).style("stroke-width", 0)

                    xg.selectAll(".tick")
                        .each(function(d, i) {
                            if (d == xname) {
                                d3.select(this)
                                    .selectAll('text')
                                    .style("font-weight", "")
                            }
                        })

                    yg.selectAll(".tick")
                        .each(function(d, i) {
                            if (d == yname) {
                                d3.select(this)
                                    .selectAll('text')
                                    .style("font-weight", "")
                            }
                        })


                })

            .append("svg:title")
                .text(function(d, i) {
                    var val = d.data
                    if (isNaN(val)) {
                        val = "NA"
                    }

                    return val;
                })

            xg = svg.append("g")
                .attr("class", "x axis")
                .attr("transform", "translate(" + (margin + cellWidth / 2) + "," + height + ")")
                .call(xAxis)


            xg.selectAll("text")
                .style("text-anchor", "end")
                .attr("dx", "-.8em")
                .attr("dy", ".15em")
                .attr("transform", function(d,i) {
                    return "rotate(-65)"
                });

            yg = svg.append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + (margin-1) + "," + cellHeight / 2 + ")")
                .call(yAxis);

            var scale = d3.range(min, max, (max - min) / 10);

                ytrans = (height + margin) / 2.5
                svg.append("g")
                    .attr("transform","translate(15,"+ytrans+") rotate(-90)")
                    .append("text")
                    .attr("class", "label-heatmap")
                    .text(ylabel)
                
                ytrans = height + margin - 10; 
                xtrans = (width + margin*1.25) / 2 

                svg.append("g")
                    .attr("transform","translate("+xtrans+","+ytrans+")")
                    .append("text")
                    .attr("class", "label-heatmap")
                    .text(xlabel)


            d3.selectAll("#"+tissueA+" g.y.axis g")  
                .attr("transform", function(d, i){
                    console.log(d,i) 
                    return "translate(0,"+(i*cellMargin)+")"
                }) 
                .on("mouseover", function(){

                    d3.select(this).selectAll("text").attr("font-weight", "bold")
                        
                    t = getTranslation(d3.select(this).attr("transform")) 

                    x = t[0] + margin,
                    y = t[1] - 1;

                    x2 = width + margin - 1
                    y2 = y 

                    svg.append("line")
                        .attr("id", "row-select") 
                        .attr("x1", x)
                        .attr("y1", y)
                        .attr("x2", x2)
                        .attr("y2", y)
                        .attr("stroke-width", 1) 
                        .attr("stroke", "black") 

                    svg.append("line")
                        .attr("id", "row-select") 
                        .attr("x1", x)
                        .attr("y1", y + cellHeight - 1)
                        .attr("x2", x2)
                        .attr("y2", y + cellHeight - 1 )
                        .attr("stroke-width", 1) 
                        .attr("stroke", "black") 
                })
                .on("mouseout", function(){
                    d3.selectAll("line#row-select").remove() 
                    d3.select(this).selectAll("text").attr("font-weight", "")
                })
                
            d3.selectAll("#"+tissueA+" g.x.axis g")
                .attr("transform", function(d, i){
                    return "translate("+(i*cellMargin)+",0)"
                }) 
                .on("mouseover", function(){
                    d3.select(this).selectAll("text").attr("font-weight", "bold")
                    t = getTranslation(d3.select(this).attr("transform")) 
                    x = t[0] + margin - 1,
                    y = t[1];

                    y2 = height - 1
                    x2 = x

                    svg.append("line")
                        .attr("id", "column-select") 
                        .attr("x1", x)
                        .attr("y1", y)
                        .attr("x2", x2)
                        .attr("y2", y2)
                        .attr("stroke-width", 1) 
                        .attr("stroke", "black") 

                    x = x + cellWidth
                    x2 = x

                    svg.append("line")
                        .attr("id", "column-select") 
                        .attr("x1", x)
                        .attr("y1", y)
                        .attr("x2", x2)
                        .attr("y2", y2)
                        .attr("stroke-width", 1) 
                        .attr("stroke", "black") 
                })
                .on("mouseout", function(){
                    d3.selectAll("line#column-select").remove() 
                    d3.select(this).selectAll("text").attr("font-weight", "")
                }) 

            plotScale(tissueA)

            loadStop();
        })


}





function strip(str, substring) {
    for (var i = 0; i < str.length; i++) {
        str[i] = str[i].replace("ME", "")
    }
    return str
}

function replicate(what, howmany) {
    return Array(howmany + 1).join(1).split('').map(function() {
        return what;
    })
}

function toFloats(strings) {
    floats = []
    for (var i = 0; i < strings.length; i++) {
        floats.push(parseFloat(strings[i]))
    }
    return floats
}

function plotScale(tissue){
    console.log("color", color) 
    console.log(max)

    steps = [0,1.30103,2,3,4]
    
    var svg = d3.select("svg#legend-"+tissue);

    svg.attr("height", cellHeight);
    svg.attr("width", steps.length*cellMargin);

    var square = svg.selectAll("rect").data(steps)

    square.enter().append("rect")
        .attr("width", cellHeight)
        .attr("height", cellWidth)
        .attr("x", function(d,i){
            return i*cellMargin ;
        })
        .attr("y",0)
        .attr("fill", function(d){
            return color(d);
        })
        .append("title").text(function(d){
            return d;
        }) 
    
    d3.selectAll("p#legendmin").text("0")
    d3.selectAll("p#legendmax").text("4")
}

function getTranslation(transform) {
  // Create a dummy g for calculation purposes only. This will never
  // be appended to the DOM and will be discarded once this function 
  // returns.
  var g = document.createElementNS("http://www.w3.org/2000/svg", "g");
  
  // Set the transform attribute to the provided string value.
  g.setAttributeNS(null, "transform", transform);
  
  // consolidate the SVGTransformList containing all transformations
  // to a single SVGTransform of type SVG_TRANSFORM_MATRIX and get
  // its SVGMatrix. 
  var matrix = g.transform.baseVal.consolidate().matrix;
  
  // As per definition values e and f are the ones for the translation.
  return [matrix.e, matrix.f];
}

