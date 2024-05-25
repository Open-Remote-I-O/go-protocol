package goprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const version = uint16(1)

func readBytes(reader io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}
	if n != len(buf) {
		return nil, fmt.Errorf("unexpected data length, expected %d bytes, got %d", len(buf), n)
	}
	return buf, nil
}

func Unmarshal(reader io.Reader) (*OrioPayload, error) {
	var protocol OrioPayload

	// Read length field
	headerVersion, err := readBytes(reader, 2)
	if err != nil {
		return nil, err
	}
	if err := binary.Read(bytes.NewReader(headerVersion), binary.BigEndian, &protocol.Header.Version); err != nil {
		return nil, fmt.Errorf("error decoding length: %w", err)
	}

	return &protocol, nil
}
