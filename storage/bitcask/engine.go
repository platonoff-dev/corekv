package bitcask

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type pointer struct {
	fileID        int
	valuePosition int64
	valueSize     int
	timestamp     time.Time
}

type fileStat struct {
	size int64
}

type Engine struct {
	dirPath string
	keyDir  map[string]pointer
	files   []*os.File
	stats   []*fileStat
}

func New(path string) (*Engine, error) {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return nil, err
	}

	dataFilePath := filepath.Join(path, "active.db")
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
		dirPath: path,
		keyDir:  make(map[string]pointer),
		files:   []*os.File{file},
		stats: []*fileStat{
			{
				size: stat.Size(),
			},
		},
	}, nil
}

func (e *Engine) Put(key []byte, value []byte) error {
	ts := time.Now()
	entry := DataEntry{
		Key:   key,
		Value: value,
		TS:    ts,
	}
	data, err := entry.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	headerSize := len(data) - len(value)
	pos := e.stats[0].size + int64(headerSize)

	_, err = e.files[0].Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	e.stats[0].size += int64(len(data))
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

func (e *Engine) Delete(_ []byte) error {
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
