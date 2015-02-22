$("#tissue-select" ).change(function() {
    console.log( "winds of shit");
    var tissue = $("#tissue-select" ).val()
    var str = location.href.split("/");
    str[str.length-1] = tissue;
    var url = str.join([separator = '/']); 
    location.assign(url)
});
