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
	"errors"
	"fmt"
	"github.com/dfindex/dfi/dht"
	log "github.com/sirupsen/logrus"
)

type MapNode struct {
	// the address is treated like an id
	Address string `json:"id"`
	Name    string `json:"name"`
}

type MapLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// this file creates a JSON map of the network, compatible with d3.js

// This takes a node to start with, and recurses through all seeds/seeding
func CreateNetMap(entry dht.Entry, db *dht.DHT, currentNodes map[string]bool, currentLinks map[string]bool) ([]MapNode, []MapLink) {
	// BUG: Fix the duplicate nodes, luckily there are no duplicate links afaik.
	// Ensure that all links have associated nodes, and no duplicates.
	nodes := make([]MapNode, 0)
	links := make([]MapLink, 0)

	if _, ok := currentNodes[string(entry.Address.Raw)]; !ok {
		nodes = append(nodes, MapNode{Address: entry.Address.StringOr(""), Name: entry.Name})
	}

	createMap := func(i []byte) error {
		address := dht.Address{Raw: i}

		e, err := db.Query(address)

		if err != nil {
			log.Error(err)
			return err
		}

		if _, ok := currentLinks[string(e.Address.Raw)+string(entry.Address.Raw)]; !ok {
			currentLinks[string(e.Address.Raw)+string(entry.Address.Raw)] = true
			links = append(links, MapLink{Source: e.Address.StringOr(""), Target: entry.Address.StringOr("")})
		} else {
			return errors.New("continue")
		}

		if _, ok := currentNodes[string(e.Address.Raw)]; !ok {
			fmt.Println(e.Address.StringOr(""))
			currentNodes[string(e.Address.Raw)] = true
			nodes = append(nodes, MapNode{Address: e.Address.StringOr(""), Name: e.Name})
		} else {
			return errors.New("continue")
		}

		n, l := CreateNetMap(*e, db, currentNodes, currentLinks)

		nodes = append(nodes, n...)
		links = append(links, l...)

		return nil
	}

	for _, i := range entry.Seeding {
		err := createMap(i)

		if err != nil {
			continue
		}
	}

	return nodes, links
}
