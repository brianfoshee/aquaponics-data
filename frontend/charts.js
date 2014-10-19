google.load('visualization', '1', {packages: ['corechart']});

var waterTempData = [['created_at', 'temperature']];
var waterTempChart;
var waterTempChartData;
var waterTempChartOptions;

var phData = [['Time', 'PH']];
var phChart;
var phChartData;
var phChartOptions;

var tdsData = [['created_at', 'tds']];
var tdsChart;
var tdsChartData;
var tdsChartOptions;

var URL = 'http://gowebz.herokuapp.com/readings';
var requestEnvironmentData = serverRequest();
requestEnvironmentData.open("GET", URL, false);

requestEnvironmentData.onload = function()
{
	if (requestEnvironmentData.status == 200)
		{
			//alert(requestEnvironmentData.responseText);
			environmentData = JSON.parse(requestEnvironmentData.responseText);

			if (environmentData != null)
			{
				for (var i = 0; i < environmentData.length; i++)
				{
					var waterTemp = [[new Date(environmentData[i].created_at).format("H:i"), environmentData[i].temperature]];
					var ph = [[new Date(environmentData[i].created_at).format("H:i"), environmentData[i].ph]];
					var tds = [[new Date(environmentData[i].created_at).format("H:i"), environmentData[i].tds]];

					waterTempData.push.apply(waterTempData, waterTemp);
					phData.push.apply(phData, ph);
					tdsData.push.apply(tdsData, tds);
				}
			}
			else
			{
				waterTempData.push.apply(waterTempData, [['Waiting for data', 70]]);
				phData.push.apply(phData,[['Waiting for data', 7]]);
				tdsData.push.apply(tdsData,[['Waiting for data', 500]]);
			}
		}

};
requestEnvironmentData.send(null);

google.setOnLoadCallback(drawCharts);
function drawCharts()
{
    //Draw Water Temperature chart
    waterTempChartData = google.visualization.arrayToDataTable(waterTempData);
    waterTempChart = new google.visualization.LineChart(document.getElementById('waterTempChart'));
    waterTempChartOptions = {
    		title: 'Water Temperature',
    		animation: {duration: 20000, easing: 'linear'},
    		curveType: 'function',
    		legend: {position: 'none'},
    		vAxis: {minValue: 32, maxValue: 110}
    		};
    waterTempChart.draw(waterTempChartData, waterTempChartOptions);

    // Draw PH Chart
    phChartData = google.visualization.arrayToDataTable(phData);
    phChart = new google.visualization.LineChart(document.getElementById('phChart'));
    phChartOptions = {
    		title: 'PH',
    		animation: {duration: 20000, easing: 'linear'},
    		curveType: 'function',
    		legend: {position: 'none'},
    		vAxis: {maxValue: 10, minValue: 5}
    		};
    phChart.draw(phChartData, phChartOptions);

    //Draw TDS Chart
    tdsChartData = google.visualization.arrayToDataTable(tdsData);
    tdsChart = new google.visualization.LineChart(document.getElementById('tdsChart'));
    tdsChartOptions = {
    	title: 'Total Dissolved Solids',
    	animation: {duration: 20000, easing: 'linear'},
    	curveType: 'function',
    	legend: {position: 'none'},
    	vAxis: {minValue: 0}
    };
    tdsChart.draw(tdsChartData, tdsChartOptions);

}

function updateCharts(Time, waterTempData, phData, tdsData)
{

	waterTempRecords = waterTempChartData.getNumberOfRows();
	waterTempLastUpdate = waterTempChartData.getValue(waterTempRecords-1,0);
	if (waterTempLastUpdate == 'Waiting for data'){
		waterTempChartData.remove(0);}
	if (Time != waterTempLastUpdate){
		waterTempChartData.addRow([Time, waterTempData]);
		waterTempChart.draw(waterTempChartData, waterTempChartOptions);}

	phRecords = phChartData.getNumberOfRows();
	phLastUpdate = phChartData.getValue(phRecords-1,0);
	if (phLastUpdate == 'Waiting for data'){
		phChartData.remove(0);}
	if (Time != phLastUpdate){
		phChartData.addRow([Time, phData]);
		phChart.draw(phChartData, phChartOptions);}

	tdsRecords = tdsChartData.getNumberOfRows();
	tdsLastUpdate = tdsChartData.getValue(tdsRecords-1,0);
	if (tdsLastUpdate == 'Waiting for data'){
		tdsChartData.remove(0);}
	if (Time != tdsLastUpdate){
		tdsChartData.addRow([Time, tdsData]);
		tdsChart.draw(tdsChartData, tdsChartOptions);}
}

function serverRequest()
{
	try { var newRequest = new XMLHttpRequest(); }
	catch (requestError)
	{
		alert("Exception caught in Charts.serverRequest()");
		var newRequest=null;
	}
	return(newRequest)
}
