var log = d3.scale.log()

var modules;
var ymodules = []

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



cellHeight = 15;
cellWidth = cellHeight;
margin = 120;


$(function() {
    svg = d3.select("#heatmap")
        .append("svg")

    legend = d3.select("#legend")
        .append("svg")

})

function loadEigengeneHeatmap(tissueA, tissueB) {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/eigengene"
    heatmap(url, tissueA, tissueB)
}

function loadGeneOverlapHeatmap(tissueA, tissueB) {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/overlap"
    heatmap(url, tissueA, tissueB)
}


function loadROIHeatmap(tissueA, tissueB) {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/roi"
    heatmap(url, tissueA, tissueB)
}


function loadPatientRankHeatmap(tissueA, tissueB) {
    var url = "/tissues/" + tissueA + "/" + tissueB + "/patientrank"
    heatmap(url, tissueA, tissueB)
}


function heatmap(url, tissueA, tissueB) {


    modules = []
    ymodules = []

    min = 100
    max = 0

    d3.csv(url, function(d) {

            ymodules.push(d.module)
            delete d.module

            modules = Object.keys(d)
            var data = Object.keys(d).map(function(key) {
                var num = -log(parseFloat(d[key]))
                if (num > 10) {
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
                modules: modules,
                data: data
            };
        },

        function(error, csvRows) {

            $("#heatmap svg").html("")
            $("#legend svg").html("")



            modules = strip(modules, "ME")
            ymodules = strip(ymodules, "ME")

            cellMargin = cellHeight + 2
            width = cellMargin * modules.length,
                height = cellMargin * ymodules.length;


            xScale = d3.scale.ordinal()
                .domain(modules)
                .rangePoints([0, width - cellWidth]);

            yScale = d3.scale.ordinal()
                .domain(ymodules)
                .rangePoints([0, height - cellHeight]);

            xAxis = d3.svg.axis().scale(xScale)
                .orient("bottom")
                .ticks(modules.length)
                .innerTickSize(2)
                .outerTickSize(0)


            yAxis = d3.svg.axis().scale(yScale)
                .orient("left")
                .ticks(ymodules.length);

            color = d3.scale.linear()
                .domain([min, max])
                .range(["#fde0dd", "#c51b8a"])

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

            row = rows.selectAll("svg")
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
                .attr("xlink:href", function(d,i){
                    
                    xname = modules[i];
                    yname = ymodules[d.index];
                    return "/compare/"+tissueA+"/"+tissueB + "/"+yname+"/"+xname
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

                    xname = modules[i];
                    yname = ymodules[d.index];
    
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
                .attr("transform", function(d) {
                    return "rotate(-65)"
                });

            yg = svg.append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + margin + "," + cellHeight / 2 + ")")
                .call(yAxis);

            var scale = d3.range(min, max, (max - min) / 10);

            legend.attr("width", function() {
                    var w = scale.length + 2
                    w = w * cellWidth
                    return w
                })
                .attr("height", cellHeight);

            legendg = legend.selectAll("g")
                .data(scale)
                .enter()
                .append("g")
                .attr("transform", function(d, i) {
                    return "translate(" + i * cellWidth + ",0)"
                })

            legendg.append("rect")
                .attr("width", cellWidth)
                .attr("height", cellHeight)
                .style("fill", function(d) {
                    return color(d)
                })

            legendg.append("text")
                .attr("transform", function(d, i) {
                    if (i == 0) {
                        return "translate(3," + 13 + ")"
                    }
                    if (i == scale.length - 1) {
                        var l = max.toFixed(1).toString().length
                        l = -3 * l
                        return "translate(" + l + "," + 13 + ")"
                    }
                })
                .style("fill", function(d, i) {
                    //if (i == scale.length - 1) {
                    //    return "white"
                    //}
                })
                .text(function(d, i) {
                    if (i == 0) {
                        return min.toFixed(0)
                    }
                    if (i == scale.length - 1) {
                        return max.toFixed(1)
                    }
                    return ""
                })
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
