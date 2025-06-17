// Package main 实现了区块链网络的示例应用
package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/felixkuang/titanchain/core"
	"github.com/felixkuang/titanchain/crypto"
	"github.com/sirupsen/logrus"

	"github.com/felixkuang/titanchain/network"
)

func main() {
	// 创建两个本地传输节点，分别命名为 LOCAL 和 REMOTE
	trLocal := network.NewLocalTransport(network.LocalTransportOpts{Addr: "LOCAL"})
	trRemoteA := network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTEA"})
	trRemoteB := network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTEB"})
	trRemoteC := network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTEC"})

	// 建立双向连接
	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)
	trRemoteA.Connect(trLocal)

	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	// 启动一个 goroutine，定时向 LOCAL 节点发送交易
	go func() {
		for {
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)

		trLate := network.NewLocalTransport(network.LocalTransportOpts{Addr: "LATE_REMOTE"})
		trRemoteC.Connect(trLate)
		lateServer := makeServer(string(trLate.Addr()), trLate, nil)

		go lateServer.Start()
	}()

	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &privKey)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{tr},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

// sendTransaction 构造并发送一笔交易到指定节点
// tr: 发送方传输层
// to: 目标节点地址
// 返回发送过程中的错误
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	//privKey := crypto.GeneratePrivateKey()
	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	tx := core.NewTransaction(data)
	//tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	// 构造网络消息并发送
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
