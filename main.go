// Package main 实现了区块链网络的示例应用
package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/felixkuang/titanchain/core"
	"github.com/sirupsen/logrus"

	"github.com/felixkuang/titanchain/network"
)

func main() {
	// 创建两个本地传输节点，分别命名为 LOCAL 和 REMOTE
	trLocal := network.NewLocalTransport(network.LocalTransportOpts{Addr: "LOCAL"})
	trRemote := network.NewLocalTransport(network.LocalTransportOpts{Addr: "REMOTE"})

	// 建立双向连接
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	// 启动一个 goroutine，定时向 LOCAL 节点发送交易
	go func() {
		for {
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 创建并配置服务器，包含两个传输节点
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal, trRemote},
	}

	s := network.NewServer(opts)
	s.Start()
}

// sendTransaction 构造并发送一笔交易到指定节点
// tr: 发送方传输层
// to: 目标节点地址
// 返回发送过程中的错误
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	//privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000000000)), 10))
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
