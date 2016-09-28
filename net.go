// a few network helpers

package zif

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
)

type ConnHeader struct {
	conn   net.Conn
	header ProtocolHeader
}

func net_recvall(buf []byte, conn net.Conn) error {
	read := 0

	for read < len(buf) {

		if conn == nil {
			return errors.New("Cannot read, connection nil")
		}

		r, err := conn.Read(buf[read:])

		if err != nil {
			return err
		}

		read += r
	}

	return nil
}

// len returns an int
// so why uint64?
// One day this protocol may be implemented in not-Go, and I'd just rather not
// be constrained to signed integers that could be either 32 or 64 bit :)
// this way it is known what is going on
func net_sendlength(conn net.Conn, length uint64) error {
	length_b := make([]byte, 8)
	binary.PutUvarint(length_b, length)

	_, err := conn.Write(length_b)

	return err
}

func net_recvlength(conn net.Conn) (uint64, error) {
	length_b := make([]byte, 8)
	err := net_recvall(length_b, conn)

	if err != nil {
		return 0, err
	}

	length, _ := binary.Uvarint(length_b)

	return length, nil
}

func net_sendpost(conn net.Conn, post Post) error {
	json, err := post.Json()

	if err != nil {
		return err
	}

	net_sendlength(conn, uint64(len(json)))

	ok := make([]byte, 2)
	net_recvall(ok, conn)

	if !bytes.Equal(proto_ok, ok) {
		return errors.New("Peer refused entry")
	}

	conn.Write(json)

	return nil
}

func net_recvpost(conn net.Conn) (*Post, error) {
	length, err := net_recvlength(conn)

	if err != nil {
		return nil, err
	}

	if length > MaxPostSize {
		return nil, errors.New("Post too large")
	}

	conn.Write(proto_ok)

	buf := make([]byte, length)
	net_recvall(buf, conn)

	var post Post
	json.Unmarshal(buf, &post)

	return &post, nil
}

func check_ok(conn net.Conn) bool {
	buf := make([]byte, 2)

	net_recvall(buf, conn)

	return bytes.Equal(buf, proto_ok)
}

func recieve_entry(conn net.Conn) (Entry, error) {
	length_b := make([]byte, 8)
	net_recvall(length_b, conn)
	length, _ := binary.Varint(length_b)

	if length > EntryLengthMax {
		return Entry{}, errors.New("Peer entry larger than max")
	}

	entry_json := make([]byte, length)
	net_recvall(entry_json, conn)

	entry, err := JsonToEntry(entry_json)

	err = ValidateEntry(&entry)

	return entry, err
}
