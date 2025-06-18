// Package crypto 实现了区块链所需的加密功能
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/felixkuang/titanchain/types"
)

var (
	ErrInvalidPublicKey = errors.New("invalid public key length")
)

const (
	PublicKeyLength = 33
)

// PrivateKey 封装了ECDSA私钥
type PrivateKey struct {
	key *ecdsa.PrivateKey // 底层ECDSA私钥
}

// Sign 使用私钥对数据进行签名
// 参数：
//   - data: 要签名的数据
//
// 返回：
//   - 签名对象和可能的错误
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

// GeneratePrivateKey 生成一个新的私钥
// 使用P256曲线和安全的随机数生成器
func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

// PublicKey 从私钥中获取对应的公钥
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

// PublicKey 封装了ECDSA公钥
type PublicKey struct {
	Key *ecdsa.PublicKey // 底层ECDSA公钥
}

// ToSlice 将公钥转换为字节切片（压缩格式）
// 使用elliptic.MarshalCompressed进行编码，便于网络传输和存储
func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(elliptic.P256(), k.Key.X, k.Key.Y)
}

// ToPublicKey 从压缩字节切片还原公钥
// 参数：
//   - b: 压缩格式的公钥字节切片
//
// 返回：
//   - 还原后的PublicKey对象和错误信息
func ToPublicKey(b []byte) (PublicKey, error) {
	if len(b) != PublicKeyLength {
		return PublicKey{}, ErrInvalidPublicKey
	}
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), b)
	if x == nil || y == nil {
		return PublicKey{}, ErrInvalidPublicKey
	}
	return PublicKey{Key: &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}}, nil
}

// Address 从公钥生成地址
// 使用公钥的SHA256哈希的最后20字节作为地址
func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

// Signature 表示ECDSA签名
type Signature struct {
	S, R *big.Int // 签名的两个组成部分
}

// Verify 验证签名是否有效
// 参数：
//   - pubKey: 用于验证的公钥
//   - data: 原始数据
//
// 返回：
//   - 签名是否有效
func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}
