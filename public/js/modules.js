// Open first tissue tab 
$("#tissue-tab li:eq(0) a").tab('show'); 



$("#module-select-btn").click(function(){
    console.log("button clicked");
     var str = "";

     var tissue = $("ul#tissue-tab li.active").attr("tissue")
     modules = $("#"+tissue+"-module-select").val();
     str = modules.join([separator = '+']); 

    cohort= $("#cohort-select").val();

    location.assign(location.href+"/"+tissue+"/"+str+"/cohort/"+cohort)
}); 


$("#module-download-btn").click(function(){
    var str = "";

    var tissue = $("ul#tissue-tab li.active").attr("tissue")
    modules = $("#"+tissue+"-module-select").val();

    baseUrl = location.href+"/"+tissue+"/"; 
    str = modules.join([separator = '+']); 
    window.open(baseUrl + str + "/genes");
}); 


//// sigma instance
//var s, filter; 
//
//$(function(){
// $('#tissue-tab a:first').tab('show')
//    var h = $("#module-selector").height();
//    var w = $("#module-selector").width();
// 
//    $("#graph-container").height(h)
//    $("#graph-container").width(w)
//
//    filter = new sigma.plugins.filter(s);
//    var tissue = $("ul#tissue-tab li.active").attr("tissue")
//    TOMGraph(s, tissue); 
//
//
//}); 


//// module selected 
//$("select").click(function(){
//    filter.undo().apply();
//        
//     var tissue = $("ul#tissue-tab li.active").attr("tissue")
//    modules = $("#"+tissue+"-module-select").val();
//     if(modules != null){ 
//        filter.nodesBy(function(n){
//            if($.inArray(n.module, modules) != -1){
//                return true
//            } else {
//                return false
//            }
//        }).apply();
//     }
//})
//
//
//
//// select tissue 
//$("ul#tissue-tab li").mouseup(function(){
//    setTimeout(function(){
//        var tissue = $("ul#tissue-tab li.active").attr("tissue")
//        TOMGraph(s, tissue) 
//    }, 10);
//})


