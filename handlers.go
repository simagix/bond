/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * handlers.go
 */

package bond

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
)

const (
	T_MIGRATION_TIME  = "migration_time"
	T_MIGRATION_STATS = "migration_stats"
	T_CHUNK_SPLITS    = "splits"
)

var charts = map[string]Chart{
	"instruction": {0, "select a chart", "", ""},
	T_MIGRATION_TIME: {1, "Migration Time",
		"Display average migration time", "/bond/charts/" + T_MIGRATION_TIME},
	T_MIGRATION_STATS: {2, "Migration Stats",
		"Display average migration time", "/bond/charts/" + T_MIGRATION_STATS},
	T_CHUNK_SPLITS: {3, "Chunk Splits",
		"Display chunk splits", "/bond/charts/" + T_CHUNK_SPLITS},
}

type Chart struct {
	Index int
	Title string
	Descr string
	URL   string
}

// ChartsHandler responds to API calls
func ChartsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/* APIs
	 * /bond/charts/:attr
	 */
	config := GetConfigDB()
	attr := params.ByName("attr")
	templ, err := GetChartTemplate(attr)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err})
		return
	}
	title := charts[attr].Title
	doc := map[string]interface{}{"Chart": charts[attr], "Config": config, "Title": title}
	if err = templ.Execute(w, doc); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
		return
	}
}

// DataHandler responds to API calls
func DataHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/* APIs
	 * /api/bond/v1.0/data/:attr
	 */
	attr := params.ByName("attr")
	config := GetConfigDB()
	if attr == "info" {
		json.NewEncoder(w).Encode(config)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": "unsupported attribute " + attr})
	}
}

// InfoHandler responds to API calls
func InfoHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/* APIs
	 * /bond/info
	 */
	config := GetConfigDB()
	templ, err := GetInfoTemplate()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err})
		return
	}

	var colls []ConfigCollection
	for _, v := range config.CollectionsMap {
		colls = append(colls, v)
	}
	sort.Slice(colls, func(i, j int) bool {
		return colls[i].Chunks > colls[j].Chunks
	})
	doc := map[string]interface{}{"Actionlog": config.Actions.Stats, "Collections": colls,
		"Changelog": config.Changes.Stats, "Config": config, "TOP_N": TOP_N}
	if err = templ.Execute(w, doc); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
		return
	}
}
