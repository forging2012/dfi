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
package dht_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dfindex/dfi/dht"
	"github.com/dfindex/dfi/util"
	"golang.org/x/crypto/ed25519"
)

// this is helpful for testing
// thanks to: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func fatalErr(err error, t testing.TB) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func makeTesting() {
	os.Mkdir(".testing", 0777)
}

func removeTesting() {
	os.RemoveAll(".testing")
}

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func randomAddress(t testing.TB) *dht.Address {
	addr, err := dht.RandomAddress()

	if err != nil {
		t.Fatal(err.Error())
	}

	return addr
}

func dbWithRandomAddress(t testing.TB) *dht.NetDB {
	// pretty much just tests that the SQL gets prepared properly
	addr := randomAddress(t)

	db, err := dht.NewNetDB(*addr, ".testing/"+addr.StringOr(""))

	if err != nil {
		t.Fatal(err.Error())
	}

	return db
}

func randomEntry(t testing.TB) dht.Entry {
	name := randString(util.RandInt(5, 25))
	desc := randString(util.RandInt(5, 144))

	pub, priv, err := ed25519.GenerateKey(nil)
	addr := dht.Address{}
	addr.Generate(pub)

	entry := dht.Entry{
		Name:          name,
		Desc:          desc,
		Address:       addr,
		PublicKey:     pub,
		PublicAddress: "localhost",
		Port:          5050,
	}

	dat, err := entry.Bytes()

	if err != nil {
		t.Fatal(err)
	}

	sig := ed25519.Sign(priv, dat)

	entry.Signature = sig

	return entry
}

func TestMain(m *testing.M) {

	makeTesting()
	ret := m.Run()
	removeTesting()
	os.Exit(ret)
}

func TestNewNetDB(t *testing.T) {
	dbWithRandomAddress(t)
}

// Tests Insert, and by extension len and tablelen
func TestNetDBInsertAndLen(t *testing.T) {
	db := dbWithRandomAddress(t)

	entry := randomEntry(t)

	_, err := db.Insert(entry)

	if err != nil {
		t.Fatal(err)
	}

	length, err := db.Len()
	if err != nil {
		t.Fatal(err.Error())
	} else if length != 1 {
		t.Fatal("Database insert failed")
	}

	length = db.TableLen()

	if length != 1 {
		t.Fatal("Database insert failed")
	}
}

func TestInsert(t *testing.T) {
	db := dbWithRandomAddress(t)

	entry := randomEntry(t)

	_, err := db.Insert(entry)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if l, _ := db.Len(); l != 1 {
		t.Fatalf("Insert failed, database len: %d", l)
	}
}

func TestInsertSeed(t *testing.T) {
	db := dbWithRandomAddress(t)
	entry := randomEntry(t)
	seed := randomEntry(t)

	// insert the entries first
	_, err := db.Insert(entry)
	fatalErr(err, t)

	_, err = db.Insert(seed)
	fatalErr(err, t)

	// then register some seeds :)
	fatalErr(db.InsertSeed(entry.Address, seed.Address), t)
	t.Log("Inserted seeds")

	seeds, err := db.QuerySeeds(entry.Address)
	fatalErr(err, t)

	seeding, err := db.QuerySeeding(seed.Address)
	fatalErr(err, t)

	if len(seeds) != 1 {
		t.Fatalf("Seeds not registered properly, length: %d", len(seeds))
	}

	if len(seeding) != 1 {
		t.Fatal("Seeding not registered properly")
	}

	if !seeds[0].Equals(&seed.Address) {
		t.Fatal("Seed address not correct")
	}

	if !seeding[0].Equals(&entry.Address) {
		t.Fatal("Seeding address not correct")
	}
}

func BenchmarkInsert(b *testing.B) {
	makeTesting()
	db := dbWithRandomAddress(b)

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		entry := randomEntry(b)
		b.StartTimer()

		db.Insert(entry)
	}

	removeTesting()
}

func BenchmarkQuery(b *testing.B) {
	makeTesting()
	db := dbWithRandomAddress(b)

	for n := 0; n < b.N; n++ {
		// don't time the actual inserting!
		// this is not what we are measuring in this case
		b.StopTimer()
		entry := randomEntry(b)
		db.Insert(entry)
		b.StartTimer()

		db.Query(entry.Address)
	}

	removeTesting()
}
