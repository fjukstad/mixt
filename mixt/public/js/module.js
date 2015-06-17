$(function(){
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


        $("td.set").mouseover(function(d){
            var genesetname = d.target.innerText; 
            var url = "http://"+ location.hostname + ":" + location.port+"/geneset/abstract/"+genesetname
            var n = $(this)
            $.get(url, function(d){
                n.attr("title", d) 
                n.tooltip('show') 
                console.log(d)
            }); 
        })


 new Tablesort(document.getElementById('gene-table'));

 $('#settab a:first').tab('show')
});
