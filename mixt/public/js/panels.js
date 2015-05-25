$("#tissue-select" ).change(function() {
    console.log( "winds of shit");
    var tissue = $("#tissue-select" ).val()

    $.cookie("tissue", tissue, { domain: '', path: '' });



    var str = location.href.split("/");
    str[str.length-1] = tissue;
    var url = str.join([separator = '/']); 
    location.assign(url)
});
