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
package util

import (
	"bufio"
	crand "crypto/rand"
	"io"
	"math/big"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func CryptoRandBytes(size int) ([]byte, error) {
	buf := make([]byte, size)

	if size <= 0 {
		return buf, nil
	}
	_, err := crand.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func CryptoRandInt(min, max int64) int64 {
	if max-min <= 0 {
		return 0 // so random
	}

	num, err := crand.Int(crand.Reader, big.NewInt(max-min))

	if err != nil {
		log.Error(err.Error())
		return min
	}

	return num.Int64() + min
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func ReadPost(r io.Reader, delim byte) {
	br := bufio.NewReader(r)

	br.ReadString(delim)
}

// Returns everything in two that is not in one
func SliceDiff(one [][]byte, two [][]byte) [][]byte {
	// first build a map out of one
	encountered := make(map[string]bool)
	result := make([][]byte, 0, len(one)+len(two))

	for _, i := range one {
		encountered[string(i)] = true
	}

	// then add to result if not in the map
	for _, i := range two {
		if _, ok := encountered[string(i)]; !ok {
			result = append(result, i)
		}
	}

	return result
}

func MergeSeeds(one [][]byte, two [][]byte) [][]byte {
	// make a map
	encountered := make(map[string]bool)
	result := make([][]byte, 0, len(one)+len(two))

	for _, i := range one {
		encountered[string(i)] = true
	}

	for _, i := range two {
		encountered[string(i)] = true
	}

	for k, _ := range encountered {
		result = append(result, []byte(k))
	}

	return result
}

func ShuffleBytes(slice [][]byte) {
	for i := range slice {
		j := rand.Intn(i + 1)

		slice[i], slice[j] = slice[j], slice[i]
	}
}
