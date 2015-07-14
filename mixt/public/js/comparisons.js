$(function(){
    // navigate to new page based on user action
    $("td#comp").click(function(){
        var comp = $(this).attr("comp");
        location.assign(location.href+"/"+comp);
    });
}); 
