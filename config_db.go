/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * config_db.go
 */

package bond

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var instance *ConfigDB

// GetLogv2 returns Logv2 instance
func GetConfigDB() *ConfigDB {
	if instance == nil {
		instance = &ConfigDB{}
	}
	return instance
}

type ConfigChunk struct {
	Chunks int              `bson:"chunks"`
	Jumbo  int              `bson:"jumbo"`
	NS     string           `bson:"ns"`
	Shard  string           `bson:"shard"`
	UUID   primitive.Binary `bson:"uuid"`
}

type ConfigCollection struct {
	Chunks    int    `bson:"chunks"`
	Dropped   bool   `bson:"dropped"`
	ID        string `bson:"_id"`
	Key       bson.D `bson:"key"`
	NoBalance bool   `bson:"noBalance"`
	Unique    bool   `bson:"unique"`
	UUID      primitive.Binary
}

type ConfigDatabase struct {
	ID          string `bson:"_id"`
	Partitioned bool   `bson:"partitioned"`
	Primary     string `bson:"primary"`
}

type ConfigMongos struct {
	AdvisoryHostFQDNs *[]string           `bson:"advisoryHostFQDNs"`
	Created           *primitive.DateTime `bson:"created"`
	ID                *string             `bson:"_id"`
	MongoVersion      *string             `bson:"mongoVersion"`
	Ping              *primitive.DateTime `bson:"ping"`
	Up                *int64              `bson:"up"`
	Waiting           *bool               `bson:"waiting"`
}

type ConfigShard struct {
	Chunks  int     `bson:"chunks"`
	Host    *string `bson:"host"`
	Jumbo   int     `bson:"jumbo"`
	ID      *string `bson:"_id"`
	MaxSize *int    `bson:"maxSize"`
	State   *int    `bson:"state"`
}

type ConfigDB struct {
	MongoVersion string `bson:"version"`
	Version      struct {
		ClusterID            *primitive.ObjectID `bson:"clusterId"`
		CurrentVersion       *int                `bson:"currentVersion"`
		ID                   *int                `bson:"_id"`
		MinCompatibleVersion *int                `bson:"minCompatibleVersion"`
	} `bson:"config.version"`
	Mongos         []ConfigMongos              `bson:"config.mongos"`
	ShardsMap      map[string]ConfigShard      `bson:"shards"`
	Databases      []ConfigDatabase            `bson:"config.databases"`
	CollectionsMap map[string]ConfigCollection `bson:"collections"`
	Chunks         []ConfigChunk               `bson:"config.chunks"`

	Actions  *ActionLog          `bson:"actions"`
	Changes  *ChangeLog          `bson:"changes"`
	LastPing *primitive.DateTime `bson:"lastPing"`
	Warnings []string            `bson:"warnings"`

	IsUpgrade     bool
	IsUserVersion bool
	MajorVersion  string

	client       *mongo.Client
	clusterType  string
	serverStatus mdb.ServerStatus
	verbose      bool

	uuid2NS map[string]string
}

func NewConfigDB(uri string, version string, verbose bool) (*ConfigDB, error) {
	cfg := ConfigDB{CollectionsMap: map[string]ConfigCollection{}, MongoVersion: version,
		ShardsMap: map[string]ConfigShard{}, verbose: verbose,
		uuid2NS: map[string]string{}}
	cs, err := mdb.ParseURI(uri)
	if err != nil {
		return nil, err
	}
	cfg.client, err = mdb.NewMongoClient(cs.String())
	if err != nil {
		return nil, err
	}
	s := uri
	if cs.Password != "" {
		s = strings.Replace(s, url.QueryEscape(cs.Password), "xxxxxx", 1)
	}
	cfg.serverStatus, err = mdb.GetServerStatus(cfg.client)
	if err != nil {
		return nil, err
	}
	cfg.clusterType = mdb.GetClusterType(cfg.serverStatus)
	if cfg.clusterType == mdb.Sharded {
		cfg.MongoVersion = cfg.serverStatus.Version
		log.Println("connected to mongos", s)
		log.Println(cfg.clusterType, "cluster, version", cfg.MongoVersion)
	} else if version != "" { // can be from a config mongodump
		cfg.IsUserVersion = true
		log.Println("connected to mongod", s)
		log.Println("given mongo version", version)
	} else {
		message := "mongo version is required, use -mongo <version>"
		log.Println("connected to mongod", s)
		return nil, errors.New(message)
	}
	toks := strings.Split(cfg.MongoVersion, ".")
	if len(toks) > 1 {
		cfg.MajorVersion = strings.Join(toks[:2], ".")
		log.Println("major version", cfg.MajorVersion)
	}
	instance = &cfg
	return instance, nil
}

func (ptr *ConfigDB) GetShardingInfo() error {
	log.Println("GetShardingInfo()")
	ctx := context.Background()
	db := ptr.client.Database("config")
	// check version
	err := db.Collection("version").FindOne(ctx, bson.D{}).Decode(&ptr.Version)
	if err != nil {
		return err
	}

	// get shards
	log.Println("GetShardingInfo() get shards")
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "_id", Value: 1}})
	cursor, err := db.Collection("shards").Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc ConfigShard
		cursor.Decode(&doc)
		ptr.ShardsMap[*doc.ID] = doc
	}
	defer cursor.Close(ctx)

	// get mongos
	log.Println("GetShardingInfo() get mongos")
	opts.SetSort(bson.D{{Key: "ping", Value: -1}})
	cursor, err = db.Collection("mongos").Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc ConfigMongos
		cursor.Decode(&doc)
		ptr.Mongos = append(ptr.Mongos, doc)
		if len(ptr.Mongos) == 1 {
			ptr.LastPing = doc.Ping
		}
	}
	defer cursor.Close(ctx)

	// get databases
	log.Println("GetShardingInfo() get databases")
	opts.SetSort(bson.D{{Key: "partitioned", Value: -1}, {Key: "_id", Value: 1}})
	cursor, err = db.Collection("databases").Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc ConfigDatabase
		cursor.Decode(&doc)
		ptr.Databases = append(ptr.Databases, doc)
	}
	defer cursor.Close(ctx)

	// get collections
	log.Println("GetShardingInfo() get collections")
	opts.SetSort(bson.D{{Key: "_id", Value: 1}})
	filter := bson.D{
		//		{Key: "_id", Value: bson.D{{Key: "$ne", Value: "config.system.sessions"}}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "dropped", Value: false}},
			bson.D{{Key: "dropped", Value: bson.D{{Key: "$exists", Value: false}}}},
		}},
	}
	cursor, err = db.Collection("collections").Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc ConfigCollection
		cursor.Decode(&doc)
		ptr.CollectionsMap[doc.ID] = doc
		ptr.uuid2NS[string(doc.UUID.Data)] = doc.ID
	}
	defer cursor.Close(ctx)

	// check chunks
	if err = ptr.GetChunksInfo(); err != nil {
		return err
	}
	return nil
}

// CheckLogs checks actionlog and changelog
func (ptr *ConfigDB) CheckLogs() error {
	log.Println("CheckLogs()")
	actions, err := NewActionLog(ptr.client)
	if err != nil {
		return err
	}
	if err = actions.GetStats(); err != nil {
		return err
	}
	if err = actions.GetBalancerRounds(); err != nil {
		return err
	}
	ptr.Actions = actions

	changes, err := NewChangeLog(ptr.client)
	if err != nil {
		return err
	}
	if err = changes.GetSplits(); err != nil {
		return err
	}
	ptr.Changes = changes
	return nil
}

func (ptr *ConfigDB) GetChunksInfo() error {
	log.Println("GetChunksInfo()")
	ctx := context.Background()
	db := ptr.client.Database("config")
	var chunk bson.M
	err := db.Collection("chunks").FindOne(ctx, bson.D{{}}).Decode(&chunk)
	if err != nil {
		return err
	}
	pipeline := bson.A{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "shard", Value: "$shard"}, {Key: "uuid", Value: "$uuid"}}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "jumbo", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$cond", Value: bson.D{
							{Key: "if", Value: bson.D{
								{Key: "$eq", Value: bson.A{"$jumbo", true}},
							}},
							{Key: "then", Value: 1},
							{Key: "else", Value: 0},
						}},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "_id.shard", Value: 1}, {Key: "_id.uuid", Value: 1}}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "shard", Value: "$_id.shard"},
				{Key: "uuid", Value: "$_id.uuid"},
				{Key: "chunks", Value: "$total"},
				{Key: "jumbo", Value: "$jumbo"},
				{Key: "_id", Value: 0},
			}},
		},
	}
	if chunk["uuid"] == nil {
		pipeline = bson.A{
			bson.D{
				{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{{Key: "shard", Value: "$shard"}, {Key: "ns", Value: "$ns"}}},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
					{Key: "jumbo", Value: bson.D{
						{Key: "$sum", Value: bson.D{
							{Key: "$cond", Value: bson.D{
								{Key: "if", Value: bson.D{
									{Key: "$eq", Value: bson.A{"$jumbo", true}},
								}},
								{Key: "then", Value: 1},
								{Key: "else", Value: 0},
							}},
						}},
					}},
				}},
			},
			bson.D{
				{Key: "$sort", Value: bson.D{
					{Key: "_id.shard", Value: 1},
					{Key: "_id.ns", Value: 1},
				},
				},
			},
			bson.D{
				{Key: "$project", Value: bson.D{
					{Key: "shard", Value: "$_id.shard"},
					{Key: "ns", Value: "$_id.ns"},
					{Key: "chunks", Value: "$total"},
					{Key: "jumbo", Value: "$jumbo"},
					{Key: "_id", Value: 0},
				},
				},
			},
		}
	}
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := db.Collection("chunks").Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc ConfigChunk
		cursor.Decode(&doc)
		tally := ptr.ShardsMap[doc.Shard]
		tally.Jumbo += doc.Jumbo
		tally.Chunks += doc.Chunks
		ptr.ShardsMap[doc.Shard] = tally

		if chunk["uuid"] != nil {
			doc.NS = ptr.uuid2NS[string(doc.UUID.Data)]
		}
		ptr.Chunks = append(ptr.Chunks, doc)
		cally := ptr.CollectionsMap[doc.NS]
		cally.Chunks += doc.Chunks
		ptr.CollectionsMap[doc.NS] = cally
	}
	defer cursor.Close(ctx)
	return nil
}

func (ptr *ConfigDB) PrintInfo() {
	log.Println("config.version:")
	log.Println(Stringify(ptr.Version))

	log.Println("config.shards:")
	i := 1
	for _, shard := range ptr.ShardsMap {
		i++
		log.Println(i, ":", Stringify(shard))
	}

	log.Println("config.mongos:")
	for i, mongos := range ptr.Mongos {
		log.Println(i+1, ":", Stringify(mongos))
	}
}

// CheckWarnings checks any misconfigrations and anomalies
func (ptr *ConfigDB) CheckWarnings() error {
	/*
	 * if # of collections greater than 1,000
	 * if any mistached major version of mongos instances
	 * if maxSize is configured in shards
	 * if config.actionlog a capped collection if exists
	 * if config.changelog a capped collection
	 * if version in the MongoDB required upgrade list
	 */
	printer := message.NewPrinter(language.English)
	if len(ptr.CollectionsMap) > 10000 {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Hey dude, you have a lot of collections: %d.", len(ptr.CollectionsMap)))
	}
	if len(ptr.Mongos) == 0 {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("No mongos instance was found, probably a restored cluster from a config database dump."))
	} else {
		count := 0
		for _, value := range ptr.Mongos {
			toks := strings.Split(*value.MongoVersion, ".")
			if len(toks) < 2 {
				continue
			}
			if strings.Join(toks[:2], ".") != ptr.MajorVersion {
				count++
			}
		}
		if count > 0 {
			ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Mismatched major version of  mongos: %d.", count))
		}
	}
	count := 0
	for _, shard := range ptr.ShardsMap {
		if shard.MaxSize != nil {
			count++
		}
	}
	if count > 0 {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("%d shards have 'maxSize' configured.", count))
	}

	if ptr.Actions.Stats.Capped == nil {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Collection config.actionlog doesn't exist."))
	} else if !*ptr.Actions.Stats.Capped {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Collection config.actionlog is not a capped collection."))
	}
	if !*ptr.Changes.Stats.Capped {
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Collection config.changelog is not a capped collection."))
	}

	str := CheckUpgradeRecommendation(ptr.MongoVersion)
	if str != "" {
		ptr.IsUpgrade = true
		ptr.Warnings = append(ptr.Warnings, printer.Sprintf("Suggest upgrade to latest MongoDB version, see %s for details.", str))
	}
	return nil
}
