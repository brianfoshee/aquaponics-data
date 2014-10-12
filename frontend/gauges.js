google.load('visualization', '1', {packages: ['gauge']});

var waterTempGauge;
var waterTempGaugeOptions;
var waterTempGaugeData;
var idealWaterTemp = 75;
var minWaterTemp = idealWaterTemp - 30;
var maxWaterTemp = idealWaterTemp + 30
var minWarningWaterTemp = idealWaterTemp - 10;
var minCriticalWaterTemp = idealWaterTemp - 20;
var maxWarningWaterTemp = idealWaterTemp + 10;
var maxCriticalWaterTemp = idealWaterTemp + 20;

var phGauge;
var phGaugeOptions;
var phGaugeData;
var idealPH = 6;
var minPH = idealPH - 1.5;
var maxPH = idealPH + 1.5;
var minCriticalPH = idealPH - 1;
var minWarningPH = idealPH - 0.5;
var maxWarningPH = idealPH + 0.5;
var maxCriticalPH = idealPH + 1;

var tdsGauge;
var tdsGaugeOptions;
var tdsGaugeData;
var idealTDS = 800;
var minTDS = idealTDS - 300;
var maxTDS = idealTDS + 300;
var minWarningTDS = idealTDS - 100;
var minCriticalTDS = idealTDS - 200;
var maxWarningTDS = idealTDS + 100;
var maxCriticalTDS = idealTDS + 200;

google.setOnLoadCallback(drawGauges);
function drawGauges()
{
	waterTempGauge = new google.visualization.Gauge(document.getElementById('waterTempGauge'));
    waterTempGaugeOptions = {
    		min: minWaterTemp,
    		max: maxWaterTemp,
    		greenFrom: minWaterTemp, greenTo: maxWaterTemp, greenColor: "#FF0000",
    		redFrom: minWarningWaterTemp, redTo: maxWarningWaterTemp, redColor: "#00FF00",
    		yellowFrom: minCriticalWaterTemp, yellowTo: maxCriticalWaterTemp, yellowColor: "#FFFF00",
    		majorTicks: [minWaterTemp, minCriticalWaterTemp, minWarningWaterTemp, idealWaterTemp, maxWarningWaterTemp, maxCriticalWaterTemp, maxWaterTemp],
    		minorTicks: 5,
    		animation:{duration: 2000, easing: 'linear'}};
    waterTempGaugeData = google.visualization.arrayToDataTable([['H20 Temp'],[75]]);
	waterTempGauge.draw(waterTempGaugeData, waterTempGaugeOptions);

	phGauge = new google.visualization.Gauge(document.getElementById('phGauge'));
	phGaugeOptions = {
			min: minPH,
			max: maxPH,
			greenFrom: minPH, greenTo: maxPH, greenColor: "#FF0000",
			redFrom: minWarningPH, redTo: maxWarningPH, redColor: "#00FF00",
			yellowFrom: minCriticalPH, yellowTo: maxCriticalPH, yellowColor: "#FFFF00",
			majorTicks: [minPH, minCriticalPH, minWarningPH, idealPH, maxWarningPH, maxCriticalPH, maxPH],
			minorTicks: 5,
			animation:{duration: 2000, easing: 'linear'}};
    phGaugeData = google.visualization.arrayToDataTable([['PH'],[7]]);
	phGauge.draw(phGaugeData, phGaugeOptions);

	tdsGauge = new google.visualization.Gauge(document.getElementById('tdsGauge'));
	tdsGaugeOptions = {
			min: minTDS,
			max: maxTDS,
			greenFrom: minTDS, greenTo: maxTDS, greenColor: "#FF0000",
			redFrom: minWarningTDS, redTo: maxWarningTDS, redColor: "#00FF00",
			yellowFrom: minCriticalTDS, yellowTo: maxCriticalTDS, yellowColor: "#FFFF00",
			majorTicks: [minTDS, minCriticalTDS, minWarningTDS, idealTDS, maxWarningTDS, maxCriticalTDS, maxTDS],
			minorTicks: 5,
			animation:{duration: 2000, easing: 'linear'}};
	tdsGaugeData = google.visualization.arrayToDataTable([['TDS'],[250]]);
	tdsGauge.draw(tdsGaugeData, tdsGaugeOptions);
}

function updateGauges(waterTempData, phData, tdsData)
{
	waterTempGaugeData.setValue(0, 0, waterTempData);
	waterTempGauge.draw(waterTempGaugeData, waterTempGaugeOptions);

	phGaugeData.setValue(0, 0, phData);
	phGauge.draw(phGaugeData, phGaugeOptions);

	tdsGaugeData.setValue(0, 0, tdsData);
	tdsGauge.draw(tdsGaugeData, tdsGaugeOptions);
}

function updateGaugeOptions()
{
	//alert('updateGaugeOptions.js - Entering');
	var URL = "/envdata/getGaugeOptions?plantType=Aquaponics";
	var updateGaugeOptions = serverRequest();
	updateGaugeOptions.open("GET", URL, true);
	//alert('updateGaugeOptions Request Sent');
	updateGaugeOptions.onreadystatechange = function()
	{
		//alert('Waiting for updateGaugeOptions readyState==4 and status==200');
		if (updateGaugeOptions.readyState == 4 && updateGaugeOptions.status == 200)
		{
			//alert('ReadyState == 4 and status == 200');
			JSONData = JSON.parse(updateGaugeOptions.responseText);

		 	idealWaterTemp = JSONData[0].idealWaterTemp;
			idealWaterTemp = 75;
		 	minWaterTemp = idealWaterTemp - 30;
		 	maxWaterTemp = idealWaterTemp + 30
		 	minWarningWaterTemp = idealWaterTemp - 10;
		 	minCriticalWaterTemp = idealWaterTemp - 20;
		 	maxWarningWaterTemp = idealWaterTemp + 10;
		 	maxCriticalWaterTemp = idealWaterTemp + 20;

			idealPH = JSONData[0].idealPH;
			minPH = idealPH - 1.5;
			maxPH = idealPH + 1.5;
			minCriticalPH = idealPH - 1;
			minWarningPH = idealPH - 0.5;
			maxWarningPH = idealPH + 0.5;
			maxCriticalPH = idealPH + 1;

			idealTDS = JSONData[0].idealTDS;
			minTDS = idealTDS - 600;
			maxTDS = idealTDS + 600;
			minWarningTDS = idealTDS - 200;
			minCriticalTDS = idealTDS - 400;
			maxWarningTDS = idealTDS + 200;
			maxCriticalTDS = idealTDS + 400;

		 	drawGauges();
		 	alert('Draw Gauges() Called');

			refreshPage();
			alert('refreshPage() Called');
		}
	}
	updateGaugeOptions.send(null);
}
