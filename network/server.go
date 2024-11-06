package network

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"projectx/core"
	"projectx/crypto"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == 0 {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	s := &Server{
		ServerOpts:  opts,
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
	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}
			if err = s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				_ = s.createNewBlock()
			}
		}
	}
	fmt.Println("Server shutdown")
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
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in memPool")
		return nil
	}
	if err := tx.Verify(); err != nil {
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	logrus.WithFields(logrus.Fields{
		"hash":           tx.Hash(core.TxHasher{}),
		"memPool length": s.memPool.Len(),
	}).Info("adding new tx to memPool")
	go s.broadcastTx(tx)
	return s.memPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGolTxEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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
