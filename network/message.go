package network

type GetBlocksMessage struct {
	From uint32
	// If To is 0 the maximum blocks will be returned.
	To uint32
}

// GetStatusMessage 用于节点间请求状态的网络消息结构体
// 主要用于节点间同步区块高度、ID等信息
// 一般由节点主动发起状态请求时发送
// 无需携带额外字段
type GetStatusMessage struct{}

// StatusMessage 用于节点间返回状态的网络消息结构体
// 包含节点ID、版本号、当前区块高度等信息
// 用于节点间状态同步和健康检查
//
// 字段说明：
//
//	ID: 节点唯一标识
//	Version: 节点软件版本号
//	CurrentHeight: 当前区块高度
type StatusMessage struct {
	// 节点唯一标识
	ID string
	// 节点软件版本号
	Version uint32
	// 当前区块高度
	CurrentHeight uint32
}
