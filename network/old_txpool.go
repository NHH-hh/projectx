package network

import (
	"projectx/core"
	"projectx/types"
	"sort"
)

type OldTxMapSorter struct {
	transactions []*core.Transaction
}

func NewOldTxMapSorter(txMap map[types.Hash]*core.Transaction) *OldTxMapSorter {
	txx := make([]*core.Transaction, len(txMap))
	i := 0
	for _, v := range txMap {
		txx[i] = v
		i++
	}
	s := &OldTxMapSorter{transactions: txx}
	sort.Sort(s)
	return s
}

func (s *OldTxMapSorter) Len() int {
	return len(s.transactions)
}

func (s *OldTxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

func (s *OldTxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

type OldTxPool struct {
	transports map[types.Hash]*core.Transaction
}

func NewOldTxPool() *OldTxPool {
	return &OldTxPool{
		transports: make(map[types.Hash]*core.Transaction),
	}
}

func (p *OldTxPool) Transactions() []*core.Transaction {
	s := NewOldTxMapSorter(p.transports)
	return s.transactions
}

// Add adds a transaction to the poo, the caller is responsible
// checking if the tx already exist
func (p *OldTxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	p.transports[hash] = tx
	return nil
}

func (p *OldTxPool) Has(hash types.Hash) bool {
	_, ok := p.transports[hash]
	return ok
}

func (p *OldTxPool) Len() int {
	return len(p.transports)
}

func (p *OldTxPool) Flush() {
	p.transports = make(map[types.Hash]*core.Transaction)
}
