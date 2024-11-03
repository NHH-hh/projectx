package network

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"projectx/core"
	"strconv"
	"testing"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)
	_ = core.NewTransaction([]byte("foo"))
	assert.Equal(t, p.Len(), 1)
	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransaction(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(100000)))
		assert.Nil(t, p.Add(tx))
	}
	assert.Equal(t, p.Len(), txLen)
	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
