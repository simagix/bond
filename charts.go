/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * charts.go
 */

package bond

import (
	"html/template"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetChartTemplate returns HTML
func GetChartTemplate(chartType string) (*template.Template, error) {
	html := GetContentHTML()
	if chartType == T_MIGRATION_STATS {
		html += MigrationStatsChartHTML
	} else if chartType == T_MIGRATION_TIME {
		html += MigrationTimeChartHTML
	} else if chartType == T_CHUNK_SPLITS {
		html += ChunkSplitsChartHTML
	}
	html += `
  	<div id='bondChart' style="width: 100%; clear: left;"></div>
  
		</body></html>`

	return template.New("bond").Funcs(template.FuncMap{
		"ISODate": func(t *primitive.DateTime) string {
			if t == nil {
				return ""
			}
			layout := "2006-01-02T15:04:05Z"
			return t.Time().Format(layout)
		},
		"toSeconds": func(n float64) float64 {
			return n / 1000
		},
		"substr": func(str string, n int) string {
			return str[:n]
		},
		"epoch": func(d string, s string) int64 {
			dfmt := "2016-01-02T23:59:59"
			sdt, _ := time.Parse("2006-01-02T15:04:05", s+dfmt[len(s):])
			dt, _ := time.Parse("2006-01-02T15:04:05", d+dfmt[len(d):])
			return dt.Unix() - sdt.Unix()
		}}).Parse(html)
}

const (
	MigrationTimeChartHTML = `
	<script>
		setChartType();
	</script>
{{ if and (.Config.Actions) (gt (len .Config.Actions.BalancerRounds) 0) }}
<script>
	google.charts.load('current', {'packages':['corechart']});
	google.charts.setOnLoadCallback(drawChart);

	function drawChart() {
		var data = google.visualization.arrayToDataTable([
			['Date/Time', 'Execution Time'],

	{{range $i, $v := .Config.Actions.BalancerRounds}}
		[new Date("{{ISODate $v.Time}}"), {{$v.AverageExecutionTime}}],
	{{end}}
		]);
		// Set chart options
		var options = {
			'backgroundColor': { 'fill': 'transparent' },
			'title': '{{.Title}}',
			'hAxis': { slantedText: true, slantedTextAngle: 30 },
			'vAxis': {title: 'Millisecond', minValue: 0},
			'width': '100%',
			'height': 480,
			'titleTextStyle': {'fontSize': 20},
			'explorer': { actions: ['dragToZoom', 'rightClickToReset'] },
			'legend': { 'position': 'right' } };
		// Instantiate and draw our chart, passing in some options.
		var chart = new google.visualization.ColumnChart(document.getElementById('bondChart'));
		// var chart = new google.visualization.LineChart(document.getElementById('bondChart'));
		chart.draw(data, options);
	}
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`

	MigrationStatsChartHTML = `
	<script>
		setChartType();
	</script>
{{ if and (.Config.Actions) (gt (len .Config.Actions.BalancerRounds) 0) }}
<script>
	google.charts.load('current', {'packages':['corechart']});
	google.charts.setOnLoadCallback(drawChart);

	function drawChart() {
		var data = google.visualization.arrayToDataTable([
			['Date/Time', 'Chunks Moved', 'Error Count'],

	{{range $i, $v := .Config.Actions.BalancerRounds}}
		[new Date("{{ISODate $v.Time}}"), {{$v.TotalChunksMoved}}, {{$v.TotalErrors}}],
	{{end}}
		]);
		// Set chart options
		var options = {
			'backgroundColor': { 'fill': 'transparent' },
			'title': '{{.Title}}',
			'hAxis': { slantedText: true, slantedTextAngle: 30 },
			'vAxis': {title: 'Chunks Moved per Hour', minValue: 0},
			'width': '100%',
			'height': 480,
			'isStacked': true,
			'titleTextStyle': {'fontSize': 20},
			'explorer': { actions: ['dragToZoom', 'rightClickToReset'] },
			'legend': { 'position': 'right' } };
		// Instantiate and draw our chart, passing in some options.
		var chart = new google.visualization.ColumnChart(document.getElementById('bondChart'));
		// var chart = new google.visualization.LineChart(document.getElementById('bondChart'));
		chart.draw(data, options);
	}
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`

	ChunkSplitsChartHTML = `
	<script>
		setChartType();
	</script>
{{ if and (.Config.Changes) (gt (len .Config.Changes.Splits) 0) }}
<script>
	google.charts.load('current', {'packages':['corechart']});
	google.charts.setOnLoadCallback(drawChart);

	function drawChart() {
		var data = google.visualization.arrayToDataTable([
			['Date/Time', 'Splits'],

	{{range $i, $v := .Config.Changes.Splits}}
		[new Date("{{ISODate $v.Time}}"), {{$v.Total}}],
	{{end}}
		]);
		// Set chart options
		var options = {
			'backgroundColor': { 'fill': 'transparent' },
			'title': '{{.Title}}',
			'hAxis': { slantedText: true, slantedTextAngle: 30 },
			'vAxis': {title: 'Millisecond', minValue: 0},
			'width': '100%',
			'height': 480,
			'titleTextStyle': {'fontSize': 20},
			'explorer': { actions: ['dragToZoom', 'rightClickToReset'] },
			'legend': { 'position': 'right' } };
		// Instantiate and draw our chart, passing in some options.
		var chart = new google.visualization.ColumnChart(document.getElementById('bondChart'));
		// var chart = new google.visualization.LineChart(document.getElementById('bondChart'));
		chart.draw(data, options);
	}
</script>
{{else}}
<div align='center' class='btn'><span style='color: red'>no data found</span></div>
{{end}}`
)
