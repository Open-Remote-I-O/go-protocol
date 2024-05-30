package goprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"testing"

	errors "github.com/Open-Remote-I-O/orio-go-protocol/internal"
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
	err = binary.Write(&w, binary.BigEndian, rand.Uint32())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
	}
	binary.BigEndian.AppendUint32([]byte{}, rand.Uint32())
	err = binary.Write(&w, binary.BigEndian, uint16(rand.UintN(math.MaxUint16)))
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
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
					Version: version,
				},
			},
			expectedError: nil,
		},
		{
			name:                "header EOF err",
			mockProtocolRequest: bytes.NewBuffer([]byte{}),
			expected:            nil,
			expectedError:       fmt.Errorf("%s: %w", errors.ErrHeaderFormat, io.EOF),
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

func Fuzz_Unmarshal(f *testing.F) {
	// init fuzz corpus values

	var validHeaderVal bytes.Buffer
	err := binary.Write(&validHeaderVal, binary.BigEndian, version)
	err = binary.Write(&validHeaderVal, binary.BigEndian, uint32(32000))
	err = binary.Write(&validHeaderVal, binary.BigEndian, uint16(42))
	if err != nil {
		f.Fatal("something went wrong while creazint fuzz corpus values: %w", err)
	}

	fuzzCorpus := [][]byte{validHeaderVal.Bytes()}

	for _, v := range fuzzCorpus {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		_, err := Unmarshal(bytes.NewReader(b))
		if err != nil && err != errors.ErrHeaderFormatEOF {
			t.Errorf("given test case %v;\n caused: %s", b, err)
		}
	})
}

// TODO: verify the best way in order to benchmark unmarhsalling behaivour
func Benchmark_Unmarshal(b *testing.B) {
	val := generateMockProtocolBuffer(b)
	for n := 0; n < b.N; n++ {
		_, _ = Unmarshal(val)
	}
}
