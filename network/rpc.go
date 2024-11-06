package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"projectx/core"
)

type MessageType byte

const (
	MessageTypeTx    MessageType = 0x1
	MessageTypeBlock MessageType = 0x2
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

type DecodeMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodeMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}
	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")
	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	_ = gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type RPCHandler interface {
	HandleRPC(RPC) error
}

type DefaultRPCHandler struct {
	p RPCProcessor
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		p: p,
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodeMessage) error
}
