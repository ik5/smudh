// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
package smudh

import "fmt"

type Encoding byte

// UDH Message Encoding
const (
	// gsm-7bit encoding
	GSM Encoding = iota

	// ASCII/IA5 encoding
	ASCII

	// 8-bit binary encoding
	Binary8Bit1

	// iso-8859-1 encoding
	Latin1

	// 8-bit binary encoding
	Binary8Bit2

	// JIS (X 0208-1990) encoding
	JIS

	// iso-8859-5 encoding
	Cyrillic

	// iso-8859-8 encoding
	Hebrew

	// UCS2 (UTF-16BE) encoding
	UCS2

	// Pictogram - Cellular icons encoding support
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

	// support for UTF8 - hardly used
	UTF8
)

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
