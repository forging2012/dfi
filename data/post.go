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
package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"time"
)

const (
	TitleMax    = 144
	TagsMax     = 256
	MaxPostSize = TitleMax + TagsMax + 1024
)

type Post struct {
	Id         int
	InfoHash   string
	Title      string
	Size       int
	FileCount  int
	Seeders    int
	Leechers   int
	UploadDate int
	Tags       string
	Meta       string
}

func (p Post) Json() ([]byte, error) {
	json, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func (p *Post) Bytes(sep, term []byte, seedLeech bool) []byte {
	buf := bytes.Buffer{}

	p.Write(string(sep), string(term), seedLeech, &buf)

	return buf.Bytes()
}

func (p *Post) String(sep, term string, seedLeech bool) string {
	return string(p.Bytes([]byte(sep), []byte(term), seedLeech))
}

// Includes an option to include seed/leech or not. Other than seed and leech
// count, posts are immutable. This prevents other peers from changing their data.
func (p *Post) Write(sep, term string, seedLeech bool, w io.Writer) {
	w.Write([]byte(strconv.Itoa(p.Id)))
	w.Write([]byte(sep))
	w.Write([]byte(p.InfoHash))
	w.Write([]byte(sep))
	w.Write([]byte(p.Title))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.Size)))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.FileCount)))

	if seedLeech {
		w.Write([]byte(sep))
		w.Write([]byte(strconv.Itoa(p.Seeders)))
		w.Write([]byte(sep))
		w.Write([]byte(strconv.Itoa(p.Leechers)))
	}

	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.UploadDate)))
	w.Write([]byte(sep))
	w.Write([]byte(p.Tags))
	w.Write([]byte(sep))
	w.Write([]byte(p.Meta))
	w.Write([]byte(sep))
	w.Write([]byte(term))

	/*
		The above seems to be a little faster, though mildly more awkward code.
		I suppose because it avoids allocating a buffer every write?
		bw := bufio.NewWriter(w)
		bw.WriteString(strconv.Itoa(p.Id))
		bw.WriteString(sep)
		bw.WriteString(p.InfoHash)
		bw.WriteString(sep)
		bw.WriteString(p.Title)
		bw.WriteString(sep)
		bw.WriteString(strconv.Itoa(p.Size))
		bw.WriteString(sep)
		bw.WriteString(strconv.Itoa(p.FileCount))
		bw.WriteString(sep)
		bw.WriteString(strconv.Itoa(p.Seeders))
		bw.WriteString(sep)
		bw.WriteString(strconv.Itoa(p.Leechers))
		bw.WriteString(sep)
		bw.WriteString(strconv.Itoa(p.UploadDate))
		bw.WriteString(sep)
		bw.WriteString(p.Tags)
		bw.WriteString(sep)
		bw.WriteString(term)
		bw.Flush()*/
}

func (p *Post) Valid() error {
	if len(p.Title) > 140 {
		return errors.New("Title too long")
	}

	if p.UploadDate > int(time.Now().Unix()) {
		return errors.New("Upload data cannot be in the future")
	}

	return nil
}
