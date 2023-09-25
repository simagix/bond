/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * handlers.go
 */

package bond

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
)

const (
	T_MIGRATION_TIME  = "migration_time"
	T_MIGRATION_STATS = "migration_stats"
	T_CHUNK_SPLITS    = "splits"
)

type Chart struct {
	Index int
	Title string
	Descr string
	URL   string
}

var charts = map[string]Chart{
	"instruction": {0, "select a chart", "", ""},
	T_MIGRATION_TIME: {1, "Average Chunk Migration Time",
		"Display average migration time", "/bond/charts/" + T_MIGRATION_TIME},
	T_MIGRATION_STATS: {2, "Chunk Migration Counts",
		"Display average migration time", "/bond/charts/" + T_MIGRATION_STATS},
	T_CHUNK_SPLITS: {3, "No. of Chunk Splits",
		"Display chunk splits", "/bond/charts/" + T_CHUNK_SPLITS},
}

type NameValue struct {
	Name  string `bson:"name"`
	Value int    `bson:"value"`
}

// ChartsHandler renders charts for UI selections
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

// ShardChartHandler renders pie charts for percentage of collections distribution for a given shard
func ShardChartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/* APIs
	 * /bond/chart/shards/:shard
	 */
	config := GetConfigDB()
	shard := params.ByName("shard")
	templ, err := GetPieChartTemplate()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err})
		return
	}

	var others int
	var docs []NameValue
	var namespaces []NameValue
	for key, value := range config.ShardsMap[shard].namespaces {
		namespaces = append(namespaces, NameValue{key, value})
	}
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Value > namespaces[j].Value
	})
	for _, value := range namespaces {
		if len(docs) > 10 {
			others += value.Value
			continue
		}
		nv := NameValue{value.Name, value.Value}
		docs = append(docs, nv)
	}
	if others > 0 {
		nv := NameValue{"'Beyond the top 10'", others}
		docs = append(docs, nv)
	}
	title := fmt.Sprintf("Collection Chunk Distribution within a Shard\n%s", shard)
	doc := map[string]interface{}{"NameValues": docs, "Title": title}
	if err = templ.Execute(w, doc); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
		return
	}
}

// NamespaceChartHandler renders pie charts for percentage of shards distribution for a given collection
func NamespaceChartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/* APIs
	 * /bond/chart/namespaces/:ns
	 */
	config := GetConfigDB()
	ns := params.ByName("ns")
	templ, err := GetPieChartTemplate()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err})
		return
	}

	var others int
	var docs []NameValue
	var shards []NameValue
	for key, value := range config.CollectionsMap[ns].shards {
		shards = append(shards, NameValue{key, value})
	}
	sort.Slice(shards, func(i, j int) bool {
		return shards[i].Value > shards[j].Value
	})
	for _, value := range shards {
		if len(docs) > 10 {
			others += value.Value
			continue
		}
		nv := NameValue{value.Name, value.Value}
		docs = append(docs, nv)
	}
	if others > 0 {
		nv := NameValue{"'Beyond the top 10'", others}
		docs = append(docs, nv)
	}
	title := fmt.Sprintf("Collection Chunk Distribution among Shards\n%s", ns)
	doc := map[string]interface{}{"NameValues": docs, "Title": title}
	if err = templ.Execute(w, doc); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
		return
	}
}
