package goprotocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	errors "github.com/Open-Remote-I-O/orio-go-protocol/internal"
)

// Reads from reader n bytes in newly instantiated n byte slice and return a reader of n bytes
func readBytes(reader io.Reader, n int) (*bytes.Reader, error) {
	buf := make([]byte, n)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}
	if n != len(buf) {
		return nil, fmt.Errorf("unexpected data length, expected %d bytes, got %d", len(buf), n)
	}
	return bytes.NewReader(buf), nil
}

// check if the buffered reader has the required number of bytes to be read from
// accepts the number of bytes to read and the reader itself
func checkValidLen(n int, buffrd *bufio.Reader) error {
	if _, err := buffrd.Peek(n); err != nil {
		return err
	}
	return nil
}

// Unmarshal Validates if r != EOF and has necessary bytes for each struct parameter and then
// deserialize the provided reader into the protocol payload struct
func Unmarshal(r io.Reader) (*OrioPayload, error) {
	var protocolHeader Header
	buffrdHeader := bufio.NewReaderSize(r, headerByteSize)
	err := checkValidLen(headerByteSize, buffrdHeader)
	if err != nil {
		if err == io.EOF {
			return nil, errors.ErrHeaderFormatEOF
		}
		return nil, fmt.Errorf("%s: %w", errors.ErrHeaderFormat, err)
	}
	headerVersion, err := readBytes(buffrdHeader, 2)
	if err != nil {
		return nil, fmt.Errorf("error reading version bytes: %w", err)
	}

	if err := binary.Read(headerVersion, binary.BigEndian, &protocolHeader.Version); err != nil {
		return nil, fmt.Errorf("error decoding length: %w", err)
	}

	deviceID, err := readBytes(buffrdHeader, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading deviceId bytes: %w", err)
	}

	if err := binary.Read(deviceID, binary.BigEndian, &protocolHeader.DeviceID); err != nil {
		return nil, fmt.Errorf("error decoding length: %w", err)
	}

	payloadLen, err := readBytes(buffrdHeader, 2)
	if err != nil {
		return nil, fmt.Errorf("error reading payloadLen bytes: %w", err)
	}

	if err := binary.Read(payloadLen, binary.BigEndian, &protocolHeader.PayloadLen); err != nil {
		return nil, fmt.Errorf("error decoding length: %w", err)
	}

	return &OrioPayload{
		Header: protocolHeader,
	}, nil
}
