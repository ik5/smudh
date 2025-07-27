package smudh_test

import (
	"fmt"
	"reflect"

	"github.com/ik5/smudh"
)

func ExampleMessage_udh() {
	fmt.Printf("%s\n", smudh.Message("05000312010168656C6C6F20776F726C64"))
}

func ExampleMessage_standalone() {
	fmt.Printf("%s\n", smudh.Message("776F726C64"))
}

func ExampleMessage_ParseElements() {
	msg := smudh.Message("05000312010168656C6C6F20776F726C64")
	elements, err := msg.ParseElements(smudh.GSM)
	if err != nil {
		panic(err)
	}

	expectedElements := smudh.MessageElements{
		HeaderLength:  0x05,
		Element:       0x00,
		ElementLength: 0x03,
		Reference:     []byte{0x12},
		TotalParts:    0x01,
		CurrentPart:   0x01,
		Encoding:      smudh.GSM,
		RawMessage:    []byte("hello world"),
		Message:       "hello world",
	}

	if !reflect.DeepEqual(*elements, expectedElements) {
		fmt.Println("Error - expected elements and parsed elements are not the same")
	}
}

func ExampleMessageFragmentations_Add() {
	fragmentation := smudh.MessageFragmentations{}

	// 2nd fragmentation
	err := fragmentation.Add(smudh.GSM, smudh.Message("050003A5020265722074657374696E67"))
	if err != nil {
		panic(err)
	}

	// 1st fragmentation
	err = fragmentation.Add(smudh.GSM, smudh.Message("050003A50201546869732069732061206C6F6E676572206D6573736167652074686174206E6565647320746F2062652073706C697420696E746F206D756C7469706C6520706172747320746F2064656D6F6E73747261746520534D5320636F6E636174656E6174696F6E20696E20534D50502070726F746F636F6C20776974682047534D20372D62697420656E636F64696E6720666F722070726F70"))
	if err != nil {
		panic(err)
	}

	// should print the full text in order:
	// This is a longer message that needs to be split into multiple parts to demonstrate SMS concatenation in SMPP protocol with GSM 7-bit encoding for proper testing
	fmt.Printf("%s\n", fragmentation.String())

}
