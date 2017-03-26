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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"strings"

	dfi "github.com/dfindex/dfi"
	data "github.com/dfindex/dfi/data"
	dht "github.com/dfindex/dfi/dht"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// these two are inserted by the makefile at build time
var (
	Version   = "N/A"
	BuildTime = "N/A"
)

func SetupLocalPeer(addr string) *dfi.LocalPeer {
	var lp dfi.LocalPeer

	if lp.ReadKey() != nil {
		lp.GenerateKey()
		lp.WriteKey()
	}
	lp.Setup()

	return &lp
}

func main() {

	log.SetLevel(log.DebugLevel)
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "15:04:05"
	log.SetFormatter(formatter)

	os.Mkdir("./data", 0777)

	SetupConfig()

	addr := viper.GetString("bind.dfi")
	fmt.Println(addr)

	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])

	lp := SetupLocalPeer(fmt.Sprintf("%s", addr))
	lp.LoadEntry()

	log.WithFields(log.Fields{
		"version": Version,
		"built":   BuildTime,
	}).Info("Starting dfid")

	if viper.GetBool("tor.enabled") {
		_, onion, err := dfi.SetupDFITorService(port, viper.GetInt("tor.control"),
			fmt.Sprintf("%s/cookie", viper.GetString("tor.cookiePath")))

		if err == nil {
			lp.PublicAddress = onion
			lp.Entry.PublicAddress = onion
			lp.SetSocks(true)
			lp.SetSocksPort(viper.GetInt("tor.socks"))
			lp.Peer.Streams().Socks = true
			lp.Peer.Streams().SocksPort = viper.GetInt("tor.socks")
		} else {
			panic(err)
		}

		// should this override tor?
	} else if viper.GetBool("socks.enabled") {
		lp.SetSocks(true)
		lp.SetSocksPort(viper.GetInt("socks.port"))
		lp.Peer.Streams().Socks = true
		lp.Peer.Streams().SocksPort = viper.GetInt("socks.port")

		// TODO: configurable public address
	} else {
		if lp.Entry.PublicAddress == "" {
			log.Debug("Local peer public address is nil, attempting to fetch")
			ip := dfi.ExternalIp()
			log.Debug("External IP is ", ip)
			lp.Entry.PublicAddress = ip
		}
	}

	lp.Entry.Port = port
	lp.Entry.SetLocalPeer(lp)
	lp.SignEntry()
	lp.SaveEntry()

	err := lp.SaveEntry()

	if err != nil {
		panic(err)
	}

	lp.Database = data.NewDatabase(viper.GetString("database.path"))

	err = lp.Database.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	lp.Listen(viper.GetString("bind.dfi"))

	log.Info("My name: ", lp.Entry.Name)
	s, _ := lp.Address().String()
	log.Info("My address: ", s)

	commandServer := dfi.NewCommandServer(lp)
	var httpServer dfi.HttpServer
	httpServer.CommandServer = commandServer
	go httpServer.ListenHttp(viper.GetString("bind.http"))

	err = lp.StartExploring()

	if err != nil {
		log.Error(err.Error())
	}

	addr1 := dht.Address{
		Raw: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	addr1str, _ := addr1.String()

	fmt.Println(addr1str)

	// Listen for SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for _ = range sigchan {
		lp.Close()

		os.Exit(0)
	}
}
