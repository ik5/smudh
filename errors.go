package udh

import "errors"

var (
	ErrHexStringMustHaveAnEvenNumberOfChars      = errors.New("hex string must have an even number of characters")
	ErrBinaryTextLengthIsNotEvenForUTF16Decoding = errors.New("binary text length is not even for UTF-16 decoding")
	ErrInputTooShortForUDH                       = errors.New("input too short for UDH")
	ErrUDHLengthExceedsInputLength               = errors.New("UDH length exceeds input length")
	ErrMessageNotComplete                        = errors.New("message is not complete yet")
	ErrMissingPart                               = errors.New("missing part")
	ErrInvalidReferenceNumber                    = errors.New("invalid reference number")
	ErrUnsupportedIEI                            = errors.New("unsupported IEI")
	ErrUnsupportedEncoding                       = errors.New("unsupported encoding")
	ErrUnknownEncoding                           = errors.New("unknown encoding")
)
