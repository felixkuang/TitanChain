// Package core 实现了TitanChain区块链的核心功能
package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/felixkuang/titanchain/crypto"
	"github.com/felixkuang/titanchain/types"
)

// Header 表示区块链中区块的元数据结构
type Header struct {
	Version       uint32     // 协议版本号
	DataHash      types.Hash // 区块中所有交易的Merkle根哈希
	PrevBlockHash types.Hash // 前一个区块的哈希值
	Height        uint32     // 区块在链中的高度
	Timestamp     int64      // 区块创建时间戳
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(h)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

// Block 表示区块链中的完整区块
type Block struct {
	*Header                        // 嵌入区块头
	Transactions []Transaction     // 区块中包含的交易列表
	Validator    crypto.PublicKey  // 验证者的公钥
	Signature    *crypto.Signature // 验证者对区块的签名
	hash         types.Hash        // 缓存的区块头哈希值，用于提升性能
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

// NewBlock 创建一个新的区块实例
// h: 区块头信息
// txx: 要包含在区块中的交易列表
func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

// Sign 使用给定的私钥对区块进行签名
// privKey: 用于签名的私钥
// 返回可能发生的错误
func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

// Verify 验证区块的签名是否有效
// 返回验证过程中可能发生的错误
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	return nil
}

// Decode 从输入流中解码区块数据
// r: 输入流
// dec: 解码器实例
// 返回解码过程中可能发生的错误
func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

// Encode 将区块数据编码到输出流
// w: 输出流
// enc: 编码器实例
// 返回编码过程中可能发生的错误
func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

// Hash 计算并返回区块的哈希值
// hasher: 用于计算哈希的实例
// 如果哈希值已经计算过，则直接返回缓存的值
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}
