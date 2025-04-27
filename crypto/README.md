# Crypto 包

这个包实现了 TitanChain 区块链的加密功能，提供了密钥生成、签名和验证等核心加密操作。

## 文件说明

### keypair.go
实现了基于 ECDSA（椭圆曲线数字签名算法）的密钥对功能：
- `PrivateKey`: 私钥结构及其操作
  - 生成新的私钥
  - 签名数据
  - 获取对应的公钥
- `PublicKey`: 公钥结构及其操作
  - 转换为字节序列
  - 生成区块链地址
- `Signature`: 签名结构及其操作
  - 签名验证功能

### keypair_test.go
包含了密钥对功能的单元测试：
- 测试签名和验证的成功场景
- 测试签名验证的失败场景
- 验证密钥对的正确性

## 核心功能

### 密钥管理
- 使用 P256 椭圆曲线
- 安全的随机数生成
- 私钥和公钥的封装
- 密钥格式转换

### 数字签名
- ECDSA 签名生成
- 签名验证
- 防篡改保护

### 地址生成
- 基于公钥生成地址
- 使用 SHA256 哈希
- 兼容区块链地址格式

## 安全特性
- 使用标准密码库
- 安全的随机数生成
- 密钥数据保护
- 签名不可伪造

## 使用示例

```go
// 生成新的密钥对
privateKey := GeneratePrivateKey()
publicKey := privateKey.PublicKey()

// 签名数据
message := []byte("要签名的数据")
signature, err := privateKey.Sign(message)
if err != nil {
    // 处理错误
}

// 验证签名
isValid := signature.Verify(publicKey, message)
```

## 后续开发计划
1. 添加密钥序列化功能
2. 实现密钥导入导出
3. 支持多种签名算法
4. 添加密钥派生功能
5. 实现分层确定性钱包
6. 增强密钥安全存储 