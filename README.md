# Golang Utils 工具库

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmzlingyin/utils)](https://goreportcard.com/report/github.com/dmzlingyin/utils)

一个功能丰富的 Go 语言工具库，提供企业级应用开发所需的各种常用功能模块，助力快速构建高质量的 Go 应用程序。

## 📋 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [核心模块](#核心模块)
- [使用示例](#使用示例)
- [配置说明](#配置说明)
- [测试](#测试)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 🚀 功能特性

### 🏗️ 核心功能
- **IOC容器**：轻量级依赖注入框架，简化服务管理
- **懒加载**：高效的资源延迟初始化机制
- **配置管理**：灵活的配置文件读取和管理
- **缓存系统**：支持内存和Redis的统一缓存接口
- **日志系统**：可配置的日志记录和输出
- **类型转换**：强大的类型转换工具，支持结构体与Map互转

### 🌐 云服务集成
- **阿里云服务**：OSS对象存储、短信服务、内容审核、机器翻译
- **腾讯云服务**：COS对象存储、短信服务、邮件服务
- **多平台OAuth2**：Google、微信、Apple、Facebook、Discord、Twitter、Casdoor

### 💳 支付系统
- **移动支付**：Apple Pay、Google Pay（支持沙盒环境）
- **第三方支付**：PayPal、Stripe
- **国内支付**：微信支付、支付宝
- **平台支付**：抖音支付、快手支付

### 🛠️ 实用工具
- **网络请求**：HTTP客户端封装，内置退避重试机制
- **URL验证**：完整的URL格式验证
- **数学运算**：常用数学函数
- **图像处理**：基础图像操作
- **加密解密**：常用加密算法
- **随机生成**：各种随机数据生成
- **短链生成**：URL短链服务
- **退避重试**：智能重试机制，支持线性退避

### 🌐 Web开发
- **Gin集成**：路由自动注册和管理
- **CORS支持**：跨域资源共享配置
- **服务器封装**：HTTP服务器快速启动

### 🗄️ 数据库支持
- **MongoDB集成**：完整的MongoDB客户端封装，支持事务处理

## 🚀 快速开始

### 安装

```bash
go get github.com/dmzlingyin/utils
```

### Hello World 示例

```go
package main

import (
    "context"
    "github.com/dmzlingyin/utils/config"
    "github.com/dmzlingyin/utils/ioc"
    "github.com/dmzlingyin/utils/log"
    "github.com/dmzlingyin/utils/misc"
)

func main() {
    // 加载配置
    config.Load("config.json")

    // 初始化日志
    log.Init()

    // 验证URL
    if misc.IsValidURL("https://example.com") {
        log.Info("URL is valid")
    }

    // 使用IOC容器
    container := ioc.NewContainer()
    container.Register("service", func() interface{} {
        return "my service"
    })

    service := container.Get("service")
    log.Infof("Service: %v", service)
}
```

## 📦 核心模块

### IOC容器 (`ioc`)
轻量级的依赖注入容器，支持服务注册、依赖解析和生命周期管理。

### 缓存系统 (`cache`)
统一的缓存接口，支持内存缓存和Redis缓存。

### 配置管理 (`config`)
支持JSON配置文件的读取和全局配置管理。

### 日志系统 (`log`)
可配置的日志记录器，支持不同级别的日志输出。

### 类型转换 (`cast`)
提供强大的类型转换功能，包括基本类型转换和结构体与Map的相互转换。

### 懒加载 (`lazy`)
高效的资源延迟初始化机制，支持线程安全的单次初始化。

### 退避重试 (`backoff`)
智能重试机制，支持线性退避策略，适用于网络请求等不稳定操作。

### 网络请求 (`net`)
HTTP客户端封装，内置退避重试机制，支持JSON自动解析。

### OAuth2认证 (`oauth2`)
支持多平台的OAuth2登录认证，包含JWT令牌解析功能。

### 支付系统 (`payment`)
统一的支付接口，支持多种支付平台，包括Apple Pay和Google Pay的沙盒环境支持。

### 云服务集成

#### 阿里云服务 (`aliyun`)
集成阿里云OSS对象存储、短信服务、内容审核、机器翻译等功能。

#### 腾讯云服务 (`tencent`)
集成腾讯云COS对象存储、短信服务、邮件服务等功能。

### 数据库支持 (`database`)

#### MongoDB (`database/mongo`)
提供MongoDB客户端封装，支持连接池配置和事务处理。

### 实用工具 (`misc`)
提供随机字符串生成、图像处理、短链生成、数学运算、URL验证等实用功能。

## ⚙️ 配置说明

### 项目配置文件
项目使用 JSON 格式的配置文件，支持应用设置、OAuth2 配置、支付配置、短信服务配置、腾讯云服务配置、MongoDB配置等。配置文件采用分层结构，便于管理不同服务的配置项。

### 配置文件示例

```json
{
  "app": {
    "name": "my-app",
    "version": "1.0.0"
  },
  "pay": {
    "apple": {
      "key_path": "./apple_key.p8",
      "key_id": "your-key-id",
      "bundle_id": "com.your.app",
      "issuer": "your-issuer",
      "sandbox": true
    },
    "google": {
      "key_path": "./google_key.json"
    }
  },
  "mongo": {
    "uri": "mongodb://localhost:27017",
    "database": "mydb",
    "max_pool_size": 10,
    "min_pool_size": 2,
    "tx_enabled": false
  }
}
```

## 🧪 测试

项目包含完整的单元测试，支持运行所有测试、特定模块测试以及测试覆盖率检查。

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./payment/...

# 检查测试覆盖率
go test -cover ./...
```

## 📝 使用示例

工具库支持快速构建完整的Web服务，包括配置加载、IOC容器初始化、服务注册、路由配置和服务器启动等功能。

### MongoDB 使用示例

```go
// 初始化MongoDB数据库
db := mongo.NewDatabase()
scope := mongo.NewScope(db)

// 使用事务
err := scope.Transact(context.Background(), func(ctx context.Context) error {
    // 在事务中执行数据库操作
    return nil
})
```

### URL验证示例

```go
// 验证URL格式
isValid := misc.IsValidURL("https://example.com")
```

### 结构体与Map转换示例

```go
// 结构体转Map
m, err := cast.StructToMap(myStruct)

// Map转结构体
err := cast.MapToStruct(myMap, &myStruct)
```

### 退避重试示例

```go
// 带重试的函数调用
err := backoff.Retry(3, func() error {
    // 可能失败的操作
    return someUnstableOperation()
})
```

### 网络请求示例

```go
// 发送HTTP请求并自动解析JSON响应
headers := map[string]string{"Content-Type": "application/json"}
response, err := net.Request[MyResponseType]("POST", "https://api.example.com", headers, requestBody)
```

### 懒加载示例

```go
// 创建懒加载值
lazyValue := &lazy.Value[string]{
    New: func() (string, error) {
        // 昂贵的初始化操作
        return "expensive initialization result", nil
    },
}

// 获取值（首次调用时初始化）
value, err := lazyValue.Get()
```

## 🤝 贡献指南

我们欢迎并感谢任何形式的贡献！

### 贡献方式
1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范
- 遵循 Go 语言规范和最佳实践
- 添加适当的测试用例
- 更新相关文档
- 保持代码风格一致
- 使用 `go fmt` 格式化代码

### 问题反馈
如果您发现了bug或有新功能建议，请通过 [Issues](https://github.com/dmzlingyin/utils/issues) 提交。

## 📋 待办事项

- [ ] 优化性能和内存使用
- [ ] 增加更多实用工具函数
- [ ] 完善文档和示例
- [ ] 添加更多云服务集成

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [Go官方网站](https://golang.org/)
- [Gin Web框架](https://gin-gonic.com/)
- [Redis官方文档](https://redis.io/documentation)
- [MongoDB官方文档](https://www.mongodb.com/docs/)
- [阿里云开发者文档](https://help.aliyun.com/)
- [腾讯云开发者文档](https://cloud.tencent.com/document)
- [Go Report Card](https://goreportcard.com/report/github.com/dmzlingyin/utils)

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！