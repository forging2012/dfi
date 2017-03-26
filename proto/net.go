// a few network helpers
package proto

import "github.com/dfindex/dfi/dht"

type ConnHeader struct {
	Client       Client
	Entry        dht.Entry
	Capabilities MessageCapabilities
}
