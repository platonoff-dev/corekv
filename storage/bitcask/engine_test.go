package bitcask

import (
	"testing"

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

func TestDelete(t *testing.T) {
	engine, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	key := []byte("exampleKey")
	value := []byte("exampleValue")

	err = engine.Put(key, value)
	require.NoError(t, err)

	v, err := engine.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, v)

	err = engine.Delete(key)
	require.NoError(t, err)

	v, err = engine.Get(key)
	require.NoError(t, err)
	assert.Nil(t, v)
}
