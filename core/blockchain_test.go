package core

import (
	"github.com/stretchr/testify/assert"
	"projectx/types"
	"testing"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(t, 0, types.RandomHash()))
	assert.Nil(t, err)
	return bc
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHashBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(uint32(0)))
	assert.False(t, bc.HasBlock(uint32(1)))
	assert.False(t, bc.HasBlock(uint32(100)))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlock := 1000
	for i := 0; i < lenBlock; i++ {
		b := randomBlock(t, uint32(i+1), getPrevBlocHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))
	}
	assert.Equal(t, uint32(lenBlock), bc.Height())
	assert.Equal(t, lenBlock+1, len(bc.headers))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, uint32(89), types.RandomHash())))
}

func TestAddBlockToHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.Nil(t, bc.AddBlock(randomBlock(t, uint32(1), getPrevBlocHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, uint32(3), types.RandomHash())))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlock := 1000
	for i := 0; i < lenBlock; i++ {
		b := randomBlock(t, uint32(i+1), getPrevBlocHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))
		header, err := bc.GetHeader(b.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, b.Header)
	}
}

func getPrevBlocHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
