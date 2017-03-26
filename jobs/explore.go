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
package jobs

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/dfindex/dfi/common"
	"github.com/dfindex/dfi/dht"
)

const ExploreFrequency = time.Minute * 2
const ExploreBufferSize = 100

// This job runs every two minutes, and tries to build the netdb with as many
// entries as it possibly can
func ExploreJob(in chan dht.Entry, data ...interface{}) <-chan dht.Entry {
	ret := make(chan dht.Entry, ExploreBufferSize)

	connector := data[0].(func(dht.Address) (interface{}, error))
	me := data[1].(dht.Address)
	seed := data[2].(func(ret chan dht.Entry))

	ticker := time.NewTicker(ExploreFrequency)

	go exploreTick(in, ret, me, connector, seed)

	go func() {
		for _ = range ticker.C {
			go exploreTick(in, ret, me, connector, seed)
		}

	}()

	return ret
}

func exploreTick(in chan dht.Entry, ret chan dht.Entry, me dht.Address, connector common.ConnectPeer, seed func(chan dht.Entry)) {
	i := <-in
	s, _ := i.Address.String()

	if i.Address.Equals(&me) {
		return
	}

	log.WithField("peer", s).Info("Exploring")

	if err := explorePeer(i.Address, me, ret, connector); err != nil {
		log.Error(err.Error())
	}

	if len(in) == 0 {
		seed(in)
		log.Info("Seeding peer explore")
	}
}

func explorePeer(addr dht.Address, me dht.Address, ret chan<- dht.Entry, connectPeer common.ConnectPeer) error {
	peer, err := connectPeer(addr)
	p := peer.(common.Peer)

	if err != nil {
		return err
	}

	randAddr, err := dht.RandomAddress()

	if err != nil {
		return err
	}

	log.Debug("Exploring random")
	closest, err := p.FindClosest(*randAddr)

	if err != nil {
		return err
	}

	for _, i := range closest {
		if !i.(*dht.Entry).Address.Equals(&me) {
			ret <- *(i.(*dht.Entry))
		}
	}

	log.Debug("Exploring closest to self")
	closestToMe, err := p.FindClosest(me)

	if err != nil {
		return err
	}
	log.Debug("Explored closest")

	for _, i := range closestToMe {
		if !i.(*dht.Entry).Address.Equals(&me) {
			ret <- *(i.(*dht.Entry))
		}
	}

	return nil
}
