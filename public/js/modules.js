// Open first tissue tab 
$("#tissue-tab li:eq(0) a").tab('show'); 

$("#module-select-btn").click(function(){
    var str = "";

    var tissue = $("ul#tissue-tab li.active").attr("tissue")
    modules = $("#"+tissue+"-module-select").val();
    str = modules.join([separator = '+']); 

    cohort= $("#cohort-select").val();

    loadStart();

    location.assign(location.href+"/"+tissue+"/"+str+"/cohort/"+cohort)
}); 


$("#module-download-btn").click(function(){
    var str = "";

    loadStart();
    var tissue = $("ul#tissue-tab li.active").attr("tissue")
    modules = $("#"+tissue+"-module-select").val();

    baseUrl = location.href+"/"+tissue+"/"; 
    str = modules.join([separator = '+']); 
    window.open(baseUrl + str + "/genes");
    loadStop();
}); 

