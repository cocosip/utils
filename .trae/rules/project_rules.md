# Trae 项目规则 for go-utils

本文档定义了 `go-utils` 工具库的开发规范、代码标准和最佳实践。遵循这些规则确保代码质量、一致性和可维护性。

## 1. 项目概述

`go-utils` 是一个综合性的 Go 工具库，提供以下功能模块：

- **crypto**: 加密解密工具（AES、MD5、SHA、SM2/SM3 国密算法）
- **database**: 数据库操作工具（支持 MySQL、PostgreSQL、SQLite、SQL Server）
- **util**: 通用工具函数（字符串、切片、文件、路径、URL 处理）
- **httpx**: HTTP 客户端工具
- **time**: 时间处理工具
- **snowflake**: 分布式 ID 生成器
- **slogx**: 结构化日志工具
- **encoding**: 编码解码工具
- **buildinfo**: 构建信息工具
- **context**: 上下文工具
- **daemon**: 守护进程工具
- **ptrutil**: 指针工具包

## 2. 技术规范

### 2.1 基础配置

- **Go 版本**: 1.25+
- **模块路径**: `github.com/cocosip/utils`
- **依赖管理**: Go Modules
- **测试框架**: `github.com/stretchr/testify`

### 2.2 主要依赖

- **GORM**: 数据库 ORM 框架
- **国密算法**: `github.com/tjfoc/gmsm`
- **日志轮转**: `gopkg.in/natefinch/lumberjack.v2`
- **系统服务**: `github.com/kardianos/service`

## 3. 代码规范

### 3.1 包结构规范

- 每个功能模块独立成包
- 包名使用小写字母，简洁明了
- 每个包必须包含对应的测试文件（`*_test.go`）
- 公共常量和变量使用 `PascalCase`
- 私有变量使用 `camelCase`

### 3.2 命名规范

- **函数命名**: 公共函数使用 `PascalCase`，私有函数使用 `camelCase`
- **结构体命名**: 使用 `PascalCase`
- **接口命名**: 使用 `PascalCase`，通常以 `-er` 结尾
- **常量命名**: 使用 `PascalCase` 或 `UPPER_SNAKE_CASE`
- **错误变量**: 以 `Err` 开头，如 `ErrNotSupport`

### 3.3 注释规范

- **包注释**: 每个包必须有包级注释，描述包的功能
- **函数注释**: 所有公共函数必须有注释，格式为 `// FunctionName description`
- **结构体注释**: 公共结构体必须有注释
- **注释语言**: 使用英文编写注释
- **参数说明**: 复杂函数需要说明参数和返回值

### 3.4 错误处理规范

- 使用 `error` 接口返回错误
- 错误信息使用英文
- 提供带错误返回的函数版本（如 `ParseIntE`）和默认值版本（如 `ParseIntOrDefault`）
- 自定义错误使用 `fmt.Errorf` 或实现 `error` 接口

### 3.5 函数设计规范

- **构造函数**: 使用 `New` 前缀，如 `NewAES()`
- **链式调用**: 支持方法链，返回结构体指针，如 `WithKey()`, `WithIV()`
- **选项模式**: 复杂配置使用选项模式或链式调用
- **默认值**: 提供合理的默认配置

## 4. 开发工作流程

### 4.1 代码质量检查

在提交代码前，必须执行以下检查：

```bash
# 代码格式化
gofmt -w .

# 静态分析
go vet ./...

# 依赖整理
go mod tidy

# 运行测试
go test -v ./...

# 测试覆盖率
go test -cover ./...
```

### 4.2 测试要求

- **单元测试**: 每个公共函数必须有对应的单元测试
- **测试文件**: 测试文件命名为 `*_test.go`
- **测试函数**: 测试函数以 `Test` 开头
- **基准测试**: 性能敏感的函数需要基准测试（`Benchmark*`）
- **示例测试**: 复杂功能提供示例测试（`Example*`）
- **测试覆盖率**: 目标覆盖率 > 80%

### 4.3 跨平台兼容性

- 代码必须支持 Windows、Linux、macOS
- 文件路径使用 `filepath` 包处理
- 避免使用平台特定的系统调用
- 使用构建标签处理平台差异

## 5. 模块特定规范

### 5.1 crypto 模块

- 提供默认密钥和 IV 常量
- 支持多种加密模式（CBC、ECB、CFB、OFB）
- 支持多种填充模式（PKCS7、Zero、ANSIX923）
- 提供 Base64 和 Hex 编码输出

### 5.2 database 模块

- 支持主流数据库（MySQL、PostgreSQL、SQLite、SQL Server）
- 使用 GORM 作为 ORM 框架
- 提供数据库连接池配置
- 支持数据库连接的优雅关闭

### 5.3 util 模块

- 提供类型安全的转换函数
- 同时提供带错误返回和默认值版本
- 字符串处理支持常见操作
- 切片操作保证类型安全

## 6. 版本管理

### 6.1 语义化版本

遵循 [Semantic Versioning](https://semver.org/) 规范：

- **MAJOR**: 不兼容的 API 变更
- **MINOR**: 向后兼容的功能新增
- **PATCH**: 向后兼容的问题修复

### 6.2 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**类型 (type)**:
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档更新
- `style`: 代码格式化
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**范围 (scope)**: 模块名称，如 `crypto`, `database`, `util` 等

**示例**:
```
feat(crypto): add AES-GCM encryption support
fix(database): resolve connection timeout issue
docs(util): update string utility function examples
test(snowflake): add benchmark tests for ID generation
```

## 7. 性能要求

- 避免不必要的内存分配
- 使用对象池减少 GC 压力
- 提供基准测试验证性能
- 关键路径函数需要性能优化
- 支持并发安全的操作

## 8. 安全要求

- 加密模块使用安全的随机数生成器
- 避免硬编码敏感信息
- 输入验证和边界检查
- 防止常见的安全漏洞（如注入攻击）
- 定期更新依赖库以修复安全漏洞