window.onload = function() {
     refreshPage();
     setInterval(refreshPage, 5000);}

function refreshPage()
{
	gaugeRequest = serverRequest();
	var URL = "http://localhost:5000/readings/latest";
	gaugeRequest.open("GET", URL, true);
    gaugeRequest.onreadystatechange = function()
    {
		if (gaugeRequest.readyState == 4 && gaugeRequest.status == 200)
      {
        JSONData = JSON.parse(gaugeRequest.responseText);
        var waterTempData = JSONData.temperature;
        var phData = JSONData.ph;
        var tdsData = JSONData.tds;
        var Time = new Date(JSONData.created_at).format("H:i");

        updateCharts(Time, waterTempData, phData, tdsData);
        updateGauges(waterTempData, phData, tdsData);
      }
	}
    gaugeRequest.send(null);
}

function serverRequest()
{
	try { var newRequest = new XMLHttpRequest(); }
	catch (requestError){
		alert("Exception caught in Charts.serverRequest()");
		var newRequest=null;}
	return(newRequest)
}
