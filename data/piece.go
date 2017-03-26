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

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

const PieceSize = 1000

type Piece struct {
	Id    uint
	Posts []Post
	hash  hash.Hash
}

func (p *Piece) Setup() {
	p.hash = sha3.New256()
}

func (p *Piece) Add(post Post, store bool) error {
	if len(p.Posts) > PieceSize {
		return errors.New("Piece full")
	}

	if store {
		p.Posts = append(p.Posts, post)
	}

	data := post.String("|", "", false)
	p.hash.Write([]byte(data))

	return nil
}

func (p *Piece) Hash() []byte {
	var ret []byte

	ret = p.hash.Sum(nil)

	return ret
}

func (p *Piece) Rehash() ([]byte, error) {
	p.hash = sha3.New256()

	for _, i := range p.Posts {
		data := i.Bytes([]byte("|"), []byte(""), false)
		p.hash.Write([]byte(data))
	}

	log.Info("Piece rehashed")

	return p.hash.Sum(nil), nil
}
