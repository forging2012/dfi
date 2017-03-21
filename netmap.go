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
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zif/zif/dht"
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
