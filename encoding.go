package smudh

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import "fmt"

// Encoding define a unique SMPP text encoding code.
type Encoding byte

// Encoding defines the encoding type for message content.
const (
	// GSM 7-bit encoding
	GSM Encoding = iota

	// ASCII/IA5 encoding
	ASCII

	// 8-bit binary encoding
	Binary8Bit1

	// ISO-8859-1 encoding
	Latin1

	// 8-bit binary encoding
	Binary8Bit2

	// JIS (X 0208-1990) encoding
	JIS

	// ISO-8859-5 encoding
	Cyrillic

	// ISO-8859-8 encoding
	Hebrew

	// UCS2 (UTF-16BE) encoding
	UCS2

	// Cellular pictogram icons encoding support
	Pictogram

	// ISO-2022-JP (Music Codes) encoding
	ISO2022JP

	// Not commonly used
	Reserved1

	// Not commonly used
	Reserved2

	// Extended Kanji JIS (X 0212-1990) encoding
	EXTJIS

	// KS C 5601 encoding
	KSC5601

	//GSM 7-bit (extended) - GSM 7-bit with national language extensions
	GSMExtended

	// UTF-8 encoding (rarely used)
	UTF8
)

// Returns the string representation of the Encoding type.
func (enc Encoding) String() string {
	switch enc {
	case GSM:
		return "GSM-7"
	case ASCII:
		return "ASCII"
	case Binary8Bit1:
		return "BINARY-1"
	case Latin1:
		return "Latin1"
	case Binary8Bit2:
		return "BINARY-2"
	case JIS:
		return "JIS"
	case Cyrillic:
		return "ISO8859-5 (Cyrillic)"
	case Hebrew:
		return "ISO8859-8 (Hebrew)"
	case UCS2:
		return "UCS2 (UTF-16BE)"
	case Pictogram:
		return "Pictogram"
	case ISO2022JP:
		return "ISO2022JP (music)"
	case EXTJIS:
		return "EXT-JS (X 0212-1990)"
	case KSC5601:
		return "KSC-5601"
	case GSMExtended:
		return "GSM-7 (Extended)"
	case UTF8:
		return "UTF-8"
	}

	return fmt.Sprintf("%d", enc)
}
