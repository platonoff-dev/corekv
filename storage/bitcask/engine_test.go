package bitcask

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	engine, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	key := []byte("exampleKey")
	value := []byte("exampleValue")

	err = engine.Put(key, value)
	require.NoError(t, err)

	retrievedValue, err := engine.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, retrievedValue)
}

func TestDataEntryEncodeDecode(t *testing.T) {
	tests := []struct {
		name  string
		entry DataEntry
	}{
		{
			name: "basic entry",
			entry: DataEntry{
				Key:   []byte("test-key"),
				Value: []byte("test-value"),
				TS:    time.Now(),
			},
		},
		{
			name: "empty key",
			entry: DataEntry{
				Key:   []byte(""),
				Value: []byte("value"),
				TS:    time.Unix(1234567890, 123456789),
			},
		},
		{
			name: "empty value",
			entry: DataEntry{
				Key:   []byte("key"),
				Value: []byte(""),
				TS:    time.Unix(0, 0),
			},
		},
		{
			name: "large values",
			entry: DataEntry{
				Key:   bytes.Repeat([]byte("k"), 1000),
				Value: bytes.Repeat([]byte("v"), 10000),
				TS:    time.Unix(1234567890, 999999999),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := tt.entry.Encode()
			require.NoError(t, err)
			require.NotEmpty(t, encoded)

			// Decode
			var decoded DataEntry
			reader := bytes.NewReader(encoded)
			err = decoded.Decode(reader)
			require.NoError(t, err)

			// Verify
			assert.Equal(t, tt.entry.Key, decoded.Key)
			assert.Equal(t, tt.entry.Value, decoded.Value)
			assert.Equal(t, tt.entry.TS.UnixNano(), decoded.TS.UnixNano())
		})
	}
}

func TestDataEntryDecodeInvalidChecksum(t *testing.T) {
	entry := DataEntry{
		Key:   []byte("test"),
		Value: []byte("value"),
		TS:    time.Now(),
	}

	encoded, err := entry.Encode()
	require.NoError(t, err)

	// Corrupt the checksum (first 4 bytes)
	encoded[0] ^= 0xFF

	var decoded DataEntry
	reader := bytes.NewReader(encoded)
	err = decoded.Decode(reader)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "checksum validation failed")
}

func TestDataEntryDecodeInvalidLength(t *testing.T) {
	entry := DataEntry{
		Key:   []byte("test"),
		Value: []byte("value"),
		TS:    time.Now(),
	}

	encoded, err := entry.Encode()
	require.NoError(t, err)

	// Truncate the data to make it incomplete
	truncated := encoded[:len(encoded)-2]

	var decoded DataEntry
	reader := bytes.NewReader(truncated)
	err = decoded.Decode(reader)
	assert.Error(t, err)
}
