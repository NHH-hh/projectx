package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	_ = tra.Connect(trb)
	_ = trb.Connect(tra)
	//	assert.Equal(t, tra.peers[trb.Addr()], trb)
	//	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	_ = tra.Connect(trb)
	_ = trb.Connect(tra)
	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}
