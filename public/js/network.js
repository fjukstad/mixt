

function TOMGraph(s, tissue){
    //s.graph = s.graph.clear() 

    var h = $("#module-selector").height();
    var w = $("#module-selector").width();
 
    $("#graph-container-"+tissue).height(300)
    $("#graph-container-"+tissue).width(500)
   //s.bind("overNode", function(d){
   //    filter.neighborsOf(d.data.node.id).apply();
   //})

   //s.bind("outNode", function(n){
   //    filter.undo().apply();
   //})

    
    $.get('/tomgraph/'+tissue+'/nodes',
        function(data){
            var nodes = JSON.parse(data);
            $.get('/tomgraph/'+tissue+'/edges',
                function(data){
                    var edges = JSON.parse(data)

                    console.log("edges:", edges) 
                    console.log("nodes", nodes) 

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

                    s.startForceAtlas2();
                    setTimeout(function(){
                        s.killForceAtlas2();
                    }, 1000);
                        

                    console.log("refreshing") 

                    s.refresh(); 
                    console.log("all refreshed") 

                    // nasty window resize hack
                    window.dispatchEvent(new Event('resize'))

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
