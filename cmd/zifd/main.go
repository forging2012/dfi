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
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"strings"

	"github.com/spf13/viper"
	zif "github.com/zif/zif"
	data "github.com/zif/zif/data"
	dht "github.com/zif/zif/dht"

	log "github.com/sirupsen/logrus"
)

// these two are inserted by the makefile at build time
var (
	Version   = "N/A"
	BuildTime = "N/A"
)

func SetupLocalPeer(addr string) *zif.LocalPeer {
	var lp zif.LocalPeer

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

	addr := viper.GetString("bind.zif")
	fmt.Println(addr)

	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])

	lp := SetupLocalPeer(fmt.Sprintf("%s", addr))
	lp.LoadEntry()

	log.WithFields(log.Fields{
		"version": Version,
		"built":   BuildTime,
	}).Info("Starting zifd")

	if viper.GetBool("tor.enabled") {
		_, onion, err := zif.SetupZifTorService(port, viper.GetInt("tor.control"),
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
			ip := zif.ExternalIp()
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

	lp.Listen(viper.GetString("bind.zif"))

	log.Info("My name: ", lp.Entry.Name)
	s, _ := lp.Address().String()
	log.Info("My address: ", s)

	commandServer := zif.NewCommandServer(lp)
	var httpServer zif.HttpServer
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
