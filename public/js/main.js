$(function() {
    var panel = $('.alternate').scotchPanel({
        clickSelector: '.toggle-panels',        
        containerSelector: '.main', 
        direction: 'left',
        distanceX: '10%',
    });


    $("#settings-icon").mouseenter(function(){
        panel.open()
    }); 

    $("#settings-icon").click(function(){
        panel.toggle()
    }); 

}); 
