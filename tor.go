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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	torc "github.com/postfix/goControlTor"
	log "github.com/sirupsen/logrus"
)

func SetupZifTorService(port, tor int, cookie string) (*torc.TorControl, string, error) {
	control := &torc.TorControl{}

	serviceDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Info(serviceDir)
	servicePort := map[int]string{port: fmt.Sprintf("127.0.0.1:%d", port)}

	err := control.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tor))

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Dialed Tor control port")

	err = control.CookieAuthenticate(cookie)

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Authenticated with Tor, creating service")

	err = control.CreateHiddenService(serviceDir, servicePort)

	if err != nil {
		log.Error(err.Error())
		return nil, "", err
	}

	log.Info("Service created")

	onion, err := torc.ReadOnion(serviceDir)
	onion = strings.TrimSpace(onion)

	log.WithField("onion", onion).Info("Connecting to Tor")

	return control, onion, nil
}
