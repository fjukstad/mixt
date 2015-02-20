$(function() {
    var cache = {};
    $( "input#search" ).autocomplete({
      minLength: 2,
      source: function( request, response ) {
        var term = request["term"];
        console.log(request) 
        if ( term in cache ) {
          response( cache[ term ] );
          return;
        }
        $.getJSON("/search/"+request.term, function( data, status, xhr ) {
            console.log(data)
            var genes  = data.Terms; 
            for(var i = 0; i < genes.length; i++){
                gene = genes[i]
                cache[gene] = gene;
            }
            response(genes);
            return
            }
        );
      }
    });

    $('input#search').bind("enterKey", function(e){
        searchterm = $('input#search').val()
        entry = cache[searchterm]
        if(typeof entry === 'undefined'){
            swal({
                title: "Could not find what you we're searching for, sorry!",
                text: "You searched for '"+searchterm+"'",
                type: "warning"
            }) 
            return
        } else { 
            //window.location = window.location.origin+"/pathway/"+id
        }
    });

    $('input#search').keyup(function(e){
        if(e.keyCode == 13) {
            $(this).trigger("enterKey");
        }
    }) 
});
