// Setup gauges
google.load("visualization", "1", {packages:["gauge", "corechart"]});
google.setOnLoadCallback(initVisualizations);

function initVisualizations() {
  initPhGauge();
  initTdsGauge();
  initWaterTempGauge();
  initChart();
  updateGauges();
}

function updateGauges() {
  var url = "http://gowebz.herokuapp.com/devices/MockClient4703/readings",
    phGauge = $('.ph-gauge'),
    tdsGauge = $('.tds-gauge'),
    wtempGauge = $('.wtemp-gauge'),
    chart = $('all-chart');

  var req = $.get(url);
  req.done(function(data){
    var timestamps = Object.keys(data);
    timestamps.sort().reverse();
    var first = timestamps[0];
    var readings = data[first];
    var ph = $('#ph-gauge').data('data');
    var tds = $('#tds-gauge').data('data');
    var wtemp = $('#wtemp-gauge').data('data');

    drawGauge($('#ph-gauge'), ph(readings.ph));
    drawGauge($('#tds-gauge'), tds(readings.tds));
    drawGauge($('#wtemp-gauge'), wtemp(readings.water_temperature));
    updateChart(data);
  });
  req.fail(function(){
    console.log("failed");
  });
  setTimeout(updateGauges, 5000);
}

function drawGauge(ele, data) {
  var gauge = $(ele).data("gauge");
  var options = $(ele).data("options");
  gauge.draw(data, options);
}

function updateChart(readings) {
  var ele = $('#all-chart');
  $(ele).data('readings', readings);
  var timestamps = Object.keys(readings);
  timestamps.sort();
  var data = $(ele).data('data');
  var chart = $(ele).data('chart');
  var options = $(ele).data('options');
  var arr = $.map(timestamps, function(v,i){
    return [[v, readings[v].ph]];
  });
  chart.draw(data(arr), options);
}

function initPhGauge() {
  var ele = $('#ph-gauge');
  var gauge = new google.visualization.Gauge(ele[0]);
  var data = function(ph){
    return google.visualization.arrayToDataTable([
      ['Label', 'Value'],
      ['PH', ph],
    ]);
  };
  var options = {
    min: 4,
    max: 7,
    greenFrom: 4, greenTo: 7, greenColor: "#FF0000",     //Red
    yellowFrom: 4.5, yellowTo: 6.5, yellowColor: "#FFFF00",     //Yellow
    redFrom: 5, redTo: 6, redColor: "#00FF00",     //Green
    majorTicks: [4,4.5,5,5.5,6,6.5,7],
    minorTicks: 5,
    animation:{duration: 2000, easing: 'linear'}
  };

  $(ele).data("gauge", gauge);
  $(ele).data("options", options);
  $(ele).data("data", data);

  gauge.draw(data(8.0), options);
}

function initTdsGauge() {
  var ele = $('#tds-gauge');
  var gauge = new google.visualization.Gauge(ele[0]);
  var data = function(tds){
    return google.visualization.arrayToDataTable([
      ['Label', 'Value'],
      ['TDS', tds],
    ]);
  };

  var options = {
    min: 0,
    max: 2000,
    greenFrom: 0, greenTo: 2000, greenColor: "#FF0000", //Red
    yellowFrom: 250, yellowTo: 1750, yellowColor: "#FFFF00", //Yellow
    redFrom: 500, redTo: 1500, redColor: "#00FF00", //Green
    majorTicks: ['0',500,1000,1500,2000],
    minorTicks: 5,
    animation:{duration: 2000, easing: 'linear'}
  };

  $(ele).data("gauge", gauge);
  $(ele).data("options", options);
  $(ele).data("data", data);

  gauge.draw(data(1000), options);
}

function initWaterTempGauge() {
  var ele = $('#wtemp-gauge');
  var gauge = new google.visualization.Gauge(ele[0]);
  var data = function(temp){
    return google.visualization.arrayToDataTable([
      ['Label', 'Value'],
      ['H20 Temp', temp],
    ]);
  }
  var options = {
    min: 50,
    max: 100,
    greenFrom: 50, greenTo: 100, greenColor: "#FF0000", //Red
    yellowFrom: 60, yellowTo: 90, yellowColor: "#FFFF00", //Yellow
    redFrom: 70, redTo: 80, redColor: "#00FF00", //Green
    majorTicks: [50,60,70,80,90,100],
    minorTicks: 10,
    animation:{duration: 2000, easing: 'linear'}
  };

  $(ele).data("gauge", gauge);
  $(ele).data("options", options);
  $(ele).data("data", data);

  gauge.draw(data(72),options)
}

function initChart() {
  var data = function(readings){
    var d = new google.visualization.DataTable();
    d.addColumn('string', 'Timestamp');
    d.addColumn('number', 'Ph');
    d.addRows(readings);
    return d;
  }

  var readings = [];

  var options = {
    hAxis: {
      title: 'Time'
    },
    vAxis: {
      title: 'Ph'
    }
  };
  var ele = $('#all-chart');
  var chart = new google.visualization.LineChart(ele[0]);

  $(ele).data("chart", chart);
  $(ele).data("options", options);
  $(ele).data('data', data);

  chart.draw(data(readings), options);
}
