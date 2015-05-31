$("#module-select-btn").click(function(){
    console.log("button clicked");
     var str = "";

     var tissue = $("ul#tissue-tab li.active").attr("tissue")
     modules = $("#"+tissue+"-module-select").val();
     str = modules.join([separator = '+']); 
    location.assign(location.href+"/"+tissue+"/"+str)
}); 


$(function(){
 $('#tissue-tab a:first').tab('show')
}); 
