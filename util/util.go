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
