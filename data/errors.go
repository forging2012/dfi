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
