package bitcask

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"time"
)

type pointer struct {
	fileID        uint16
	valuePosition int64
	valueSize     int
	timestamp     time.Time
}

type Engine struct {
	keyDir   map[string]pointer
	files    []*os.File
	fileSize int64
}

func New(path string) (*Engine, error) {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return nil, err
	}

	dataFilePath := filepath.Join(path, "data.db")
	file, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat data file: %w", err)
	}

	// TODO: Fill keyDir if hint/data files exist
	return &Engine{
		keyDir:   make(map[string]pointer),
		files:    []*os.File{file},
		fileSize: stat.Size(),
	}, nil
}

func (e *Engine) Put(key []byte, value []byte) error {
	ts := time.Now()
	data := encode(key, value, ts)

	headerSize := len(data) - len(value)
	pos := e.fileSize + int64(headerSize)

	_, err := e.files[0].Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	e.fileSize += int64(len(data))
	e.keyDir[string(key)] = pointer{
		fileID:        0,
		valuePosition: pos,
		valueSize:     len(value),
		timestamp:     ts,
	}

	return nil
}

func (e *Engine) Get(key []byte) ([]byte, error) {
	p, ok := e.keyDir[string(key)]
	if !ok {
		return nil, nil
	}

	data := make([]byte, p.valueSize)
	_, err := e.files[p.fileID].ReadAt(data, p.valuePosition)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *Engine) Delete(key []byte) error {
	data := tombstoneRecord(key, time.Now())

	_, err := e.files[0].Write(data)
	if err != nil {
		return err
	}

	delete(e.keyDir, string(key))

	return nil
}

func (e *Engine) Close() error {
	for _, file := range e.files {
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func encode(k, v []byte, ts time.Time) []byte {
	var header []byte

	header = binary.AppendVarint(header, ts.UnixNano())
	header = binary.AppendVarint(header, int64(len(k)))
	header = binary.AppendVarint(header, int64(len(v)))

	entry := bytes.Join([][]byte{header, k, v}, []byte{})

	checksum := crc32.Checksum(entry, crc32.MakeTable(crc32.Castagnoli))
	checksumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(checksumBytes, checksum)

	return bytes.Join([][]byte{checksumBytes, entry}, []byte{})
}

func tombstoneRecord(k []byte, ts time.Time) []byte {
	var header []byte

	header = binary.AppendVarint(header, ts.UnixNano())
	header = binary.AppendVarint(header, int64(len(k)))
	header = binary.AppendVarint(header, 0) // Value size is 0 for tombstone

	entry := bytes.Join([][]byte{header, k}, []byte{})

	checksum := crc32.Checksum(entry, crc32.MakeTable(crc32.Castagnoli))
	checksumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(checksumBytes, checksum)

	return bytes.Join([][]byte{checksumBytes, entry}, []byte{})
}
