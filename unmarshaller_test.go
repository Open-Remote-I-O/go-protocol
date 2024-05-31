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

type testCase struct {
	name          string
	version       []byte
	deviceID      []byte
	payloadLen    []byte
	expected      *OrioPayload
	expectedError error
}

func randUint16() uint16 {
	return uint16(rand.UintN(math.MaxUint16))
}

func genMockProtocolParamBytes[T any](t testing.TB, val T) []byte {
	var w bytes.Buffer
	err := binary.Write(&w, binary.BigEndian, val)
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	return w.Bytes()
}

func generateMockProtocolBuffer(t testing.TB, tt testCase) io.Reader {
	var w bytes.Buffer
	err := binary.Write(&w, binary.BigEndian, tt.version)
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	err = binary.Write(&w, binary.BigEndian, tt.deviceID)
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
	}
	err = binary.Write(&w, binary.BigEndian, tt.payloadLen)
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
	}
	return &w
}

func Test_Unmarshal(t *testing.T) {
	tests := []testCase{
		{
			name:          "ok",
			version:       genMockProtocolParamBytes(t, version),
			deviceID:      genMockProtocolParamBytes(t, uint32(10)),
			payloadLen:    genMockProtocolParamBytes(t, uint16(100)),
			expected:      &OrioPayload{Header: Header{Version: version, DeviceID: 10, PayloadLen: 100}},
			expectedError: nil,
		},
		{
			name:          "invalid header format passed",
			version:       []byte{},
			deviceID:      nil,
			payloadLen:    nil,
			expected:      nil,
			expectedError: fmt.Errorf("%s: %w", errors.ErrHeaderFormat, io.EOF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVal := generateMockProtocolBuffer(t, tt)
			res, err := Unmarshal(mockVal)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.Equal(t, tt.expected, res)
		})
	}
}

// TODO: learn more about fuzzing and best practices
func Fuzz_Unmarshal(f *testing.F) {
	// init fuzz corpus values

	validHeaderValReader := generateMockProtocolBuffer(f, testCase{
		version:    genMockProtocolParamBytes(f, version),
		deviceID:   genMockProtocolParamBytes(f, uint32(10)),
		payloadLen: genMockProtocolParamBytes(f, uint16(100)),
	})
	validHeaderVal, err := io.ReadAll(validHeaderValReader)
	if err != nil {
		f.Errorf("unexpected error while generating fuzz corpus")
	}

	fuzzCorpus := [][]byte{validHeaderVal}

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

func generateRandomMockProtocolBuffer(t testing.TB) io.Reader {
	var w bytes.Buffer
	err := binary.Write(&w, binary.BigEndian, randUint16())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol buffer")
	}
	err = binary.Write(&w, binary.BigEndian, rand.Uint32())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
	}
	err = binary.Write(&w, binary.BigEndian, randUint16())
	if err != nil {
		t.Fatalf("something went wrong while generating mock protocol request")
	}
	return &w
}

// TODO: verify the best way in order to benchmark unmarhsalling behaivour
func Benchmark_Unmarshal(b *testing.B) {
	b.Run("pre generated random message", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		val := generateRandomMockProtocolBuffer(b)
		for n := 0; n < b.N; n++ {
			_, _ = Unmarshal(val)
		}
	})
}
