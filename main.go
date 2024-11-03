package main

import (
	"projectx/network"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")
	_ = trLocal.Connect(trRemote)
	_ = trRemote.Connect(trLocal)
	go func() {
		for {
			_ = trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			time.Sleep(time.Second)
		}
	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}
	server := network.NewServer(opts)
	server.Start()
}
