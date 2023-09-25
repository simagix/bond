/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * template.go
 */

package bond

import (
	"fmt"
	"sort"
)

const headers = `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Chen's Bond</title>
	<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
	<meta http-equiv="Pragma" content="no-cache" />
	<meta http-equiv="Expires" content="0" />

  <script src="https://www.gstatic.com/charts/loader.js"></script>
  <link href="/favicon.ico" rel="icon" type="image/x-icon" />
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
  <style>
    :root {
      --text-color: #FF0000;       /* Vibrant Red for Text */
      --header-color: #FFC984;     /* Bright Orange for Headers */
      --row-color: #FFF2CC;        /* Soft Yellow for Rows */
      --background-color: #F0F0F0; /* Neutral Gray for Background */
      --accent-color-1: #FF78A2;   /* Playful Pink for Accent 1 */
      --accent-color-2: #61C0BF;   /* Teal for Accent 2 */
      --accent-color-3: #FF3D57;   /* Vivid Red for Accent 3 */
      --border-color: #A6A6A6;     /* Light Gray for Borders */
    }
  	body {
      font-family: Helvetica, Arial, sans-serif;
      margin-top: 10px;
      margin-bottom: 10px;
      margin-right: 10px;
      margin-left: 10px;
	    background-color: var(--background-color);
    }
    div {
        margin: 5px 5px;
    }
    table {
      border-collapse:collapse;
      min-width: 400px;
      margin: 5px 5px;
    }
    caption {
      caption-side: top;
      font-size: 1.25em;
      font-weight: bold;
	    text-align: left;
    }
    table, th, td {
      border: 1px solid var(--border-color);
      vertical-align: middle;
    }
    th {
      background-color: var(--header-color);
      color: var(--text-color);
      font-weight: bold;
      padding: 0.3rem;
      font-size: 1em;
      text-align: center;
    }
    td {
      background-color: var(--row-color);
      padding: 0.1rem;
      font-size: 1em;
    }
    tr:nth-child(even) td {
      background-color: white;
    }
    .rowtitle {
      vertical-align: middle;
      font-size: 1em;
      font-weight: bold;
      word-break: break-all;
      padding: 5px 5px;
    }
    .break {
      vertical-align: middle;
      font-size: 1em;
      word-break: break-all;
      padding: 5px 5px;
    }
    table a:link {
      color: var(--text-color);
      text-decoration: none;
    }
    table a:visited {
      color: var(--header-color);
      text-decoration: none;
    }
    table a:hover {
      color: red;
      text-decoration: none;
    }
    ul, ol {
      #font-family: Consolas, monaco, monospace;
      font-size: 1em;
    }
    .btn {
      background-color: transparent;
      border: none;
      outline:none;
      color: var(--text-color);
      padding: 2px 2px;
      cursor: pointer;
      font-size: 1.5em;
      font-weight: bold;
      border-radius: .25em;
    }
    .btn:hover {
      background-color: var(--accent-color-2);
      color: #FFF;
    }
    .button { 
      background-color: var(--text-color);
      border: none; 
      outline: none;
      color: var(--background-color);
      padding: 3px 15px;
      margin: 0px 10px;
      cursor: pointer;
      font-size: 14px;
      font-weight: bold;
      border-radius: 3px;
    }
    .exclamation {
      background: none;
      color: red;
      border: none;
      outline: none;
      padding: 5px 10px;
      margin: 2px 2px;
      font-size: 1em;
      border-radius: .25em;
    }
    .tooltip {
      position: relative;
      display: inline-block;
    }
    .tooltip .tooltiptext {
      visibility: hidden;
      width: 200px;
      background-color: #555;
      color: #fff;
      text-align: left;
      border-radius: 6px;
      padding: 5px 5px;
      position: absolute;
      z-index: 1;
      bottom: 125%;
      left: 50%;
      margin-left: -100px;
      opacity: 0;
      transition: opacity 0.3s;
    }
    .tooltip .tooltiptext::after {
      content: "";
      position: absolute;
      top: 100%;
      left: 50%;
      margin-left: -5px;
      border-width: 5px;
      border-style: solid;
      border-color: #555 transparent transparent transparent;
    }
    .tooltip:hover .tooltiptext {
      visibility: visible;
      opacity: 1;
    }
    h1 {
      font-size: 1.6em;
      font-weight: bold;
    }
    h2 {
      font-size: 1.4em;
      font-weight: bold;
    }
    h3 {
      font-size: 1.2em;
      font-weight: bold;
    }
    h4 {
      font-size: 1em;
      font-weight: bold;
    }
    .footer {
      background-color: #fff;
      opacity: .75;
      position: fixed;
      left: 0;
      bottom: 0;
      width: 100%;
      color: #000;
      text-align: left;
      padding: 2px 10px;
    }
    input, select, textarea {
      font-family: "Trebuchet MS";
      appearance: auto;
      background-color: var(--row-color);
      color: var(--text-color);
      border-radius: .25em;
      font-size: .9em;
      #padding: 5px 5px;
    }
    .rotate23:hover {
      -webkit-transform: rotate(23deg);
      -moz-transform: rotate(23deg);
      -o-transform: rotate(23deg);
      -ms-transform: rotate(23deg);
      transform: rotate(23deg);
    }
    input[type="checkbox"] {
      accent-color: red;
    }
    .sort {
      color: #FFF;
    }
    .sort:hover {
      color: #DB4437;
    }
    .summary {
      font-family: Consolas, monaco, monospace;
	  background-color: #001f3f;
      color: var(--row-color);
	  padding: .5rem;
	  margin: .5rem;
      // font-size: .8em;
    }
    #loading {
      position: fixed;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      background-color: rgba(0, 0, 0, 0.5);
      z-index: 9999;
	  display: none;
    }
    .spinner {
      border: 5px solid #f3f3f3;
      border-top: 5px solid #3498db;
      border-radius: 50%;
      width: 50px;
      height: 50px;
      animation: spin 2s linear infinite;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    }
    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
    .chart {
      background-color: var(--row-color);;
      border: solid;
      padding: 10px 10px;
      margin: 10px 10px;
      border-radius: .5em;
    }
  </style>
  <script>
    function loadData(url) {
    	var loading = document.getElementById('loading');
    	loading.style.display = 'block';
    	fetch(url)
        	.then(response => response.text())
        	.then(data => {
      			loading.style.display = 'none';
				document.open();
				document.write(data);
				document.close();
        	})
        	.catch(error => {
      			loading.style.display = 'none';
        	});
    }
  </script>
</head>
<body>
  <div id="loading">
    <div class="spinner"></div>
  </div>
`

func GetContentHTML() string {
	html := headers
	html += `
<script>
	function gotoChart() {
		var sel = document.getElementById('nextChart')
		var value = sel.options[sel.selectedIndex].value;
		if(value == "") {
			return;
		}
		loadData(value);
	}
</script>
<div align='center' style='margin: 5px 5px; padding: 5px 5px;'>
  <div style="float: left;">
  	<button id="logs" onClick="javascript:loadData('/bond/info'); return false;"
		class="btn"><i class="fa fa-home"> Chen's Bond</i></button>
  </div>

  <div style="float: left;">
    <button id="chart" onClick="javascript:loadData('/bond/charts/migration_time'); return false;" 
        class="btn" style=""><i class="fa fa-bar-chart"></i></button>
    <select id='nextChart' style="margin: 0px 0px; padding: 0px 0px; font-size: 1em" onchange='gotoChart()'>`
	items := []Chart{}
	for _, chart := range charts {
		items = append(items, chart)
	}
	sort.Slice(items, func(i int, j int) bool {
		return items[i].Index < items[j].Index
	})

	html += "<option value=''>select a chart</option>"
	for i, item := range items {
		if i == 0 {
			continue
		}
		html += fmt.Sprintf("<option value='%v'>%v</option>", item.URL, item.Title)
	}

	html += `</select>
  </div>
</div>
<script>
	function setChartType() {
		var sel = document.getElementById('nextChart')
		sel.selectedIndex = {{.Chart.Index}};
	}
</script>
`
	//	html += getFooter()
	return html
}
