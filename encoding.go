package udh

type Encoding byte

// UDH Message Encoding
const (
	// gsm-7bit coding
	GSM Encoding = iota

	// ASCII/IA5 coding
	ASCII

	// 8-bit binary coding
	Binary8Bit1

	// iso-8859-1 coding
	Latin1

	// 8-bit binary coding
	Binary8Bit2

	// JIS (X 0208-1990) coding
	JIS

	// iso-8859-5 coding
	Cyrillic

	// iso-8859-8 coding
	Hebrew

	// UCS2 (UTF-16BE) coding
	UCS2

	// Pictogram - Cellular icons support
	Pictogram

	// ISO-2022-JP (Music Codes)
	ISO2022JP

	// Not commonly used
	Reserved1

	// Not commonly used
	Reserved2

	// Extended Kanji JIS (X 0212-1990)
	EXTJIS

	// KS C 5601
	KSC5601

	//GSM 7-bit (extended) - GSM 7-bit with national language extensions
	GSMEdtended

	// support for UTF8 - hardly used
	UTF8
)
