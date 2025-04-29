package udh

import (
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	ErrHexStringMustHaveAnEvenNumberOfChars      = errors.New("hex string must have an even number of characters")
	ErrBinaryTextLengthIsNotEvenForUTF16Decoding = errors.New("binary text length is not even for UTF-16 decoding")
	ErrRuneOusideUCS2Range                       = errors.New("rune outside UCS-2 range")
	ErrInputTooShortForUDH                       = errors.New("input too short for UDH")
	ErrUDHLengthExceedsInputLength               = errors.New("UDH length exceeds input length")
)

type UDHMessage []byte

// DecodeUDHSMSMessage takes a UDH based binary UTF-16BE/UCS2 text and convert it into UTF-8 text
func DecodeUDHSMSMessagemsg(msg UDHMessage) (string, error) {
	var binary []byte
	var err error

	// Assume input is a hex string; decode to binary
	if len(msg)%2 != 0 {
		return "", ErrHexStringMustHaveAnEvenNumberOfChars
	}

	binary, err = hex.DecodeString(string(msg))
	if err != nil {
		return "", errors.New("invalid hex string: " + err.Error())
	}

	// Check if binary is long enough for UDH
	if len(binary) < 1 {
		return "", ErrInputTooShortForUDH
	}

	// Read UDH length (first byte indicates number of following bytes)
	udhLength := int(binary[0]) + 1 // Include length byte
	if udhLength > len(binary) {
		return "", ErrUDHLengthExceedsInputLength
	}

	// Extract text after UDH
	textBinary := binary[udhLength:]
	// Check if text length is even for UTF-16BE
	if len(textBinary)%2 != 0 {
		return "", ErrBinaryTextLengthIsNotEvenForUTF16Decoding
	}

	// Decode UTF-16BE to UTF-8
	decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	reader := transform.NewReader(strings.NewReader(string(textBinary)), decoder)
	utf8Bytes, err := io.ReadAll(reader)

	if err != nil {
		return "", errors.New("failed to decode UTF-16BE: " + err.Error())
	}

	return string(utf8Bytes), nil
}
