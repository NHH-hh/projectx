package network

import (
	"github.com/stretchr/testify/assert"
	"io"
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
	trb := NewLocalTransport("B")
	tra := NewLocalTransport("A")
	_ = tra.Connect(trb)
	_ = trb.Connect(tra)
	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	rpc := <-trb.Consume()
	buf, err := io.ReadAll(rpc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, len(buf), len(msg))
	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")
	_ = tra.Connect(trb)
	_ = tra.Connect(trc)
	msg := []byte("hello world")
	assert.Nil(t, tra.Broadcast(msg))
	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)
	rpcc := <-trc.Consume()
	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, c)

}
