package udh

import (
        "bytes"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "io"
        "slices"
        "strings"
        "sync"

        "github.com/ik5/gostrutils"
        "golang.org/x/text/encoding/unicode"
        "golang.org/x/text/transform"
)

type Encoding byte

// UDH Message Encoding
const (
        // gsm-7bit coding
        GSM Encoding = iota

        // ascii coding
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

        // Extended Kanji JIS (X 0212-1990)
        EXTJIS Encoding = 0x0D

        // KS C 5601
        KSC5601 Encoding = 0x0E
)

// Message represents a hex-encoded SMS message
type Message []byte
]
// Elements for SMS based UDH single message.
// This struct does not hold the protocol names, but rather it's uses
type MessageElements struct {
        // UDHL
        HeaderLength byte `json:"header_length"`

        // IEI (Information Element Identifier)
        Element byte `json:"element"`

        // IE Length
        ElementLength byte `json:"element_length"`

        // Reference Number
        Reference []byte `json:"reference"`

        // How many parts are there
        TotalParts byte `json:"total_parts"`

        // Current Part
        CurrentPart byte `json:"current_part"`

        // The actual message payload part
        RawMessage []byte `json:"raw_message"`

        // Decoded UTF-8 message
        Message string `json:"message"`

        // Encoding is not part of the UDH part
        Encoding Encoding `json:"encoding"`

        // is this message stand alone (pure message)
        Standalone bool `json:"standalone"`
}

// MessageFragmentations hold all fragmented value of a given message
type MessageFragmentations []*MessageElements

type Messages struct {
        fragments map[string]*MessageFragmentations
        mtx       sync.Mutex
}

const rfc822Element byte = 0x20

var (
        ErrHexStringMustHaveAnEvenNumberOfChars      = errors.New("hex string must have an even number of characters")
        ErrBinaryTextLengthIsNotEvenForUTF16Decoding = errors.New("binary text length is not even for UTF-16 decoding")
        ErrInputTooShortForUDH                       = errors.New("input too short for UDH")
        ErrUDHLengthExceedsInputLength               = errors.New("UDH length exceeds input length")
        ErrMessageNotComplete                        = errors.New("message is not complete yet")
        ErrMissingPart                               = errors.New("missing part")
        ErrInvalidReferenceNumber                    = errors.New("invalid reference number")
        ErrUnsupportedIEI                            = errors.New("unsupported IEI")
)
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

        switch encoding {
        case GSM:
                elements.Message = gostrutils.GSM0338ToUTF8(string(elements.RawMessage))
        case UCS2:
                if len(elements.RawMessage)%2 != 0 {
                        return nil, ErrBinaryTextLengthIsNotEvenForUTF16Decoding
                }

                decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
                reader := transform.NewReader(strings.NewReader(string(elements.RawMessage)), decoder)
                utf8Bytes, err := io.ReadAll(reader)
                if err != nil {
                        return nil, fmt.Errorf("%w", err)
                }

                elements.Message = string(utf8Bytes)
        }

        return &elements, nil
}

func (elem MessageElements) IsSingleMessage() bool {
        return elem.Standalone || elem.TotalParts == 1
}

func (elem MessageElements) ToJSON() (string, error) {
        result, err := json.Marshal(elem)

        if err != nil {
                return "", fmt.Errorf("%w", err)
        }

        return string(result), nil
}

func MessageElementFromJSON(rawJSON string) (*MessageElements, error) {
        var result MessageElements

        err := json.Unmarshal([]byte(rawJSON), &result)
        if err != nil {
                return nil, fmt.Errorf("%w", err)
        }

        return &result, nil
}

func (msgs MessageFragmentations) ToJSON() (string, error) {
        result, err := json.Marshal(msgs)

        if err != nil {
                return "", fmt.Errorf("%w", err)
        }

        return string(result), nil
}

func (msgs *MessageFragmentations) FromJSON(rawJSON string) error {
        err := json.Unmarshal([]byte(rawJSON), msgs)
        if err != nil {
                return nil, fmt.Errorf("%w", err)
        }

        return &result, nil

}

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
func (msgs MessageFragmentations) HaveAllFragments() bool {
        msgsLen := len(msgs)
        if msgsLen == 0 {
                return false
        }


        first := msgs[0]

        if first.TotalParts == 0 && first.CurrentPart == 0 && len(first.Message) > 0 {
                return true
        }

        return msgsLen == int(first.TotalParts)
}

func (msgs *MessageFragmentations) String() string {
        msgs.Sort()

        buffer := bytes.Buffer{}

        for _, info := range *msgs {
                _, _ = buffer.WriteString(info.Message)
        }

        return buffer.String()
}

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

func (msgs MessageFragmentations) Reference() []byte {
        if len(msgs) == 0 {
                return nil
        }

        first := msgs[0]
        return first.Reference
}
func InitMessages() *Messages {
        messages := &Messages{
                fragments: make(map[string]*MessageFragmentations),
                mtx:       sync.Mutex{},
        }

        return messages
}

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

func (msgs *Messages) GetMessageFragments(reference []byte) *MessageFragmentations {
        msgs.mtx.Lock()
        defer msgs.mtx.Unlock()

        messages, found := msgs.fragments[string(reference)]
        if !found {
                return nil
        }

        return messages
}

func (msgs *Messages) ListAll() []*MessageFragmentations {
        msgs.mtx.Lock()
        defer msgs.mtx.Unlock()

        results := []*MessageFragmentations{}

        for _, fragmentations := range msgs.fragments {
                results = append(results, fragmentations)
        }

        return results
}
