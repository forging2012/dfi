// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <http://unlicense.org/>
package dht

/*
	This file stores all the SQL queries needed for the NetDB.
	It will also be used to prepare all SQL statements :)
*/

const (
	/*
		id             - primary key
		address        - the encoded dfi address
		name           - the user-defined name of the node
		desc           - the user-defined description of the node
		publicAddress  - the user-defined address of the node on the internet/tor/etc
		port           - the port used to connect to the node
		publicKey      - the publicKey for the node, can be checked with the address
		signature      - this entry, signed
		collectionHash - the root hash of all the posts
		postCount      - the number of posts this node has
		seedCount      - the number of seeds this node has
		updated        - when this entry was last updated by the node, or another adding seeds
		seen           - when this node was last seen online

		DFI addresses are stored encoded mostly because it makes debugging *far*
		easier, at the code of some extra encoding and decoding.
	*/
	sqlCreateEntriesTable = `
			CREATE TABLE IF NOT EXISTS
				entry(
					id INTEGER PRIMARY KEY NOT NULL,
					address STRING(40) UNIQUE,
					name STRING(64) NOT NULL,
					desc STRING(256),
					publicAddress STRING(256) NOT NULL,
					port INT,
					publicKey BLOB(32) NOT NULL,
					signature BLOB(64),
					collectionHash BLOB(32),
					postCount INT,
					seedCount INT,
					seedingCount INT,
					updated INT,
					seen INT
				)
	`

	// Create the seeds table, using to link together seeds and the actual node
	// constraint should make sure we don't end up with duplicate seeds
	// TODO: Make sure the constraint is only one way. IE, allow both x,y and y,x
	// to exist.
	sqlCreateSeedsTable = `
		CREATE TABLE IF NOT EXISTS 
				seed(
					id INTEGER PRIMARY KEY NOT NULL,
					seed INTEGER NOT NULL,
					for INTEGER NOT NULL,
					UNIQUE(seed, for) ON CONFLICT REPLACE
				)
	`
	// The full text search virtual table, allowing for the search of a node by
	// description and name.
	sqlCreateFtsTable = `
			CREATE VIRTUAL TABLE IF NOT EXISTS
				ftsEntry using fts4(
					content="entry",
					name,
					desc
				)
	`
	sqlUpdateEntry = `
			UPDATE entry SET 
				name=?,
				desc=?,
				publicAddress=?,
				port=?,
				publicKey=?,
				signature=?,
				collectionHash=?,
				postCount=?,
				seedCount=?,
				seedingCount=?,
				updated=?,
				seen=?
			WHERE address=?
	`

	sqlInsertEntry = `
			INSERT OR IGNORE INTO entry (
				address,
				name,
				desc,
				publicAddress,
				port,
				publicKey,
				signature,
				collectionHash,
				postCount,
				seedCount,
				seedingCount,
				updated,
				seen
			)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	sqlInsertSeed = `
			INSERT OR IGNORE INTO seed (
				seed,
				for
			) VALUES (?, ?)
	`

	sqlInsertFtsEntry = `
			INSERT OR IGNORE INTO ftsEntry (
				docid,
				name,
				desc
			) VALUES(?, ?, ?)
	`

	// We need an index on addresses, as nodes wll be fetched by index really
	// quite often. Most of the time actually! It's probably a good idea to cache
	// in RAM for n seconds.
	sqlIndexAddresses = `
			CREATE INDEX IF NOT EXISTS
				addressIndex ON entry(address)
	`

	sqlQueryAddress = `
		SELECT * FROM entry WHERE address=?
	`

	sqlQueryIdByAddress = `
		SELECT id FROM entry WHERE address=?
	`

	// Get all the seeders for a given address
	sqlQuerySeeds = `
		SELECT entry.address FROM entry
			JOIN seed
				ON entry.id = seed.seed
			WHERE seed.for = ?
	`

	// pretty much the opposite of the above, get a list of addresses that the
	// peer is seeding
	sqlQuerySeeding = `
		SELECT entry.address FROM entry
			JOIN seed
				ON entry.id = seed.for
			WHERE seed.seed = ?
	`

	sqlEntryLen = `
		SELECT MAX(id) FROM entry
	`

	sqlQueryLatest = `
		SELECT * FROM entry ORDER BY id DESC LIMIT 20
	`

	sqlSearchEntries = `
		SELECT entry.address FROM entry 
			WHERE entry.id=(
				SELECT docid FROM ftsEntry WHERE name MATCH ? OR desc MATCH ?
			)
		LIMIT ?,?
	`
)
