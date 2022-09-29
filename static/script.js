var lastDatetime = "";
var latestCounterEndpoint = "http://localhost:8080/latest-counter";
var latestCounter = $("#latest_counter");
var style = "";

function getLatestCounterLog() {
    console.log("Hello");
    $.ajax( {
        url     : latestCounterEndpoint
      , success : function( content ){

        console.log(content);

        lastDatetime = content.datetime
        if( content.value >= 5 ){
            style = "background:green;";
        }else{
            style = "background:red;";
        }
        $("#latest_counter").html(content.value).attr("style",style);
      }
      , error    : function(){
          alert( 'Failed to fetch result, please try again.' );
      }
    } );
}

setInterval(function() {
    getLatestCounterLog();
}, 5 * 1000); 
