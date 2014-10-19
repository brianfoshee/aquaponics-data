window.onload = function() {
     refreshPage();
     setInterval(refreshPage, 60000);}

function refreshPage() {
	gaugeRequest = serverRequest();
	var URL = "http://gowebz.heroku.com/devices/fe5d2287d64e448971febd91a4229b616e3a4bf6641916ade12a931a/readings?number=1";
	gaugeRequest.open("GET", URL, true);
    	gaugeRequest.onreadystatechange = function() {
    	if (gaugeRequest.readyState == 4 && gaugeRequest.status == 200) {
    		JSONData = JSON.parse(gaugeRequest.responseText);
    		var waterTempData = JSONData[0].temperature;
    		var phData = JSONData[0].ph;
    		var tdsData = JSONData[0].tds;
    		var Time = new Date(JSONData[0].created_at).format("H:i");

    		updateCharts(Time, waterTempData, phData, tdsData);
    		updateGauges(waterTempData, phData, tdsData);
    		}
    	}
    gaugeRequest.send(null);
}

function serverRequest() {
	try { 
		var newRequest = new XMLHttpRequest(); 
	}
	catch (requestError){
		alert("Exception caught in Charts.serverRequest()");
		var newRequest=null;
	}
	return(newRequest)
}
