/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * changelog.go
 */

package bond

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Split struct {
	Time  primitive.DateTime `bson:"time"`
	Total int                `bson:"total"`
}

type ChangeLog struct {
	Splits []Split `bson:"splits"`

	Stats struct {
		Capped      *bool `bson:"capped"`
		MaxSize     int64 `bson:"maxSize"`
		TotalSplits int   `bson:"total"`
	} `bson:"stats"`
	client *mongo.Client
}

func NewChangeLog(client *mongo.Client) (*ChangeLog, error) {
	changes := ChangeLog{client: client}
	var doc bson.M
	client.Database("config").RunCommand(context.Background(), bson.D{{Key: "collStats", Value: "changelog"}}).Decode(&doc)
	if doc["capped"] != nil {
		capped := doc["capped"].(bool)
		changes.Stats.Capped = &capped
	}
	if doc["capped"] != nil && doc["maxSize"] != nil {
		changes.Stats.MaxSize = ToInt64(doc["maxSize"])
	}
	return &changes, nil
}

func (ptr *ChangeLog) GetSplits() error {
	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "what", Value: bson.D{
					{Key: "$in", Value: bson.A{"split", "multi-split"}},
				}},
				{Key: "time", Value: bson.D{
					{Key: "$exists", Value: true},
					{Key: "$ne", Value: primitive.Null{}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "hour", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m-%dT%H:00:00.000Z"},
						{Key: "date", Value: "$time"},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "hour", Value: "$hour"},
				}},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "time", Value: bson.D{{Key: "$toDate", Value: "$_id.hour"}}},
				{Key: "total", Value: 1},
			}},
		},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "time", Value: 1}}}},
	}
	ctx := context.Background()
	db := ptr.client.Database("config")
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := db.Collection("changelog").Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc Split
		cursor.Decode(&doc)
		ptr.Splits = append(ptr.Splits, doc)
		ptr.Stats.TotalSplits += doc.Total
	}
	defer cursor.Close(ctx)
	return nil
}
