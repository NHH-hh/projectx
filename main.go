package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"math/rand"
	"projectx/core"
	"projectx/crypto"
	"projectx/network"
	"strconv"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")
	_ = trLocal.Connect(trRemote)
	_ = trRemote.Connect(trLocal)
	go func() {
		for {
			//_ = trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(time.Second)
		}
	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}
	server := network.NewServer(opts)
	server.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(10000000000)), 10))
	tx := core.NewTransaction(data)
	if err := tx.Sign(privKey); err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGolTxEncoder(buf)); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())
}
