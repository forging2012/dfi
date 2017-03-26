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
package proto

import (
	"net"

	"github.com/dfindex/dfi/common"
	"github.com/dfindex/dfi/dht"
	"github.com/hashicorp/yamux"
)

type ProtocolHandler interface {
	common.Signer
	NetworkPeer

	HandleAnnounce(*Message) error
	HandleQuery(*Message) error
	HandleFindClosest(*Message) error
	HandleSearch(*Message) error
	HandleRecent(*Message) error
	HandlePopular(*Message) error
	HandleHashList(*Message) error
	HandlePiece(*Message) error
	HandleAddPeer(*Message) error

	HandleHandshake(ConnHeader) (NetworkPeer, error)
	HandleCloseConnection(*dht.Address)

	GetNetworkPeer(dht.Address) NetworkPeer
	SetNetworkPeer(NetworkPeer)
	GetCapabilities() *MessageCapabilities
}

// Allows the protocol stuff to work with Peers, while libdfi/peer can interface
// peers with the DHT properly.
type NetworkPeer interface {
	Session() *yamux.Session
	AddStream(net.Conn)

	Address() *dht.Address
	Query(dht.Address) (common.Verifier, error)
	FindClosest(dht.Address) ([]common.Verifier, error)
	SetCapabilities(MessageCapabilities)
	UpdateSeen()
}
