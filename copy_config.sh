#! /bin/bash
# Copyright 2023-present Kuei-chun Chen. All rights reserved.
# copy_config.sh

source="mongodb://admin:password@localhost/?authSource=admin"
target="mongodb://localhost:27003/"
if [ "$1" != "" ] && [ "$2" != "" ]; then
    source="$1"
    target="$2"
fi

mongodump --uri="$source" --db config --archive=testdata/bond-confiig.gz --gzip
mongodump --uri="$source" --db config --archive=testdata/bond-config-mongos.gz --gzip --collection mongos

mongosh "$target" --eval "db.getSiblingDB('config').dropDatabase();"
mongorestore --uri="$target" --archive=testdata/bond-confiig.gz --gzip
mongorestore --uri="$target" --archive=testdata/bond-config-mongos.gz --gzip

# update config.shards
mongosh "$target" --eval "
    db.getSiblingDB('config').shards.updateOne({'_id': 'shard01'}, 
        {'\$set': {'host': 'shard01/bond-ATL-1.simagix.com:27016,bond-ATL-2.simagix.com:27016,bond-ATL-3.simagix.com:27016'}});
    db.getSiblingDB('config').shards.updateOne({'_id': 'shard02'}, 
        {'\$set': {'host': 'shard02/bond-SEA-1.simagix.com:27016,bond-SEA-2.simagix.com:27016,bond-SEA-3.simagix.com:27016'}});

    var doc = db.getSiblingDB('config').mongos.findOne({'_id': 'Mac-mini-server.local:27017'});
    if (doc != null) {
        doc._id = 'bond-ATL-1.simagix.com:27017';
        doc.advisoryHostFQDNs = ['bond-ATL-1.simagix.com:27017'];
        db.getSiblingDB('config').mongos.insertOne(doc);
        db.getSiblingDB('config').mongos.remove({'_id': 'Mac-mini-server.local:27017'});
    }

    db.getSiblingDB('config').mongos.deleteOne({'_id': 'bond-SEA-1.simagix.com:27017'});
    db.getSiblingDB('config').mongos.insertOne({'_id': 'bond-SEA-1.simagix.com:27017', 
        advisoryHostFQDNs: [ 'bond-SEA-1.simagix.com' ], mongoVersion: '4.4.6', ping: ISODate('2023-08-26T09:16:10.312Z'), waiting: false});
"
