/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * info.go
 */

package bond

import (
	"fmt"
	"html/template"
	"math/rand"
	"strings"
	"time"

	"github.com/simagix/gox"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// GetInfoTemplate returns HTML
func GetInfoTemplate() (*template.Template, error) {
	html := GetContentHTML()
	html += `<div style='margin: 5px 5px; width=100%; clear: left;'>
	  <table style='border: none; margin: 10px 10px; width=100%; clear: left;' width='100%'>
		{{$flag := coinToss}}
		<tr><td style='border:none; vertical-align: top; padding: 5px; background-color: var(--background-color);'>
			<img class='rotate23' src='data:image/png;base64,{{ assignConsultant $flag }}'></img></td>
			<td class='summary'>{{consultantIntro $flag}} ` + SummaryHTML + "</td></tr></table></div>"
	html += InfoHTML + LogsHTML + BondVideoHTML + MongosHTML + ShardsHTML + DatabasesHTML + CollectionsHTML
	html += "</body></html>"
	return template.New("bond").Funcs(template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
		"assignConsultant": func(sage bool) string {
			if sage {
				return SAGE_PNG
			}
			return SIMONE_PNG
		},
		"checkVersion": func(version string, majorVersion string) bool {
			toks := strings.Split(version, ".")
			if len(toks) < 2 {
				return false
			}
			return (strings.Join(toks[:2], ".") == majorVersion)
		},
		"coinToss": func() bool {
			rand.Seed(time.Now().UnixNano())
			randomNum := rand.Intn(2)
			return (randomNum%2 == 0)
		},
		"consultantIntro": func(sage bool) string {
			if sage {
				return "Hello, my name is Sage and I like to share my thoughts with you on the findings."
			}
			return "Hey there, my name is Simone, and here is the summary I have prepared for you."
		},
		"getCheckMarkSymbol": func(b bool) template.HTML {
			if !b {
				return template.HTML("")
			}
			return template.HTML("<i class='fa fa-check'></i>")
		},
		"getCountLabel": func(n int, prefix string) string {
			printer := message.NewPrinter(language.English)
			if n > TOP_N {
				return printer.Sprintf("%s %d of %d", prefix, TOP_N, n)
			}
			return printer.Sprintf("%d", n)
		},
		"getDurationFromSeconds": func(s int64) string {
			return gox.GetDurationFromSeconds(float64(s))
		},
		"getDurationFromMilliseconds": func(n interface{}) string {
			if n == nil {
				return ""
			}
			ms := ToInt64(n)
			if ms == 0 {
				return ""
			} else if ms < 1000 {
				return fmt.Sprintf("%v ms", ms)
			}
			return gox.GetDurationFromSeconds(float64(ms) / 1000)
		},
		"getHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"getStorageSize": func(size int64) string {
			return gox.GetStorageSize(float64(size))
		},
		"getUserSymbol": func(b bool) template.HTML {
			if !b {
				return template.HTML("")
			}
			return template.HTML("<i class='fa fa-user'></i>")
		},
		"getWarningSymbol": func(b bool) template.HTML {
			if b {
				return template.HTML("")
			}
			return template.HTML("<i class='fa fa-warning' style='color:red;'></i>")
		},
		"hasData": func(data []interface{}) bool {
			return len(data) > 0
		},
		"hasPrefix": func(str string, pre string) bool {
			return strings.HasPrefix(str, pre)
		},
		"ISODate": func(t *primitive.DateTime) string {
			if t == nil {
				return ""
			}
			layout := "2006-01-02T15:04:05Z"
			return t.Time().Format(layout)
		},
		"numPrinter": func(n interface{}) string {
			printer := message.NewPrinter(language.English)
			return printer.Sprintf("%v", ToInt(n))
		},
		"plural": func(num int, word string, tail string) string {
			if num == 0 {
				return "no " + word
			} else if num == 1 {
				return fmt.Sprintf("1 %s", word)
			}
			printer := message.NewPrinter(language.English)
			return printer.Sprintf("%d %s%s", num, word, tail)
		},
		"stringify": func(doc interface{}) string {
			return Stringify(doc)
		}}).Parse(html)
}

const (
	// Summary
	SummaryHTML = `
	{{$ndbs:=len .Config.Databases}}
	{{$ncolls:=len .Config.CollectionsMap}}
	Based on the data from Chen's Bond, the MongoDB cluster is running version {{.Config.MongoVersion}}.  It consists of {{plural (len .Config.ShardsMap) "shard" "s"}} and {{plural (len .Config.Mongos) "mongos instance" "s"}}. There are {{plural $ncolls "sharded collection" "s"}} across {{plural $ndbs "database" "s"}}.
	<p/>Please take a moment to review the statistics below.  All misconfigurations or anomalies will be highlighted with <i class="fa fa-warning" style="color: red;"></i> icons. 
	{{if eq (len .Config.Warnings) 0}}
	Good news! Bond didn't find anything goofy from your cluster. But, still,
	{{else}}
	Bond has uncovered a few noteworthy findings:
	<ol>
	{{range $n, $value := .Config.Warnings}}
		<li>{{getHTML $value}}</li>
	{{end}}
	</ol>
	Be sure to
	{{end}}
	review the provided statistics and don't forget to check out the charts to better understand the performance of your cluster.
	`

	// general info
	InfoHTML = `<div style='float: left;'>
		<table width=400px><caption>General Info</caption><tr><th>Metric</th><th>Value</th>
			<tr><td align='left' class='rowtitle'>Mongo Version</td><td align='center' class='break'>
				{{ .Config.MongoVersion }} {{getUserSymbol .Config.IsUserVersion}} {{getWarningSymbol (not .Config.IsUpgrade) }}</td></tr>
			<tr><td align='left' class='rowtitle'>Number of Shards</td><td align='right' class='break'>{{ len .Config.ShardsMap }}</td></tr>

		{{if gt (len .Config.Mongos) 0}}
			<tr><td align='left' class='rowtitle'>Number of mongos Found</td><td align='right' class='break'>{{ len .Config.Mongos }}</td></tr>
		{{end}}

			<tr><td align='left' class='rowtitle'>Number of Databases</td><td align='right' class='break'>{{ numPrinter (len .Config.Databases) }}</td></tr>
			<tr><td align='left' class='rowtitle'>Number of Sharded Collections</td><td align='right' class='break'>{{ numPrinter (len .Config.CollectionsMap) }}</td></tr>

		{{if ne .Actionlog.Capped nil}}
			<tr><td align='left' class='rowtitle'>Total Chunks Moved</td><td align='right' class='break'>{{ numPrinter .Actionlog.TotalChunksMoved }}</td></tr>
			<tr><td align='left' class='rowtitle'>Total Chunk Move Error</td><td align='right' class='break'>{{ .Actionlog.TotalErrors }}</td></tr>
			<tr><td align='left' class='rowtitle'>Average Chunk Move Time</td><td align='right' class='break'>{{ getDurationFromMilliseconds (.Actionlog.AverageExecutionTime) }}</td></tr>
			<tr><td align='left' class='rowtitle'>Longest Chunk Move Time</td><td align='right' class='break'>{{ getDurationFromMilliseconds (.Actionlog.MaxExecutionTime) }}</td></tr>
		{{end}}

		{{if ne .Changelog.Capped nil}}
			<tr><td align='left' class='rowtitle'>Total Chunk Splits</td><td align='right' class='break'>{{ numPrinter .Changelog.TotalSplits }}</td></tr>
		{{end}}
		</table></div>`

	// logs
	LogsHTML = `<div style='float: left;'>
		<table width=400px><caption>Action and Change Logs</caption><tr><th>Metric</th><th>Value</th>
		{{if ne .Actionlog.Capped nil}}
			<tr><td align='left' class='rowtitle'>Is config.actionlog capped?</td><td align='center' class='break'>{{getCheckMarkSymbol .Actionlog.Capped}}{{getWarningSymbol .Actionlog.Capped}}</td></tr>
			<tr><td align='left' class='rowtitle'>Size of config.actionlog</td><td align='right' class='break'>{{getStorageSize .Actionlog.MaxSize}}</td></tr>
		{{end}}

		{{if ne .Changelog.Capped nil}}
			<tr><td align='left' class='rowtitle'>Is config.changelog capped?</td><td align='center' class='break'>{{getCheckMarkSymbol .Changelog.Capped}}{{getWarningSymbol .Changelog.Capped}}</td></tr>
			<tr><td align='left' class='rowtitle'>Size of config.changelog</td><td align='right' class='break'>{{getStorageSize .Changelog.MaxSize}}</td></tr>
		{{end}}
		</table></div>`

	// shards
	ShardsHTML = `
{{if gt (len .Config.ShardsMap) 0}}
	<div style='float: left;'>
	<table><caption>Shards</caption><tr><th>#</th>
	<th>Shard Name</th><th>Host</th><th>State</th><th>Chunks</th><th>Jumbo</th><th>Max Size</th><th>-</th>
	{{$cnt:=0}}
	{{range $n, $value := .Config.ShardsMap}}
		{{$cnt = add $cnt 1}}
			<tr>
				<td align='right' class='break'>{{ $cnt }}</td>
				<td align='left' class='break'>{{$value.ID}}</td>
				<td align='left' class='break'>{{ $value.Host }}</td>
				<td align='right' class='break'>{{ $value.State }}</td>
				<td align='right' class='break'>{{ numPrinter $value.Chunks }}</td>
				<td align='right' class='break'>{{ numPrinter $value.Jumbo }}</td>
			{{if not $value.MaxSize}}
				<td align='right' class='break'></td>
			{{else}}
				<td align='right' class='break'>{{ $value.MaxSize }} {{getWarningSymbol false}}</td>
			{{end}}
				<td align='center'><button class='btn'
					onClick="javascript:loadData('/bond/chart/shards/{{$value.ID}}'); return false;">
					<i class='fa fa-pie-chart' style='font-size: .8em;'></i></button></td>
			</tr>
	{{end}}
{{end}}
	</table></div>`

	// mongos
	MongosHTML = `
	{{if gt (len .Config.Mongos) 0}}
		<div style='float: left;'>
			<table><caption>mongos Instances</caption><tr><th>#</th>
				<th>_id</th><th>Version</th><th>Ping</th><th>Up</th><th>Waiting</th><th>Created</th>
		{{$majorVersion := .Config.MajorVersion}}
		{{range $n, $value := .Config.Mongos}}
				<tr>
					<td align='right' class='break'>{{ add $n 1 }}</td>
					<td align='left' class='break'>{{ $value.ID }}</td>
				{{ if (checkVersion $value.MongoVersion $majorVersion) }}
					<td align='center' class='break'>{{$value.MongoVersion}}</td>
				{{ else }}
					<td align='center' class='break'>{{$value.MongoVersion}} {{getWarningSymbol false}}</td>
				{{ end }}
					<td align='left' class='break'>{{ ISODate $value.Ping }}</td>
					<td align='right' class='break'>{{ getDurationFromMilliseconds $value.Up }}</td>
				{{ if $value.Waiting }}
					<td align='center' class='break'>{{getCheckMarkSymbol true}}</td>
				{{ else }}
					<td align='center' class='break'>{{getWarningSymbol false}}</td>
				{{ end }}
					<td align='left' class='break'>{{ ISODate $value.Created }}</td>
				</tr>
		{{end}}
	{{end}}
		</table></div>`

	// Databases
	DatabasesHTML = `
{{if gt (len .Config.Databases) 0}}
	<div style='float: left;'>
	<table><caption>{{getCountLabel (len .Config.Databases) "First"}} Databases</caption><tr><th>#</th>
	<th>Database Name</th><th>Primary</th><th>Partitioned</th>
	{{$limit:=.TOP_N}}
	{{range $n, $value := .Config.Databases}}
		{{if lt $n $limit}}
			<tr>
				<td align='right' class='break'>{{ add $n 1 }}</td>
				<td align='left' class='break'>{{ $value.ID }}</td>
				<td align='left' class='break'>{{ $value.Primary }}</td>
				<td align='center' class='break'>{{ getCheckMarkSymbol $value.Partitioned }}</td>
			</tr>
		{{end}}
	{{end}}
{{end}}
	</table></div>`

	// Collections
	CollectionsHTML = `
{{if gt (len .Collections) 0}}
	<div style='float: left;'>
	<table><caption>{{getCountLabel (len .Collections) "Top"}} Sharded Collections</caption><tr><th>#</th>
	<th>Collection Name</th><th>Shard Key</th><th>Unique</th><th>No Balance</th><th>Chunks</th><th>-</th>
	{{$limit:=.TOP_N}}
	{{range $n, $value := .Collections}}
		{{if lt $n $limit}}
			<tr>
				<td align='right' class='break'>{{ add $n 1 }}</td>
				<td align='left' class='break'>{{$value.ID}}</td>
				<td align='left' class='break'>{{ stringify $value.Key }}</td>
				<td align='center' class='break'>{{ getCheckMarkSymbol $value.Unique }}</td>
				<td align='center' class='break'>{{ getCheckMarkSymbol $value.NoBalance }}</td>
				<td align='right' class='break'>{{ numPrinter $value.Chunks }}</td>
				<td align='center'><button class='btn'
					onClick="javascript:loadData('/bond/chart/namespaces/{{$value.ID}}'); return false;">
					<i class='fa fa-pie-chart' style='font-size: .8em;'></i></button></td>
			</tr>
		{{end}}
	{{end}}
{{end}}
	</table></div>`

	BondVideoHTML = `
	<div style='float: left;'>
	<table><caption>Bond Tutorial</caption>
	<tr><td>
	<iframe width="480" height="270" src="https://www.youtube.com/embed/tdUT_PLiRSs?si=q-WOmdGqwgJV8R0R" 
		title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; 
		picture-in-picture; web-share" allowfullscreen></iframe></td></tr></table></div>
	`
)
