// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package parser // import "miniflux.app/reader/parser"

import (
	"encoding/xml"
	"strings"

	"miniflux.app/logger"
	"miniflux.app/reader/encoding"
)

// List of feed formats.
const (
	FormatRDF     = "rdf"
	FormatRSS     = "rss"
	FormatAtom    = "atom"
	FormatJSON    = "json"
	FormatUnknown = "unknown"
)

// DetectFeedFormat tries to guess the feed format from input data.
func DetectFeedFormat(data string) string {
	if strings.HasPrefix(strings.TrimSpace(data), "{") {
		return FormatJSON
	}

	data = stripInvalidXMLCharacters(data)

	decoder := xml.NewDecoder(strings.NewReader(data))
	decoder.CharsetReader = encoding.CharsetReader

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		if element, ok := token.(xml.StartElement); ok {
			switch element.Name.Local {
			case "rss":
				return FormatRSS
			case "feed":
				return FormatAtom
			case "RDF":
				return FormatRDF
			}
		}
	}

	return FormatUnknown
}

func stripInvalidXMLCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if isInCharacterRange(r) {
			return r
		}

		logger.Debug("Strip invalid XML characters: %U", r)
		return -1
	}, input)
}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of http://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}
