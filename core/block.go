// Package core 实现了TitanChain区块链的核心功能
package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/felixkuang/titanchain/types"
)

// Header 表示区块链中区块的元数据结构
type Header struct {
	Version   uint32     // 协议版本号
	PrevBlock types.Hash // 前一个区块的哈希值
	Timestamp int64      // 区块创建时间戳
	Height    uint32     // 区块在链中的高度
	Nonce     uint64     // 用于挖矿和共识的随机数
}

// EncodeBinary 将区块头序列化为二进制格式
// 使用小端序（LittleEndian）写入所有字段
func (h *Header) EncodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, &h.Nonce)
}

// DecodeBinary 从二进制格式反序列化区块头
// 使用小端序（LittleEndian）读取所有字段
func (h *Header) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Read(r, binary.LittleEndian, &h.Nonce)
}

// Block 表示区块链中的完整区块
type Block struct {
	Header                     // 嵌入区块头
	Transactions []Transaction // 区块中包含的交易列表

	hash types.Hash // 缓存的区块头哈希值，用于提升性能
}

// Hash 计算并返回区块头的哈希值
// 使用SHA-256进行哈希计算，并缓存结果
func (b *Block) Hash() types.Hash {
	buf := &bytes.Buffer{}
	b.Header.EncodeBinary(buf)

	if b.hash.IsZero() {
		b.hash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.hash
}

// EncodeBinary 将整个区块（区块头和交易）序列化为二进制格式
func (b *Block) EncodeBinary(w io.Writer) error {
	if err := b.Header.EncodeBinary(w); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		if err := tx.EncodeBinary(w); err != nil {
			return err
		}
	}

	return nil
}

// DecodeBinary 从二进制格式反序列化整个区块（区块头和交易）
func (b *Block) DecodeBinary(r io.Reader) error {
	if err := b.Header.DecodeBinary(r); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		if err := tx.DecodeBinary(r); err != nil {
			return err
		}
	}

	return nil
}
