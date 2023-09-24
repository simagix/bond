// Copyright 2023-present Kuei-chun Chen. All rights reserved.
// create_chunks.js
function create(dealer, coll) {
	sh.shardCollection(dealer+'.'+coll, {brand: 1, year: 1});
	sh.splitAt(dealer+'.'+coll, {brand: 'M', year: 2013});
	sh.moveChunk(dealer+'.'+coll, { brand: 'M', year: 2013 }, 'shard01');

	for (let i = 'L'.charCodeAt(0); i <= 'S'.charCodeAt(0); i++) {
	    for (let y = 2000; y <= 2024; y+=4) {
			const b = String.fromCharCode(i);
			sh.splitAt(dealer+'.'+coll, {brand: b, year: y});
	    }
	}
}

sh.stopBalancer()

for (let i = 108; i <= 120; i++) {
	let dealer = 'dealer'+i;
	db.getSiblingDB(dealer).dropDatabase();
	sh.enableSharding(dealer);
	for (let j = 6000; j <= 6100; j++) {
		let coll = 'branch-'+j;
		create(dealer, coll);
	}
}

sh.startBalancer()
