// Package goprotocol had the message structure that the server expects for a specific version
// and Marshal and Unmarshal methos in order to generate and serialize protocol communication
package goprotocol

// NOTE: currently commented other protocol values in order to test basic implementation of the unmarshalling

const version = uint16(1)

const headerByteSize = 8

// Header has all metadata needed before handling actual protocol data
type Header struct {
	Version    uint16
	DeviceID   uint32
	PayloadLen uint16
}

// Data is the body sent with expected command and eventual data in order to give detail about command
type Data struct {
	// CommandId uint8
	// Len       uint16
	// Data      []byte
}

// OrioPayload is the complete payload that a client will be sending to server
type OrioPayload struct {
	Header Header
	// Data   Data
}
