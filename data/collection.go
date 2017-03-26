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

// For more information, please refer to <http://unlicense.org/>
package data

import (
	"errors"
	"hash"
	"io/ioutil"
	"math"

	"golang.org/x/crypto/sha3"
)

// A collection of pieces, by extension a structure containing all posts this
// peer has. Whether or not the pieces are *actually* there is optional, if not
// this is essentially a hash list.
type Collection struct {
	Pieces   []*Piece
	HashList []byte
	RootHash hash.Hash
}

// Create a new collection, set all it's members to the correct default values.
func NewCollection() *Collection {
	col := &Collection{}

	col.RootHash = sha3.New256()
	col.Pieces = make([]*Piece, 0, 2)
	col.HashList = make([]byte, 0)

	return col
}

// Takes a database, starting id, and piece size. From this we create a
// collection, except it does not contain any posts - consider making this optional.
func CreateCollection(db *Database, start, pieceSize int) (*Collection, error) {
	col := NewCollection()

	postCount := db.PostCount()
	pieceCount := int(math.Ceil(float64(postCount) / float64(pieceSize)))

	for i := 0; i < pieceCount; i++ {
		piece, err := db.QueryPiece(uint(i), false)

		if err != nil {
			return nil, err
		}

		col.Add(piece)
	}

	return col, nil
}

// Loads a collection from file.
// This essentially loads the hash list, the data of pieces themselves is just
// left. It's all in the database if it is really needed.
func LoadCollection(path string) (col *Collection, err error) {
	col = NewCollection()

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	if len(data)%32 != 0 {
		err = errors.New("Invalid collection data file")
		return
	}

	col.HashList = data
	col.Rehash()

	return
}

// Save the collection hash list to the given path, with permissions 0777.
func (c *Collection) Save(path string) error {
	return ioutil.WriteFile(path, c.HashList, 0777)
}

// Add a piece to the collection, storing it in c.Pieces and appending it's hash
// to the hash list.
func (c *Collection) Add(piece *Piece) {
	if uint(len(c.HashList)) < piece.Id+1 {
		c.HashList = append(c.HashList, piece.Hash()...)
	} else {
		copy(c.HashList[piece.Id*32:piece.Id*32+32], piece.Hash())
	}

	c.RootHash.Write(piece.Hash())
}

// Return the hash of the hash list, which can then go on to be signed by the
// LocalPeer. This allows proper validation of an entire collection, but the
// localpeer only needs to sign a single hash.
func (c *Collection) Hash() []byte {
	var ret []byte

	ret = c.RootHash.Sum(nil)

	return ret
}

// Regenerates the root hash from the hash list we have.
func (c *Collection) Rehash() {
	c.RootHash = sha3.New256()

	for i := 0; i < len(c.HashList)/32; i++ {
		c.RootHash.Write(c.HashList[i*32 : i*32+32])
	}
}
