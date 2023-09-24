#! /bin/bash
# Copyright 2023-present Kuei-chun Chen. All rights reserved.
# update-config.sh

source="mongodb://admin:password@localhost/?authSource=admin"
target="mongodb://localhost:27003/"
mongodump --uri="$source" --db config --archive=testdata/bond-confiig.gz --gzip
mongodump --uri="$source" --db config --archive=testdata/bond-config-mongos.gz --gzip --collection mongo

mongosh "$target" --eval "db.getSiblingDB('config').dropDatabase();"
mongorestore --uri="$target" --archive=testdata/bond-confiig.gz --gzip
mongorestore --uri="$target" --archive=testdata/bond-config-mongos.gz --gzip

# update config.shards


mongosh "$target" --eval "
    db.getSiblingDB('config').shards.updateOne({'_id': 'shard01'}, 
        {'\$set': {'host': 'shard01/bond-atlanta-1.simagix.com:27016,bond-atlanta-2.simagix.com:27016,bond-atlanta-3.simagix.com:27016'}});
    db.getSiblingDB('config').shards.updateOne({'_id': 'shard02'}, 
        {'\$set': {'host': 'shard02/bond-chicago-1.simagix.com:27016,bond-chicago-2.simagix.com:27016,bond-chicago-3.simagix.com:27016'}});

    db.getSiblingDB('config').shards.updateOne({'_id': 'shard01'}, 
        {'\$set': {'host': 'shard01/bond-atlanta-1.simagix.com:27016,bond-atlanta-2.simagix.com:27016,bond-atlanta-3.simagix.com:27016'}});
    db.getSiblingDB('config').shards.updateOne({'_id': 'shard02'}, 
        {'\$set': {'host': 'shard02/bond-chicago-1.simagix.com:27016,bond-chicago-2.simagix.com:27016,bond-chicago-3.simagix.com:27016'}});

    var doc = db.getSiblingDB('config').mongos.findOne({'_id': 'Mac-mini-server.local:27017'});
    if (doc != null) {
        doc._id = 'bond-atlanta-1.simagix.com:27017';
        doc.advisoryHostFQDNs = ['bond-atlanta-1.simagix.com:27017'];
        db.getSiblingDB('config').mongos.insertOne(doc, {'upsert': true});
        db.getSiblingDB('config').mongos.remove({'_id': 'Mac-mini-server.local:27017'});
    }

    db.getSiblingDB('config').mongos.deleteOne({'_id': 'bond-chicago-1.simagix.com:27017'});
    db.getSiblingDB('config').mongos.insertOne({'_id': 'bond-chicago-1.simagix.com:27017', 
        advisoryHostFQDNs: [ 'bond-chicago-1.simagix.com' ], mongoVersion: '4.4.6', ping: ISODate('2023-08-26T09:16:10.312Z'), waiting: false});
"
