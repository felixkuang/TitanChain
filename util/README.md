# util 工具包说明

该包为 TitanChain 区块链项目提供常用的工具函数，主要用于测试、数据生成和断言。

## 文件说明

### random.go
- 提供生成随机字节、随机哈希、随机交易的工具函数。
- 主要接口：
  - `RandomBytes(size int) []byte`：生成指定长度的随机字节切片
  - `RandomHash() types.Hash`：生成随机 32 字节哈希
  - `NewRandomTransaction(size int) *core.Transaction`：生成未签名的随机交易
  - `NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *core.Transaction`：生成带签名的随机交易
- 典型用途：
  - 单元测试中快速生成随机数据、交易、哈希
  - 模拟区块链环境下的多样化输入

### assert.go
- 提供简易断言工具，便于测试和调试。
- 主要接口：
  - `AssertEqual(a, b any)`：判断两个对象是否深度相等，不相等时终止程序并输出详细信息
- 典型用途：
  - 测试用例中的断言判断
  - 调试阶段的快速验证

## 主要功能
- 高效生成随机数据、哈希、交易，便于测试覆盖
- 简单易用的断言工具，提升测试代码可读性
- 与 core、types、crypto 等模块无缝集成

## 典型用例
- 区块链核心模块的单元测试、集成测试
- 随机化输入下的健壮性验证
- 断言区块、交易、状态等对象的等价性

## 后续开发计划
1. 增加更多随机数据生成器（如随机区块、随机地址等）
2. 丰富断言工具（如 AssertNotNil、AssertTrue 等）
3. 支持性能测试相关工具函数
4. 提供常用数据转换、格式化等通用工具
5. 完善文档和用例示范
