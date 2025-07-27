/*
Package smudh provides functionality to parse and manage Short Message (SM) User Data Header (UDH) content within the SMPP protocol, handling both standalone text messages (without headers) and messages with UDH, including fragmentation. The package processes hexadecimal-encoded content and supports various encodings.

This package addresses issues with certain Short Message Service Centers (SMSCs) that fail to properly handle fragmented or non-fragmented text messages. If your SMSC correctly processes short messages and delivers proper text, this package may not be necessary, and you can rely on standard SMSC APIs.

When the short_message field contains raw hexadecimal UDH, the package provides a Message type (a byte slice) to represent the message. The ParseElements method detects and parses both UDH-structured and standalone messages.

UDH structure example:

		0500030F030368656C6C6F20776F726C64
	    | | | | | | |
	    | | | | | | |- The text to use (hello world in ASCII/GSM03.38/UTF-8 encodings).
	    | | | | | |- Current Part number (03 in this case)
	    | | | | |- Total Parts (03 in this case)
	    | | | |- Reference Number (can be either single or multi-byte long, 0F in this case)
	    | | |- Element Length (for reference number)
	    | |- Element (type)
	    |- Header Length

Standalone text example:

	776F726C64

Represents the word "world" (in ASCII/GSM03.38/UTF-8) without UDH.

The parsing of both for UDH and stand alone are detected and parsed using the ParseElements method.

UDH and standalone messages do not include encoding details, which must be provided via another SMPP field accompanying the `short_message`.

The package uses functional naming for elements rather than official UDH terminology.
*/
package smudh

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
