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

package dfi

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/dfindex/dfi/data"
)

// Command input types

type CommandPeer struct {
	Address string `json:"address"`
}

type CommandFile struct {
	File string `json:"file"`
}

type CommandPing CommandPeer
type CommandAnnounce CommandPeer

type CommandRequestAddPeer struct {
	// The peer to send the request to
	Remote string `json:"remote"`
	// The peer we wish to be registered as a peer for
	Peer string `json:"peer"`
}

type CommandSearchEntry struct {
	Name string
	Desc string
	Page int
}

type CommandRSearch struct {
	CommandPeer
	Query string `json:"query"`
	Page  int    `json:"page"`
}
type CommandPeerSearch CommandRSearch
type CommandPeerRecent struct {
	CommandPeer
	Page int `json:"page"`
}
type CommandPeerPopular CommandPeerRecent
type CommandMirror CommandPeer
type CommandMirrorProgress CommandPeer
type CommandPeerIndex struct {
	CommandPeer
	Since int `json:"since"`
}

type CommandMeta struct {
	PId int `json:"pid"`
}

type CommandAddPost struct {
	data.Post
	Index bool
}
type CommandSelfIndex struct {
	Since int `json:"since"`
}
type CommandResolve CommandPeer
type CommandBootstrap CommandPeer

type CommandSuggest struct {
	Query string `json:"query"`
}

type CommandSelfSearch struct {
	CommandSuggest
	Page int `json:"page"`
}
type CommandSelfRecent struct {
	Page int `json:"page"`
}
type CommandSelfPopular CommandSelfRecent
type CommandAddMeta struct {
	CommandMeta
	Value string `json:"value"`
}
type CommandGetMeta CommandMeta
type CommandSaveCollection interface{}
type CommandRebuildCollection interface{}
type CommandPeers interface{}
type CommandSaveRoutingTable interface{}

// Used for setting values in the localpeer entry
type CommandLocalSet struct {
	Key   string `json:"key"`
	Value string `json:"key"`
}

type CommandLocalGet struct {
	Key string `json:"key"`
}

type CommandAddressEncode struct {
	Raw []byte `json:"raw"`
}

type CommandSetSeedLeech struct {
	Id       uint
	Seeders  uint
	Leechers uint
}

type CommandNetMap struct {
	Address string
}

// Command output types

type CommandResult struct {
	IsOK   bool        `json:"status"`
	Result interface{} `json:"value"`
	Error  error       `json:"err"`
}

func (cr *CommandResult) WriteJSON(w io.Writer) {
	e := json.NewEncoder(w)

	if cr.IsOK {
		if cr.Result == nil {
			e.Encode(struct {
				Status string `json:"status"`
			}{"ok"})
		} else {
			e.Encode(struct {
				Status string      `json:"status"`
				Value  interface{} `json:"value"`
			}{"ok", cr.Result})
		}
	} else {
		if cr.Error == nil {
			cr.Error = errors.New("Something bad happened, but we don't know what, which makes the fact much worse.")
		}

		e.Encode(struct {
			Status string `json:"status"`
			Error  string `json:"err"`
		}{"err", cr.Error.Error()})
	}
}
