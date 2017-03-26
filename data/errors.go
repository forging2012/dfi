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
	"bufio"
	"io"

	log "github.com/sirupsen/logrus"
)

type ErrorReader struct {
	reader *bufio.Reader
	Err    error
}

func NewErrorReader(r io.Reader) *ErrorReader {
	return &ErrorReader{bufio.NewReader(r), nil}
}

func (er *ErrorReader) ReadString(delim byte) string {
	var ret string

	ret, er.Err = er.reader.ReadString(delim)

	if er.Err != nil {
		log.Error(er.Err.Error())
		return ""
	}

	return ret[0 : len(ret)-1]
}

func (er *ErrorReader) ReadByte() (byte, error) {
	var ret byte

	ret, er.Err = er.reader.ReadByte()

	if er.Err != nil {
		return 0, er.Err
	}

	return ret, nil
}

type AddressResolutionError struct {
	Address string
}

func (a AddressResolutionError) Error() string {
	return "Failed to resolve address, address may not exist or is not reachable"
}
