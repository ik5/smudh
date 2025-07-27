# smudh

The `smudh` package provides functionality for handling SMPP `short_message` User Data Header (UDH) content in Go. It parses and manages both standalone text messages (without UDH) and messages with UDH, supporting fragmentation and various encodings. This package is designed for scenarios where Short Message Service Centers (SMSCs) fail to properly process fragmented or non-fragmented SMS messages, enabling developers to parse and reassemble hexadecimal-encoded message content.

## Features

- Parses hexadecimal-encoded SMPP `short_message` fields with or without UDH.
- Supports fragmented messages for reassembly of multi-part SMS.
- Handles common SMS encodings (e.g., GSM 7-bit, ASCII, UTF-16BE, ISO-8859-1).
- Provides JSON serialization/deserialization for message elements and fragments.
- Includes clear error handling for invalid inputs or unsupported encodings.

## Installation

To use `smudh` in your Go project, install it via:

```bash
go get github.com/ik5/smudh
```

Then, import it in your code:

```go
import "github.com/ik5/smudh"
```

## Requirements

- Go 1.18 or later.
- Familiarity with SMPP protocol and UDH structure (see [User Data Header](https://en.wikipedia.org/wiki/User_Data_Header)).

## Usage

The `smudh` package provides types and methods to parse and manage SMS messages:

- **`Message`**: A byte slice representing a hex-encoded SMS message.
- **`MessageElements`**: A struct holding parsed UDH components or standalone message details.
- **`MessageFragmentations`**: A slice of `MessageElements` for handling fragmented messages.
- **`Messages`**: A container for grouping message fragments by reference number.

Key methods include:
- `Message.ParseElements(encoding Encoding)`: Parses a message into its components.
- `MessageFragmentations.Add(encoding Encoding, message Message)`: Adds a message to a fragment collection.
- `Messages.GetMessageFragments(reference []byte)`: Retrieves all fragments for a given reference number.


## Contributing

Contributions are welcome! Please submit issues or pull requests to the [GitHub repository](https://github.com/ik5/smudh). Ensure your code follows Go conventions and includes tests where applicable.

## License

This project is licensed under the Mozilla Public License Version 2.0. See the [LICENSE](https://github.com/ik5/smudh/blob/main/LICENSE) file or visit [http://mozilla.org/MPL/2.0/](http://mozilla.org/MPL/2.0/) for details.

## TODO

- [ ] Add additional examples (`example_test.go`).
- [ ] Add additional testing to make sure all encoding and other methods works well.
- [ ] Write Benchmark to find better memory optimizations.
- [ ] When stable, raise to v1.0.0
