// Setup gauges
google.load("visualization", "1", {packages:["gauge", "corechart"]});
google.setOnLoadCallback(initVisualizations);

function initVisualizations() {
  drawPhGauge();
  drawTdsGauge();
  drawWaterTempGauge();
  drawChart();
}

function drawChart() {

  var data = new google.visualization.DataTable();
  data.addColumn('number', 'X');
  data.addColumn('number', 'Dogs');

  data.addRows([
    [0, 0],   [1, 10],  [2, 23],
    [3, 17],  [4, 18],  [5, 9],
    [6, 11],  [7, 27],  [8, 33],
    [9, 40],  [10, 32], [11, 35],
    [12, 30], [13, 40], [14, 42],
    [15, 47], [16, 44], [17, 48],
    [18, 52], [19, 54], [20, 42],
    [21, 55], [22, 56], [23, 57],
    [24, 60], [25, 50], [26, 52],
    [27, 51], [28, 49], [29, 53],
    [30, 55], [31, 60], [32, 61],
    [33, 59], [34, 62], [35, 65],
    [36, 62], [37, 58], [38, 55],
    [39, 61], [40, 64], [41, 65],
    [42, 63], [43, 66], [44, 67],
    [45, 69], [46, 69], [47, 70],
    [48, 72], [49, 68], [50, 66],
    [51, 65], [52, 67], [53, 70],
    [54, 71], [55, 72], [56, 73],
    [57, 75], [58, 70], [59, 68],
    [60, 64], [61, 60], [62, 65],
    [63, 67], [64, 68], [65, 69],
    [66, 70], [67, 72], [68, 75],
    [69, 80]
  ]);

  var options = {
    hAxis: {
      title: 'Time'
    },
    vAxis: {
      title: 'Popularity'
    }
  };

  var chart = new google.visualization.LineChart(document.getElementById('all-chart'));
  chart.draw(data, options);
}

function drawPhGauge() {

  var gauge = new google.visualization.Gauge(document.getElementById('ph-gauge'))
  var data = google.visualization.arrayToDataTable([
    ['Label', 'Value'],
    ['PH', 5.5],
  ]);
  
  var options = {
	min: 4,
	max: 7,
	                                      //Red
	greenFrom: 4, greenTo: 7, greenColor: "#FF0000",
				              //Yellow
	yellowFrom: 4.5, yellowTo: 6.5, yellowColor: "#FFFF00",
			                      //Green		
	redFrom: 5, redTo: 6, redColor: "#00FF00",
	majorTicks: [4,4.5,5,5.5,6,6.5,7],
	minorTicks: 5,
	animation:{duration: 2000, easing: 'linear'}
  };
  gauge.draw(data, options);

}
function drawTdsGauge() {
  var gauge = new google.visualization.Gauge(document.getElementById('tds-gauge'))
  var data = google.visualization.arrayToDataTable([
    ['Label', 'Value'],
    ['TDS', 1000],
  ]);

  var options = {
	min: 0,
	max: 2000,
	                                      //Red
	greenFrom: 0, greenTo: 2000, greenColor: "#FF0000",
				              //Yellow
	yellowFrom: 250, yellowTo: 1750, yellowColor: "#FFFF00",
			                      //Green		
	redFrom: 500, redTo: 1500, redColor: "#00FF00",
	majorTicks: ['0',500,1000,1500,2000],
	minorTicks: 5,
	animation:{duration: 2000, easing: 'linear'}
  };

  gauge.draw(data, options);

}
function drawWaterTempGauge() {
  var gauge = new google.visualization.Gauge(document.getElementById('wtemp-gauge'));
  var data = google.visualization.arrayToDataTable([
    ['Label', 'Value'],
    ['H20 Temp', 75],
  ]);
  var options = {
	min: 50,
	max: 100,
	                                      //Red
	greenFrom: 50, greenTo: 100, greenColor: "#FF0000",
				              //Yellow
	yellowFrom: 60, yellowTo: 90, yellowColor: "#FFFF00",
			                      //Green		
	redFrom: 70, redTo: 80, redColor: "#00FF00",
	majorTicks: [50,60,70,80,90,100],
	minorTicks: 10,
	animation:{duration: 2000, easing: 'linear'}
  };
  gauge.draw(data,options)
}
