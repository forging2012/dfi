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
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// This provides searching, as it is a little more comlex than just a db query.
// Search strings could provide other data that needs parsing, as well as spell
// correction that needs doing. This has to be passed through other functions
// before it hits a db query, hence this.
type SearchProvider struct {
	Loaded bool
	// if the model has been loaded, otherwise no autocomplete/spell suggestions
}

type SearchResult struct {
	Posts  []*Post `json:"posts"`
	Source string  `json:"source"`
}

func NewSearchProvider() *SearchProvider {
	sp := &SearchProvider{true}

	return sp
}

func IsAlnumWord(word string) bool {
	for _, i := range word {
		if !unicode.IsLetter(i) && !unicode.IsNumber(i) {
			return false
		}
	}

	return true
}

// Takes a string, makes it look "nice" for an autocomplete cue.
func SanitiseForAuto(in string) string {
	buffer := bytes.Buffer{}

	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if IsAlnumWord(scanner.Text()) {
			buffer.WriteString(scanner.Text())
			buffer.WriteString(" ")
		}
	}

	return buffer.String()
}

func (sp *SearchProvider) Suggest(db *Database, query string) ([]string, error) {
	checked, err := db.Suggest(fmt.Sprintf("%s%%", query))

	if err != nil {
		return nil, err
	}

	ret := make([]string, len(checked))

	for _, i := range checked {
		ret = append(ret, SanitiseForAuto(i))
	}

	return ret, nil
}

func (sp *SearchProvider) Search(source string, db *Database, query string, page int) (SearchResult, error) {
	// TODO: Instead of searching for spell-corrected versions, suggest an
	// alternate search.
	results, err := db.Search(query, page, 25)

	return SearchResult{results, source}, err
}
