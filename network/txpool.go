package network

import (
	"projectx/core"
	"projectx/types"
	"sort"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))
	i := 0
	for _, v := range txMap {
		txx[i] = v
		i++
	}
	s := &TxMapSorter{transactions: txx}
	sort.Sort(s)
	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

type TxPool struct {
	transports map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transports: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transports)
	return s.transactions
}

// Add adds a transaction to the poo, the caller is responsible
// checking if the tx already exist
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	p.transports[hash] = tx
	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transports[hash]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transports)
}

func (p *TxPool) Flush() {
	p.transports = make(map[types.Hash]*core.Transaction)
}
