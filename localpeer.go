package main

import (
	"github.com/streamrail/concurrent-map"
	"golang.org/x/crypto/ed25519"
)

type Identity struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Public string `json:"public"`
}

type LocalPeer struct {
	MirrorProgress cmap.ConcurrentMap

	Identity Identity

	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}
