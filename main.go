package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"projectx/core"
	"projectx/crypto"
	"projectx/network"
	"strconv"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")
	_ = trLocal.Connect(trRemoteA)
	_ = trRemoteA.Connect(trRemoteB)
	_ = trRemoteB.Connect(trRemoteC)
	_ = trRemoteA.Connect(trLocal)
	initRemoteServer([]network.Transport{trRemoteA, trRemoteB, trRemoteC})
	go func() {
		for {
			//		_ = trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &privKey)
	localServer.Start()
}

func initRemoteServer(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: privKey,
		ID:         id,
		Transports: []network.Transport{tr},
	}
	server, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	return server
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
