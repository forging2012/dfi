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
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// TODO: Make this check using UpNp/NAT_PMP first, then query services.
func ExternalIp() string {
	resp, err := http.Get("https://api.ipify.org/")

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	return string(ret)
}
