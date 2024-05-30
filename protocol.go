package goprotocol

// NOTE: currently commented other protocol values in order to test basic implementation of the unmarshalling

const headerByteSize = 8

type Header struct {
	Version    uint16
	DeviceId   uint32
	PayloadLen uint16
}

type Data struct {
	// CommandId uint8
	// Len       uint16
	// Data      []byte
}

type OrioPayload struct {
	Header Header
	// Data   Data
}
