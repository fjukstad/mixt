$("#module-select-btn").click(function(){
    console.log("button clicked");
     var str = "";

     var tissue = $("ul#tissue-tab li.active").attr("tissue")
     modules = $("#"+tissue+"-module-select").val();
     str = modules.join([separator = '+']); 
    location.assign(location.href+"/"+tissue+"/"+str)
}); 

// sigma instance
var s, filter; 

$(function(){
 $('#tissue-tab a:first').tab('show')
    s = new sigma({
        container: 'graph-container',
        settings: {
            edgeColor: "source",
            defaultNodeColor: "#fab" ,
            labelThreshold: 12,
            maxNodeSize: 2, 
            batchEdgesDrawing: true, 
            defaultEdgeType: "curve", 
        }
    });

    var h = $("#module-selector").height();
    var w = $("#module-selector").width();
 
    $("#graph-container").height(h)
    $("#graph-container").width(w)

    filter = new sigma.plugins.filter(s);
    var tissue = $("ul#tissue-tab li.active").attr("tissue")
    TOMGraph(s, tissue); 


}); 


// module selected 
$("select").click(function(){
    filter.undo().apply();
        
     var tissue = $("ul#tissue-tab li.active").attr("tissue")
    modules = $("#"+tissue+"-module-select").val();
     if(modules != null){ 
        filter.nodesBy(function(n){
            if($.inArray(n.module, modules) != -1){
                return true
            } else {
                return false
            }
        }).apply();
     }

        
})



// select tissue 
$("ul#tissue-tab li").mouseup(function(){
    setTimeout(function(){
        var tissue = $("ul#tissue-tab li.active").attr("tissue")
        TOMGraph(s, tissue) 
    }, 10);
})



function TOMGraph(s, tissue){
    s.graph = s.graph.clear() 

   //s.bind("overNode", function(d){
   //    filter.neighborsOf(d.data.node.id).apply();
   //})

   //s.bind("outNode", function(n){
   //    filter.undo().apply();
   //})

    
    $.get('/tomgraph/'+tissue+'/nodes',
        function(data){
            nodes = JSON.parse(data);
            $.get('/tomgraph/'+tissue+'/edges',
                function(data){
                    edges = JSON.parse(data)

                    for(var i = 0; i < nodes.length; i++){
                        nodes[i].label = nodes[i].id + " ("+nodes[i].color+")"
                        nodes[i].module = nodes[i].color 
                        nodes[i].color = getHexColor(nodes[i].color)
                        nodes[i].size = 1; 
                    }

                    s.graph.read({
                        nodes: nodes,
                        edges: edges
                    })

                    s.refresh() 

                })
        })

}

function getHexColor(colorStr) {
    var a = document.createElement('div');
    a.style.color = colorStr;
    var colors = window.getComputedStyle( document.body.appendChild(a) ).color.match(/\d+/g).map(function(a){ return parseInt(a,10); });
    document.body.removeChild(a);
    return (colors.length >= 3) ? '#' + (((1 << 24) + (colors[0] << 16) + (colors[1] << 8) + colors[2]).toString(16).substr(1)) : false;
}
