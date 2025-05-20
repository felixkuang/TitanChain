// Package network 包实现了区块链的网络传输层
package network

// NetAddr 表示网络地址，使用字符串类型
type NetAddr string

// Transport 定义了节点间网络通信的接口
type Transport interface {
	// Consume 返回一个用于接收传入RPC消息的通道
	Consume() <-chan RPC

	// Connect 与另一个传输节点建立连接
	// 如果连接失败则返回错误
	Connect(Transport) error

	// SendMessage 向指定的网络地址发送消息
	// 参数：
	//   - to: 目标地址
	//   - payload: 消息内容
	// 如果发送失败则返回错误
	SendMessage(to NetAddr, payload []byte) error

	// Addr 返回当前节点的网络地址
	Addr() NetAddr
}
