var latestCounterEndpoint = "http://localhost:8080/latest-counter";
var style = "";

function getLatestCounterLog() {
    $.ajax( {
        url     : latestCounterEndpoint
        , success : function( content ){
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

getLatestCounterLog();

setInterval(function() {
    getLatestCounterLog();
}, 5 * 1000);

// FOR THE GRAPH

var l = 2000;
var y = 0;
var data = [];
var ds = { type: "line" };
var dataPoints = [];
var lastTenCounterEndpoint = "http://localhost:8080/last-ten-counter";

window.onload = function () {
    var chart = new CanvasJS.Chart("chart_graph", {
        zoomEnabled: true,
        title: {
            text: "Counter log",
        },
        axisY: {
            includeZero: false,
        },
        data: data,
    });
    chart.render();
};

$.ajax( {
    url     : lastTenCounterEndpoint
    , success : function( result ){
        for(var i = result.length - 1; i >= 0; i--){
            dataPoints.push({x: new Date(result[i].datetime), y: result[i].value });
        }
    }
    , error    : function(){
        alert( 'Failed to fetch result, please try again.' );
    }
} );

ds.dataPoints = dataPoints;
data.push(ds);    
