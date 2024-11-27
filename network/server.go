package network

import (
	"bytes"
	"github.com/go-kit/log"
	"os"
	"projectx/core"
	"projectx/crypto"
	"projectx/types"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == 0 {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}
	// if server don't got any processor from the server options
	// we're going to use the server as default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	if s.isValidator {
		go s.validatorLoop()
	}
	return s, nil
}

func (s *Server) Start() {
	s.initTransports()
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				_ = s.Logger.Log("error", err)
			}
			if err = s.RPCProcessor.ProcessMessage(msg); err != nil {
				_ = s.Logger.Log("error", err)
			}
		case <-s.quitCh:
			break free
		}
	}
	_ = s.Logger.Log("Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	_ = s.Logger.Log("msg", "Starting validator loop", "block_time", s.BlockTime)
	for {
		<-ticker.C
		_ = s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(message *DecodeMessage) error {
	switch t := message.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}
	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		return nil
	}
	if err := tx.Verify(); err != nil {
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	_ = s.Logger.Log("msg", "adding new tx to memPool",
		"hash", hash, "memPoolLength", s.memPool.Len())
	go s.broadcastTx(tx)
	return s.memPool.Add(tx)
}

func (s *Server) broadcastBlock(b *core.Block) error {
	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGolTxEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}
	// For now, we are going to use all transactions that are in the mempool
	// Later on when we know thr internal structure of our transaction
	// we will implement same kind of complexity function to determine how
	// many transactions can be included in a block.
	txx := s.memPool.Transactions()
	block, err := core.NewBlockFromPreHeader(currentHeader, txx)
	if err != nil {
		return err
	}
	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}
	if err := s.chain.AddBlock(block); err != nil {
		return err
	}
	s.memPool.Flush()
	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:  1,
		DataHash: types.Hash{},
		Height:   0,
		//Timestamp: time.Now().UnixNano(),
		Timestamp: 000000,
	}
	b, _ := core.NewBlock(header, nil)
	return b
}
