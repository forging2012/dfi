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
	"errors"
	"io"
	"net"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	"github.com/zif/zif/common"
	"github.com/zif/zif/dht"
)

type Message struct {
	Header      string
	Stream      net.Conn
	Client      *Client
	From        *dht.Address
	Compression string

	Content []byte
}

func (m *Message) Write(iface interface{}) error {
	writer := bytes.Buffer{}
	encoder := msgpack.NewEncoder(&writer)

	err := encoder.Encode(iface)

	if err != nil {
		return err
	}

	m.Content = writer.Bytes()

	return nil
}

func (m *Message) Read(iface interface{}) error {
	if m.Content == nil {
		return errors.New("Message has no content")
	}

	reader := bytes.NewReader(m.Content)
	limiter := &io.LimitedReader{R: reader, N: common.MaxMessageContentSize}

	decoder := msgpack.NewDecoder(limiter)

	err := decoder.Decode(iface)

	return err
}

func (m *Message) ReadInt() (int, error) {
	var ret int

	err := m.Read(&ret)

	return ret, err
}

func (m *Message) Json() ([]byte, error) {
	return msgpack.Marshal(m)
}

// Ok() is just an easier way to check if the peer has sent an "ok" response,
// rather than comparing the header member to a constant repeatedly.
func (m *Message) Ok() bool {
	return m.Header == ProtoOk
}
