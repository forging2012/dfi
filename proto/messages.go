/*
 *  Zif
 *  Copyright (C) 2017 Zif LTD
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published
 *  by the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.

 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package proto

import (
	"bytes"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

// This contains the more "complex" structures that will be sent in message
// content fields.

type MessageCollection struct {
	Hash      []byte
	HashList  []byte
	Size      int
	Signature []byte
}

type MessageSearchQuery struct {
	Query string
	Page  int
}

type MessageRequestPiece struct {
	Address string
	Id      int
	Length  int
}

// Allows us to decode a pieces without also decoding all of the posts within it.
type MessagePiece struct {
	Posts interface{}
}

type MessageCapabilities struct {
	// an array of strings, each a compression type, in order of preference.
	// Index 0 is the preferred method. The method used is the shared method
	// with the lowest index.
	Compression []string
}

func (mp *MessagePiece) Hash() ([]byte, error) {
	hash := sha3.New256()

	//for _, i := range mp.Posts {
	//h := sha3.Sum256([]byte(i))
	//hash.Write(h[:])
	//}

	log.Info("Piece hashed")

	return hash.Sum(nil), nil
}

func (mhl *MessageCollection) Verify(root []byte) error {
	hash := sha3.New256()

	for i := 0; i < mhl.Size; i++ {
		hash.Write(mhl.HashList[32*i : (32*i)+32])
	}

	if !bytes.Equal(hash.Sum(nil), root) {
		return errors.New("Invalid hash list")
	}

	return nil
}

func (mhl *MessageCollection) Encode() ([]byte, error) {
	data, err := json.Marshal(mhl)
	return data, err
}

func (sq *MessageSearchQuery) Encode() ([]byte, error) {
	data, err := json.Marshal(sq)
	return data, err
}

func (mrp *MessageRequestPiece) Encode() ([]byte, error) {
	data, err := json.Marshal(mrp)
	return data, err
}
