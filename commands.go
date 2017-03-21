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
package zif

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/zif/zif/data"
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
