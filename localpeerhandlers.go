package zif

import (
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const MaxSearchLength = 256

// TODO: Move this into some sort of handler object, can handle general requests.

// TODO: While I think about it, move all these TODOs to issues or a separate
// file/issue tracker or something.

// Querying peer sends a Zif address
// This peer will respond with a list of the k closest peers, ordered by distance.
// The top peer may well be the one that is being queried for :)
func (lp *LocalPeer) HandleQuery(msg *Message) error {
	log.Info("Handling query")
	cl := Client{msg.Stream, nil, nil}

	msg.From.limiter.queryLimiter.Wait()

	address := DecodeAddress(string(msg.Content))
	log.WithField("target", address.Encode()).Info("Recieved query")

	ok := &Message{Header: ProtoOk}
	err := cl.WriteMessage(ok)

	if err != nil {
		return err
	}

	var closest []*Entry

	if address.Equals(&lp.ZifAddress) {
		log.Debug("Query for local peer")
		closest = make([]*Entry, 1)
		closest[0] = &lp.Entry
	} else {
		log.Debug("Querying routing table")
		closest = lp.RoutingTable.FindClosest(address, BucketSize)
	}

	closest_json, err := json.Marshal(closest)
	log.Debug("Query results: ", string(closest_json))

	if err != nil {
		return errors.New("Failed to encode closest peers to json")
	}

	results := &Message{
		Header:  ProtoEntry,
		Content: closest_json,
	}

	err = cl.WriteMessage(results)

	return err
}

func (lp *LocalPeer) HandleAnnounce(msg *Message) error {
	cl := Client{msg.Stream, nil, nil}
	msg.From.limiter.announceLimiter.Wait()
	lp.CheckSessions()

	defer msg.Stream.Close()

	entry := Entry{}
	err := msg.Decode(&entry)

	log.WithField("address", entry.ZifAddress.Encode()).Info("Announce")

	if err != nil {
		return err
	}

	saved := lp.RoutingTable.Update(entry)

	if saved {
		cl.WriteMessage(&Message{Header: ProtoOk})
		log.WithField("peer", entry.ZifAddress.Encode()).Info("Saved new peer")

	} else {
		cl.WriteMessage(&Message{Header: ProtoNo})
		return errors.New("Failed to save entry")
	}

	// next up, tell other people!
	closest := lp.RoutingTable.FindClosest(entry.ZifAddress, BucketSize)

	// TODO: Parallize this
	for _, i := range closest {
		if i.ZifAddress.Equals(&entry.ZifAddress) || i.ZifAddress.Equals(&msg.From.ZifAddress) {
			continue
		}

		peer := lp.GetPeer(i.ZifAddress.Encode())

		if peer == nil {
			log.Debug("Connecting to new peer")

			var p Peer
			err = p.Connect(i.PublicAddress+":"+strconv.Itoa(i.Port), lp)

			if err != nil {
				log.Warn("Failed to connect to peer: ", err.Error())
				continue
			}

			p.ConnectClient(lp)

			peer = &p
		}

		peer_stream, err := peer.OpenStream()
		defer peer_stream.Close()

		if err != nil {
			log.Error(err.Error())
			continue
		}

		peer_announce := &Message{
			Header:  ProtoDhtAnnounce,
			Content: msg.Content,
		}
		peer_stream.WriteMessage(peer_announce)
	}
	return nil

}

func (lp *LocalPeer) HandleSearch(msg *Message) error {
	if len(msg.Content) > MaxSearchLength {
		return errors.New("Search query too long")
	}

	sq := MessageSearchQuery{}
	err := msg.Decode(&sq)

	if err != nil {
		return err
	}

	log.WithField("query", sq.Query).Info("Search recieved")

	posts, err := lp.Database.Search(sq.Query, sq.Page)

	if err != nil {
		return err
	}

	json, err := PostsToJson(posts)

	if err != nil {
		return err
	}

	post_msg := &Message{
		Header:  ProtoPosts,
		Content: json,
	}

	NewClient(msg.Stream).WriteMessage(post_msg)

	return nil
}

func (lp *LocalPeer) HandleRecent(msg *Message) error {
	log.Info("Recieved query for recent posts")

	page, err := strconv.Atoi(string(msg.Content))

	if err != nil {
		return err
	}

	recent, err := lp.Database.QueryRecent(page)

	if err != nil {
		return err
	}

	recent_json, err := PostsToJson(recent)

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoPosts,
		Content: recent_json,
	}

	NewClient(msg.Stream).WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandlePopular(msg *Message) error {
	log.Info("Recieved query for popular posts")

	page, err := strconv.Atoi(string(msg.Content))

	if err != nil {
		return err
	}

	recent, err := lp.Database.QueryPopular(page)

	if err != nil {
		return err
	}

	recent_json, err := PostsToJson(recent)

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoPosts,
		Content: recent_json,
	}

	NewClient(msg.Stream).WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandleHashList(msg *Message) error {
	cl := NewClient(msg.Stream)
	address := Address{msg.Content}

	log.WithField("address", address.Encode()).Info("Collection request recieved")

	col, err := CreateCollection(lp.Database, 0, 1000)

	sig := lp.Sign(col.HashList())
	log.Info(col.Hash())
	mhl := MessageCollection{col.Hash(), col.HashList(), len(col.HashList()) / 32, sig}
	data, err := mhl.Encode()

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoHashList,
		Content: data,
	}

	cl.WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandlePiece(msg *Message) error {
	cl := NewClient(msg.Stream)

	mrp := MessageRequestPiece{}
	err := msg.Decode(&mrp)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"address": mrp.Address,
		"id":      mrp.Id,
	}).Info("New piece request")

	var piece *Piece
	piece = nil

	if mrp.Address == lp.Entry.ZifAddress.Encode() {
		piece, err = lp.Database.QueryPiece(mrp.Id, true)

		if err != nil {
			return err
		}
	} else if _, ok := lp.Databases[mrp.Address]; ok {
		piece, err = lp.Databases[mrp.Address].QueryPiece(mrp.Id, true)

		if err != nil {
			return err
		}

	} else {
		return errors.New("Piece not found")
	}

	data, err := json.Marshal(piece)

	if err != nil {
		return err
	}

	rep := &Message{
		Header:  ProtoPiece,
		Content: data,
	}

	cl.WriteMessage(rep)

	return nil
}

func (lp *LocalPeer) ListenStream(peer *Peer) {
	lp.Server.ListenStream(peer)
}
