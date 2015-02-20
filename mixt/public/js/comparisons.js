$(function(){
    // navigate to new page based on user action
    $("td#comp").click(function(){
        var comp = $("td#comp").attr("comp");
        location.assign(location.href+"/"+comp);
    });
}); 
