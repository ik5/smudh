package smudh

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/ik5/gostrutils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Message represents a hex-encoded SMS message as a byte slice.
type Message []byte

// MessageElements represents a single parsed SMS message, either standalone or part of a UDH.
type MessageElements struct {
	// UDHL - UDH Length
	HeaderLength byte `json:"header_length"`

	// IEI (Information Element Identifier)
	Element byte `json:"element"`

	// IE Length (Length of reference number)
	ElementLength byte `json:"element_length"`

	// Reference Number (single or multi-byte)
	Reference []byte `json:"reference"`

	// Total number of parts
	TotalParts byte `json:"total_parts"`

	// Current part number
	CurrentPart byte `json:"current_part"`

	// Raw message payload
	RawMessage []byte `json:"raw_message"`

	// Decoded UTF-8 message
	Message string `json:"message"`

	// Message encoding
	Encoding Encoding `json:"encoding"`

	// True if message is standalone
	Standalone bool `json:"standalone"`
}

// MessageFragmentations a slice container of MessageElements pointers - for fragmentation gathering for a specific message.
type MessageFragmentations []*MessageElements

// Messages manages a collection of message fragmentations, grouped by reference number.
type Messages struct {
	fragments map[string]*MessageFragmentations
	mtx       sync.Mutex
}

const rfc822Element byte = 0x20

// ParseElements parses the hexadecimal content of a Message into its structural components, using the provided
// encoding from the SMPP protocol.
// On success, it returns a MessageElements struct.
// For standalone text (no UDH), the Standalone flag is set to true, TotalParts and CurrentPart are set to 1, and
// Reference is set to `0x00`.
// Returns an error for invalid content.
func (msg Message) ParseElements(encoding Encoding) (*MessageElements, error) {
	if len(msg)%2 != 0 {
		return nil, ErrHexStringMustHaveAnEvenNumberOfChars
	}

	binary, err := hex.DecodeString(string(msg))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var elements MessageElements
	elements.Encoding = encoding

	if len(binary) >= 2 {
		tmpLength := int(binary[0])
		if tmpLength > 0 && tmpLength < len(binary)-1 && binary[1] < rfc822Element {
			if tmpLength+1 > len(binary) {
				return nil, ErrUDHLengthExceedsInputLength
			}
			elements.HeaderLength = binary[0]
			elements.Element = binary[1]
			elements.ElementLength = binary[2]
			switch elements.Element {
			case 0x00: // 8-bit reference
				elements.Reference = []byte{binary[3]}
				elements.TotalParts = binary[4]
				elements.CurrentPart = binary[5]
			case 0x08: // 16-bit reference
				if tmpLength < 6 { // Need at least 6 bytes for UDH
					return nil, ErrInputTooShortForUDH
				}
				elements.Reference = binary[3:5] // 2 bytes
				elements.TotalParts = binary[5]
				elements.CurrentPart = binary[6]
			default:
				return nil, ErrUnsupportedIEI
			}

			elements.RawMessage = binary[tmpLength+1:]
		} else {
			elements.Standalone = true
			elements.Reference = []byte{0}
			elements.TotalParts = 0x01
			elements.CurrentPart = 0x01
			elements.RawMessage = binary
		}
	}

	err = elements.encodeMessage()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &elements, nil
}

// setTransformCharmap translates the given RawMessage based on a given decoder.
// If successful, than the function sets the elem.Message, otherwise an error is returned.
func (elem *MessageElements) setTransformCharmap(decoder *encoding.Decoder) error {
	var (
		err       error
		reader    *transform.Reader
		utf8Bytes []byte
	)
	reader = transform.NewReader(strings.NewReader(string(elem.RawMessage)), decoder)
	utf8Bytes, err = io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	elem.Message = string(utf8Bytes)
	return nil
}

// encodeMessage looks over the encoding element, and try to decode the RawMessage element accordingly.
//
// At this time the Pictogram encoding is not supported, as well as the Reserved1 and Reserved2 encoding.
// If found, an error will return.
// If the encoding is unknown, then an error is returned on that.
// Any other error is based on the encoding decoder streaming.
func (elem *MessageElements) encodeMessage() error {
	var (
		decoder *encoding.Decoder
		err     error
	)

	switch elem.Encoding {
	case GSM, GSMExtended:
		elem.Message = gostrutils.GSM0338ToUTF8(string(elem.RawMessage))

	case ASCII, UTF8:
		elem.Message = string(elem.RawMessage)

	case Latin1:
		decoder = charmap.ISO8859_1.NewDecoder()

		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case Binary8Bit1, Binary8Bit2:
		elem.Message = hex.EncodeToString(elem.RawMessage)
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case UCS2:
		if len(elem.RawMessage)%2 != 0 {
			return ErrBinaryTextLengthIsNotEvenForUTF16Decoding
		}

		decoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case Cyrillic:
		decoder = charmap.ISO8859_5.NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case Hebrew:
		decoder = charmap.ISO8859_8.NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case ISO2022JP:
		decoder = japanese.ISO2022JP.NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case KSC5601:
		decoder = korean.EUCKR.NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case JIS, EXTJIS:
		decoder = japanese.EUCJP.NewDecoder()
		err = elem.setTransformCharmap(decoder)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	case Pictogram, Reserved1, Reserved2:
		// TODO: support these as well
		return ErrUnsupportedEncoding

	default:
		return ErrUnknownEncoding

	}

	return nil
}

// IsSingleMessage returns true when message is standalone or is not fragmented.
func (elem MessageElements) IsSingleMessage() bool {
	return elem.Standalone || elem.TotalParts == 1
}

// ToJSON Serializes a MessageElements struct to JSON.
// Returns an error if serialization fails.
func (elem MessageElements) ToJSON() (string, error) {
	result, err := json.Marshal(elem)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(result), nil
}

// MessageElementFromJSON de-serializes a JSON string into a MessageElements struct.
// Returns an error if the JSON is invalid.
func MessageElementFromJSON(rawJSON string) (*MessageElements, error) {
	var result MessageElements

	err := json.Unmarshal([]byte(rawJSON), &result)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &result, nil
}

// ToJSON serializes the MessageFragmentations slice to JSON.
// Returns an error if serialization fails.
func (msgs MessageFragmentations) ToJSON() (string, error) {
	result, err := json.Marshal(msgs)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(result), nil
}

// FromJSON de-serialize a JSON string into a MessageFragmentations slice. Returns an error if the JSON is invalid.
func (msgs *MessageFragmentations) FromJSON(rawJSON string) error {
	err := json.Unmarshal([]byte(rawJSON), msgs)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil

}

// Sort sorts the MessageFragmentations slice in ascending order based on CurrentPart.
func (msgs MessageFragmentations) Sort() {
	slices.SortFunc(msgs, func(a, b *MessageElements) int {
		if b.Element > a.Element {
			return -1
		}
		if a.Element > b.Element {
			return 1
		}

		return 0
	})
}

// HaveAllFragments returns true if the MessageFragmentations contains all parts of a fragmented message or is standalone.
func (msgs MessageFragmentations) HaveAllFragments() bool {
	msgsLen := len(msgs)
	if msgsLen == 0 {
		return false
	}

	first := msgs[0]

	if first.Standalone {
		return true
	}

	if first.TotalParts == 0 && first.CurrentPart == 0 && len(first.Message) > 0 {
		return true
	}

	return msgsLen == int(first.TotalParts)
}

// String returns a string representation of the full oredered MessageFragmentations.
//
// IMPORTANT: The function calls Sort method before collecting all of the messages.
func (msgs *MessageFragmentations) String() string {
	msgs.Sort()

	buffer := bytes.Buffer{}

	for _, info := range *msgs {
		_, _ = buffer.WriteString(info.Message)
	}

	return buffer.String()
}

// Add parses a raw Message using the specified encoding and appends the resulting MessageElements to the MessageFragmentations slice.
// The method does not reorder elements. Returns an error if parsing fails.
func (msgs *MessageFragmentations) Add(encoding Encoding, message Message) error {
	info, err := message.ParseElements(encoding)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = msgs.AddMessageElements(info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// AddMessageElements appends a MessageElements instance to the MessageFragmentations slice.
// The method does not reorder elements. Returns an error if addition fails.
func (msgs *MessageFragmentations) AddMessageElements(info *MessageElements) error {
	if len(*msgs) == 0 {
		*msgs = append(*msgs, info)
		return nil
	}

	first := (*msgs)[0]

	if bytes.Equal(first.Reference, info.Reference) {
		*msgs = append(*msgs, info)
		return nil
	}

	return ErrInvalidReferenceNumber
}

// Reference returns the reference number of the message fragments.
// Returns nil if the slice is empty.
func (msgs MessageFragmentations) Reference() []byte {
	if len(msgs) == 0 {
		return nil
	}

	first := msgs[0]
	return first.Reference
}

// InitMessages	initializes and returns a new Messages instance.
func InitMessages() *Messages {
	messages := &Messages{
		fragments: make(map[string]*MessageFragmentations),
		mtx:       sync.Mutex{},
	}

	return messages
}

// AddMessageElements adds a MessageElements instance to the Messages container.
// Returns an error if the addition is invalid.
// The function does not re-order the elements.
func (msgs *Messages) AddMessageElements(info *MessageElements) error {
	msgs.mtx.Lock()
	defer msgs.mtx.Unlock()

	var err error

	strRefer := string(info.Reference)

	if _, found := msgs.fragments[strRefer]; found {
		err = msgs.fragments[strRefer].AddMessageElements(info)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	fragments := &MessageFragmentations{}
	err = fragments.AddMessageElements(info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	msgs.fragments[strRefer] = fragments

	return nil
}

// Add Parses a raw Message using the specified encoding and adds it to the Messages container.
// Returns an error if parsing fails.
// The function does not re-order the elements.
func (msgs *Messages) Add(encoding Encoding, message Message) error {
	msgs.mtx.Lock()
	defer msgs.mtx.Unlock()

	info, err := message.ParseElements(encoding)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	strRefer := string(info.Reference)

	if _, found := msgs.fragments[strRefer]; found {
		err = msgs.fragments[strRefer].AddMessageElements(info)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	fragments := &MessageFragmentations{}
	err = fragments.AddMessageElements(info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	msgs.fragments[strRefer] = fragments

	return nil
}

// GetMessageFragments retrieves the MessageFragmentations for a given reference number, returning an ordered slice.
// Returns nil if the reference is not found.
func (msgs *Messages) GetMessageFragments(reference []byte) *MessageFragmentations {
	msgs.mtx.Lock()
	defer msgs.mtx.Unlock()

	messages, found := msgs.fragments[string(reference)]
	if !found {
		return nil
	}

	messages.Sort()
	return messages
}

// ListAll returns a slice of all MessageFragmentations in the Messages container, unsorted.
func (msgs *Messages) ListAll() []*MessageFragmentations {
	msgs.mtx.Lock()
	defer msgs.mtx.Unlock()

	results := []*MessageFragmentations{}

	for _, fragmentations := range msgs.fragments {
		results = append(results, fragmentations)
	}

	return results
}
