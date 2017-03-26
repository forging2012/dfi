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
	"bytes"
	"time"

	"github.com/dfindex/dfi/dht"
	"github.com/dfindex/dfi/util"

	log "github.com/sirupsen/logrus"
)

const SeedSearchFrequency = time.Minute * 5

// brought into it's own type to track seed data, and manage it.
// works for all peers that we are seeding, including the localpeer.
// could one day be extended as a "gossip" protocol for stuff like comments,
// methinks.

type SeedManager struct {
	// the localpeer, allows the struct to make connections, etc
	lp *LocalPeer

	// the address we are tracking seeds for
	track dht.Address
	entry *dht.Entry
	Close chan bool
}

// Creates a new seed manager, given an address to track seeds for and the
// localpeer.
func NewSeedManager(track dht.Address, lp *LocalPeer) (*SeedManager, error) {
	ret := SeedManager{
		lp:    lp,
		Close: make(chan bool),
	}

	entry, err := lp.QueryEntry(track)

	if err != nil {
		return nil, err
	}

	ret.entry = entry
	ret.track = track

	return &ret, nil
}

// Start looking for seeds
func (sm *SeedManager) Start() {
	log.WithField("peer", sm.track.StringOr("")).Info("Starting seed manager")
	go sm.findSeeds()
}

// queries all seeds to see if we can find new seeds
func (sm *SeedManager) findSeeds() {
	ticker := time.NewTicker(SeedSearchFrequency)

	find := func() {
		entry, err := sm.lp.QueryEntry(sm.track)

		if err != nil {
			log.Error(err.Error())
			return
		}

		sm.entry = entry

		log.Info("Searching for new seeds")
		for _, i := range sm.entry.Seeds {
			addr := dht.Address{Raw: i}

			if addr.Equals(sm.lp.Address()) {
				continue
			}

			s, err := addr.String()
			if err != nil {
				continue
			}

			peer, entry, err := sm.lp.ConnectPeer(addr)

			if err != nil {
				log.Error(err.Error())
				continue
			}

			es, err := entry.Address.String()

			if err != nil {
				log.Error(err.Error())
				continue
			}

			log.WithField("peer", es).Info("Querying for seeds")

			qResultVerifiable, err := peer.Query(sm.entry.Address)
			if err != nil {
				continue
			}

			qResult := qResultVerifiable.(*dht.Entry)

			result := util.SliceDiff(sm.entry.Seeds, qResult.Seeds)

			// make sure all these seeds actually link back! Otherwise they could
			// be fakes
			for n, i := range result {
				seedAddress := dht.Address{Raw: i}

				entry, err := sm.lp.Resolve(seedAddress)

				// nope, we won't be adding this one
				if err != nil {
					if n >= len(result)-1 {
						result = result[:n]
					} else {
						result = append(result[:n], result[n+1:]...)
					}
					result = append(result[:n], result[n+1:]...)
					continue
				}

				// check if the entry has registered itself as a seeder

				found := false
				for _, j := range entry.Seeding {
					if bytes.Equal(sm.track.Raw, j) {
						found = true
						break
					}
				}

				if !found {
					if n >= len(result)-1 {
						result = result[:n]
					} else {
						result = append(result[:n], result[n+1:]...)
					}
					continue
				}
			}

			if len(result) > 0 {
				sm.entry.Seeds = append(sm.entry.Seeds, result...)

				log.WithField("peer", s).Info("Found new seeds")
				sm.lp.DHT.Insert(*sm.entry)

			}
		}
	}

	find()

	for {
		select {
		case _ = <-ticker.C:
			find()
		case _ = <-sm.Close:
			return
		}
	}
}
