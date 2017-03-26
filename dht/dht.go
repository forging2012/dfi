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

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type DHT struct {
	db *NetDB
}

// sets up the dht
func NewDHT(addr Address, path string) *DHT {
	ret := &DHT{}

	db, err := NewNetDB(addr, path)

	if err != nil {
		panic(err)
	}

	ret.db = db

	log.Debug("Loading latest into DHT")
	// insert a load of new entries, keep it fresh!
	entries, err := db.QueryLatest()

	if err == sql.ErrNoRows {
		return ret
	}

	count := 0
	for _, i := range entries {
		count += 1
		db.Insert(i)
	}

	log.WithField("count", count).Debug("Inserted")

	return ret
}

func (dht *DHT) Address() Address {
	return dht.db.addr
}

func (dht *DHT) Insert(entry Entry) (int64, error) {
	// TODO: Announces
	return dht.db.Insert(entry)
}

func (dht *DHT) Query(addr Address) (*Entry, error) {
	entry, _, err := dht.db.Query(addr)

	return entry, err
}

func (dht *DHT) FindClosest(addr Address) (Entries, error) {
	return dht.db.FindClosest(addr)
}

func (dht *DHT) SaveTable(path string) {
	dht.db.SaveTable(path)
}

func (dht *DHT) LoadTable(path string) {
	dht.db.LoadTable(path)
}

func (dht *DHT) SearchEntries(name, desc string, page int) ([]Address, error) {
	return dht.db.SearchPeer(name, desc, page)
}
