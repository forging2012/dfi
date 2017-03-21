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
package proto

func ChooseCompression(client MessageCapabilities, server MessageCapabilities) string {
	// check if the peer has our caps, in order of preference
	// the server has preference
	compression := ""
	for _, i := range server.Compression {
		if compression != "" {
			break
		}

		for _, j := range client.Compression {
			if i == j {
				compression = i
			}
		}
	}

	return compression
}
