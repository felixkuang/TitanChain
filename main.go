// Package main 实现了区块链网络的示例应用
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/felixkuang/titanchain/core"
	"github.com/felixkuang/titanchain/crypto"
	"github.com/felixkuang/titanchain/network"
	"log"
	"time"
)

var transports = []network.Transport{
	network.NewLocalTransport(network.LocalTransportOpts{Addr: "LOCAL"}),
	//network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTE_B"}),
	//network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTE_C"}),
}

func main() {
	initRemoteServers(transports)
	localNode := transports[0]
	trLate := network.NewLocalTransport(network.LocalTransportOpts{Addr: "LATE_NODE"})
	//remoteNodeA := transports[1]
	// remoteNodeC := transports[3]

	//go func() {
	//	for {
	//		if err := sendTransaction(remoteNodeA, localNode.Addr()); err != nil {
	//			logrus.Error(err)
	//		}
	//		time.Sleep(2 * time.Second)
	//	}
	//}()

	go func() {
		time.Sleep(7 * time.Second)
		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
		go lateServer.Start()
	}()

	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", localNode, &privKey)
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
		Transport:  tr,
		PrivateKey: pk,
		ID:         id,
		Transports: transports,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
	var (
		getStatusMsg = new(network.GetStatusMessage)
		buf          = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

// sendTransaction 构造并发送一笔交易到指定节点
// tr: 发送方传输层
// to: 目标节点地址
// 返回发送过程中的错误
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	//data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	// 构造网络消息并发送
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
