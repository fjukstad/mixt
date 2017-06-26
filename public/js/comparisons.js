$(function(){
    $("td#comp").click(function(){
        var comp = $(this).attr("comp");
        loadStart(); 
        location.assign(location.href+"/"+comp);
    });
}); 
