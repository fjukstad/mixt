$(function(){

    var options = {
      valueNames: [ 'name', 'correlation','k','kin', 'updown' ]
    };

    var geneList = new List('genes', options);
    
    //geneList.sort('name', { order: "desc" });
    //
      $('[data-toggle="tooltip"]').tooltip({'placement': 'bottom'})

        $("td.name").mouseover(function(d){
            var name = d.target.innerText; 
            var url = "http://"+ location.hostname + ":" + location.port+"/gene/summary/"+name
            var n = $(this)
            $.get(url, function(d){
                n.attr("title", d) 
                n.tooltip('show') 
                console.log(d)
            }); 
        }); 
});
