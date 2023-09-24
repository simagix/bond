/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * actionlog.go
 */

package bond

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalancerRound struct {
	AverageExecutionTime float64            `bson:"averageExecutionTime"`
	Time                 primitive.DateTime `bson:"time"`
	TotalChunksMoved     int                `bson:"totalChunksMoved"`
	TotalErrors          int                `bson:"totalErrors"`
}

type ActionLog struct {
	BalancerRounds []BalancerRound
	Stats          struct {
		AverageExecutionTime float64 `bson:"averageExecutionTime"`
		Capped               *bool   `bson:"capped"`
		MaxExecutionTime     int64   `bson:"maxExecutionTime"`
		MaxSize              int64   `bson:"maxSize"`
		TotalChunksMoved     int     `bson:"totalChunksMoved"`
		TotalErrors          int     `bson:"totalErrors"`
	}

	client *mongo.Client
}

func NewActionLog(client *mongo.Client) (*ActionLog, error) {
	actions := ActionLog{client: client}
	var doc bson.M
	err := client.Database("config").RunCommand(context.Background(), bson.D{{Key: "collStats", Value: "actionlog"}}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	if doc["capped"] != nil {
		capped := doc["capped"].(bool)
		actions.Stats.Capped = &capped
	}
	if doc["maxSize"] != nil {
		actions.Stats.MaxSize = ToInt64(doc["maxSize"])
	}
	return &actions, nil
}

func (ptr *ActionLog) GetStats() error {
	log.Println("ActionLog.GetStats()")
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "what", Value: "balancer.round"}}}},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: primitive.Null{}},
				{Key: "totalChunksMoved", Value: bson.D{{Key: "$sum", Value: "$details.chunksMoved"}}},
				{Key: "totalErrors", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$details.errorOccured", true}}}, 1, 0},
						},
					}},
				}},
				{Key: "averageExecutionTime", Value: bson.D{{Key: "$avg", Value: "$details.executionTimeMillis"}}},
				{Key: "maxExecutionTime", Value: bson.D{{Key: "$max", Value: "$details.executionTimeMillis"}}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "totalChunksMoved", Value: 1},
				{Key: "totalErrors", Value: 1},
				{Key: "averageExecutionTime", Value: 1},
				{Key: "maxExecutionTime", Value: 1},
			}},
		},
	}
	ctx := context.Background()
	db := ptr.client.Database("config")
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := db.Collection("actionlog").Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	if cursor.Next(ctx) {
		cursor.Decode(&ptr.Stats)
	}
	log.Println(Stringify(ptr.Stats))
	defer cursor.Close(ctx)
	return nil
}

func (ptr *ActionLog) GetBalancerRounds() error {
	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "what", Value: "balancer.round"},
				{Key: "details.executionTimeMillis", Value: bson.D{{Key: "$exists", Value: true}}},
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
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "hour", Value: 1},
				{Key: "executionTimeMillis", Value: "$details.executionTimeMillis"},
				{Key: "chunksMoved", Value: "$details.chunksMoved"},
				{Key: "errorOccured", Value: "$details.errorOccured"},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "hour", Value: "$hour"},
				}},
				{Key: "averageExecutionTime", Value: bson.D{{Key: "$avg", Value: "$executionTimeMillis"}}},
				{Key: "totalChunksMoved", Value: bson.D{{Key: "$sum", Value: "$chunksMoved"}}},
				{Key: "totalErrors", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$errorOccured", true}}}, 1, 0},
						}}},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "time", Value: bson.D{{Key: "$toDate", Value: "$_id.hour"}}},
				{Key: "averageExecutionTime", Value: 1},
				{Key: "totalChunksMoved", Value: 1},
				{Key: "totalErrors", Value: 1},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "time", Value: 1},
			}},
		},
	}
	ctx := context.Background()
	db := ptr.client.Database("config")
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := db.Collection("actionlog").Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var doc BalancerRound
		cursor.Decode(&doc)
		ptr.BalancerRounds = append(ptr.BalancerRounds, doc)
	}
	defer cursor.Close(ctx)
	return nil
}
