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
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	"golang.org/x/crypto/ed25519"
)

const (
	MaxEntryNameLength          = 32
	MaxEntryDescLength          = 160
	MaxEntryPublicAddressLength = 253
	MaxEntrySeeds               = 100000
)

// This is an entry into the DHT. It is used to connect to a peer given just
// it's DFI address.
type Entry struct {
	Address       Address `json:"address"`
	Name          string  `json:"name"`
	Desc          string  `json:"desc"`
	PublicAddress string  `json:"publicAddress"`
	PublicKey     []byte  `json:"publicKey"`
	PostCount     int     `json:"postCount"`
	Updated       uint64  `json:"updated"`

	// The owner of this entry should have signed it, we need to store the
	// sigature. It's actually okay as we can verify that a peer owns a public
	// key by generating an address from it - if the address is not the peers,
	// then Mallory is just using someone elses entry for their own address.
	Signature      []byte `json:"signature"`
	CollectionHash []byte `json:"collectionHash"`
	Port           int    `json:"port"`

	Seeds   [][]byte `json:"seeds"`
	Seeding [][]byte `json:"seeding"`
	Seen    int      `json:"seed"`

	// Used in the FindClosest function, for sorting.
	distance Address
}

// true if JSON, false if msgpack
func DecodeEntry(data []byte, isJson bool) (*Entry, error) {
	var err error
	e := &Entry{}

	if isJson {
		err = json.Unmarshal(data, e)
	} else {
		err = msgpack.Unmarshal(data, e)
	}

	if err != nil {
		return nil, err
	}

	return e, nil
}

// This is signed, *not* the JSON. This is needed because otherwise the order of
// the posts encoded is not actually guaranteed, which can lead to invalid
// signatures. Plus we can only sign data that is actually needed.
func (e Entry) Bytes() ([]byte, error) {
	ret, err := e.String()
	return []byte(ret), err
}

func (e Entry) String() (string, error) {
	var str string

	addressString, err := e.Address.String()

	if err != nil {
		return "", err
	}

	postCount := strconv.Itoa(e.PostCount)
	updated := strconv.Itoa(int(e.Updated))

	str += addressString
	str += e.Name
	str += e.Desc
	str += string(e.PublicAddress)
	str += string(e.PublicKey)
	str += string(e.Port)
	str += postCount
	str += updated
	str += string(e.CollectionHash)

	for _, i := range e.Seeding {
		str += string(i)
	}

	// note that we do not, in fact, sign who the seeds are. This allows others
	// to build the swarm while this peer is not online.

	return str, nil
}

func (e Entry) Encode() ([]byte, error) {
	return msgpack.Marshal(e)
}

// Returns a JSON encoded string, not msgpack. This is because it is likely
// going to be seen by a human, otherwise it would be bytes.
func (e Entry) EncodeString() (string, error) {
	enc, err := json.Marshal(e)

	if err != nil {
		return "", err
	}

	return string(enc), err
}

func (e *Entry) SetLocalPeer(lp Node) {
	e.Address = *lp.Address()

	e.PublicKey = make([]byte, len(lp.PublicKey()))
	copy(e.PublicKey, lp.PublicKey())
	e.PublicKey = lp.PublicKey()
}

type Entries []*Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Entries) Less(i, j int) bool {
	return e[i].distance.Less(&e[j].distance)
}

// Ensures that all the members of an entry struct fit the requirements for the
// DFI libdficol. If an entry passes this, then we should be able to perform
// most operations on it.
func (entry *Entry) Verify() error {
	if entry == nil {
		return errors.New("Entry is nil")
	}

	if len(entry.Address.Raw) != 20 {
		return errors.New("Address size invalid")
	}

	if len(entry.Name) > MaxEntryNameLength {
		return errors.New("Entry name is too long")
	}

	if len(entry.Desc) > MaxEntryDescLength {
		return errors.New("Entry description is too long")
	}

	if len(entry.Seeds) > MaxEntrySeeds {
		return errors.New("Entry has too many seeds")
	}

	if len(entry.PublicKey) < ed25519.PublicKeySize {
		return errors.New(fmt.Sprintf("Public key too small: %d", len(entry.PublicKey)))
	}

	if len(entry.Signature) < ed25519.SignatureSize {
		return errors.New("Signature too small")
	}

	data, _ := entry.Bytes()
	verified := ed25519.Verify(entry.PublicKey, data, entry.Signature[:])

	if !verified {
		return errors.New("Failed to verify signature")
	}

	if len(entry.PublicAddress) == 0 {
		return errors.New("Public address must be set")
	}

	// 253 is the maximum length of a domain name
	if len(entry.PublicAddress) >= MaxEntryPublicAddressLength {
		return errors.New("Public address is too large (253 char max)")
	}

	if entry.Port > 65535 {
		return errors.New("Port too large (" + string(entry.Port) + ")")
	}

	return nil
}

func ShuffleEntries(slice Entries) {
	for i := range slice {
		j := rand.Intn(i + 1)

		slice[i], slice[j] = slice[j], slice[i]
	}
}
