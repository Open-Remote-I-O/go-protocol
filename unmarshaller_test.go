package goprotocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCases struct {
	name                string
	mockProtocolRequest io.Reader
	expected            *OrioPayload
	expectedError       error
}

func generateMockProtocolBuffer(t testing.TB) io.Reader {
	var w bytes.Buffer
	err := binary.Write(&w, binary.BigEndian, version)
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	return &w
}

func Test_Unmarshal(t *testing.T) {
	tests := []testCases{
		{
			name:                "version ok",
			mockProtocolRequest: generateMockProtocolBuffer(t),
			expected: &OrioPayload{
				Header: Header{
					Version: uint16(version),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Unmarshal(tt.mockProtocolRequest)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.Equal(t, tt.expected, res)
		})
	}
}

// TODO: verify the best way in order to benchmark unmarhsalling behaivour
func Benchmark_Unmarshal(b *testing.B) {
	val := generateMockProtocolBuffer(b)
	for n := 0; n < b.N; n++ {
		_, _ = Unmarshal(val)
	}
}
