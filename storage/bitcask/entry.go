package bitcask

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"time"
)

var endian = binary.LittleEndian

type DataEntry struct {
	Key   []byte
	Value []byte
	TS    time.Time
}

func (e *DataEntry) Encode() ([]byte, error) {
	data := []any{
		e.TS.UnixNano(),
		int64(len(e.Key)),
		int64(len(e.Value)),
	}

	var header []byte
	var err error
	for _, v := range data {
		header, err = binary.Append(header, endian, v)
		if err != nil {
			return nil, err
		}
	}

	entry := bytes.Join([][]byte{header, e.Key, e.Value}, []byte{})

	checksum := crc32.Checksum(entry, crc32.MakeTable(crc32.Castagnoli))
	checksumBytes, err := binary.Append(nil, endian, checksum)
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{checksumBytes, entry}, []byte{}), nil
}

func (e *DataEntry) Decode(reader io.Reader) error {
	// Read checksum (4 bytes)
	var checksum uint32
	if err := binary.Read(reader, endian, &checksum); err != nil {
		return err
	}

	// Read timestamp (8 bytes)
	var tsNano int64
	if err := binary.Read(reader, endian, &tsNano); err != nil {
		return err
	}

	// Read key length (8 bytes)
	var keyLen int64
	if err := binary.Read(reader, endian, &keyLen); err != nil {
		return err
	}

	// Read value length (8 bytes)
	var valueLen int64
	if err := binary.Read(reader, endian, &valueLen); err != nil {
		return err
	}

	// Validate lengths
	if keyLen < 0 || valueLen < 0 {
		return errors.New("invalid key or value length")
	}

	// Read key
	key := make([]byte, keyLen)
	if _, err := io.ReadFull(reader, key); err != nil {
		return err
	}

	// Read value
	value := make([]byte, valueLen)
	if _, err := io.ReadFull(reader, value); err != nil {
		return err
	}

	// Reconstruct entry for checksum validation
	data := []any{tsNano, keyLen, valueLen}
	var header []byte
	for _, v := range data {
		var err error
		header, err = binary.Append(header, endian, v)
		if err != nil {
			return err
		}
	}
	entry := bytes.Join([][]byte{header, key, value}, []byte{})

	// Validate checksum
	expectedChecksum := crc32.Checksum(entry, crc32.MakeTable(crc32.Castagnoli))
	if checksum != expectedChecksum {
		return errors.New("checksum validation failed")
	}

	// Set decoded values
	e.Key = key
	e.Value = value
	e.TS = time.Unix(0, tsNano)

	return nil
}
