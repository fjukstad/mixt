$("#module-select-btn").click(function(){
    console.log("button clicked");
     var str = "";
     modules = $("#module-select").val();
     str = modules.join([separator = '+']); 
    location.assign(location.href+"/"+str)
}); 
