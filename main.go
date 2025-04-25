// Package main 实现了区块链网络的示例应用
package main

import (
	"time"

	"github.com/felixkuang/titanchain/network"
)

func main() {
	trLocal := network.NewLocalTransport(network.LocalTransportOpts{Addr: "LOCAL"})
	trRemote := network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTE"})

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	// 创建并配置服务器
	opts := network.ServerOptions{
		Transports: []network.Transport{trLocal, trRemote},
	}

	s := network.NewServer(opts)
	s.Start()
}
