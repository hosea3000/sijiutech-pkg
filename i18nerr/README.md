# i18nerror

国际化错误码与翻译服务包，提供：

- 带业务码的 i18n 错误类型（`ierror` 子包）
- JSON 自动翻译字段类型（`i18ntype` 子包）
- Redis + MySQL 两级缓存的翻译服务
- Gin 中间件自动注入语言
- 协程级语言上下文管理

根目录包名为 **`i18nerr`**（与导入路径 `.../i18nerror` 不同，属 Go 常见写法）。

---

## 包结构

```
i18nerror/
├── service.go       # I18NService 翻译服务实现
├── middleware.go     # Gin 语言注入中间件
├── ctxutil.go       # 协程级语言存储（sync.Map + goroutine ID）
├── translate.go     # TranslateI18nKey 辅助函数
├── i18ntype/
│   └── i18n_string.go  # I18nString — JSON 序列化时自动翻译的字段类型
├── ierror/
│   ├── err.go       # 错误类型：I18nCode、I18nTranslationError、NewError
│   └── errors.go    # 预定义业务错误码（116 个）
├── helper/
│   └── string.go    # SQL Null 类型辅助
└── sqlc/            # sqlc 生成的数据访问层
    ├── sqlc.yaml
    ├── schema/      # 4 张 blade_i18n_* 表 DDL
    ├── sql/         # SQL 查询定义
    └── repository/  # 生成的 Go 代码
```

---

## 在项目中的三种使用模式

### 模式一：业务错误码翻译（主要用法）

在 service 层用 `ierror.NewError()` 创建错误，handler 层通过 `v1.HandleError()` 统一翻译后返回：

```go
// service 层 — 抛出带 i18n key 的错误
return nil, ierror.NewError(ierror.PlanNotFound)

// api/v1/v1.go HandleError — 自动检测 I18nTranslationError，翻译 message
var i18nErr ierror.I18nTranslationError
if errors.As(err, &i18nErr) {
    translatedValue, _ := service.Translate(ctx, lang, i18nErr.Error())
    resp := Response{Code: i18nErr.Code, Message: translatedValue, ...}
    ctx.JSON(http.StatusOK, resp)
}
```

涉及文件：`subscription.go`、`material_order.go`、`provider_stripe.go`、`provider_apple.go`、`provider_google.go`、`subscription_intent.go` 等，约 60+ 处调用。

### 模式二：响应字段自动翻译（I18nString）

`i18ntype.I18nString` 在 JSON 序列化时自动翻译：

```go
import "your.module/i18nerror/i18ntype"

// 响应结构体 — 字段类型声明
Name i18ntype.I18nString `json:"name"`

// service 层赋值
Name: i18ntype.I18nString(p.Name),
```

`MarshalJSON` 内部调用 `service.Translate()` 把 i18n key 替换为翻译后的文本。

### 模式三：手动翻译（TranslateI18nKey）

用于需要拼接多个翻译结果的场景：

```go
// internal/service/record.go — 订阅记录的 item name 拼接
itemName = i18nerr.TranslateI18nKey(itemName, s.logger)
itemName = fmt.Sprintf("%s (%s)", itemName, i18nerr.TranslateI18nKey(FreeTrailI18nKey, s.logger))
```

---

## 语言传递链路

```
HTTP 请求 → InjectLanguage 中间件（解析 Accept-Language 头）
         → ctxutil.SetI18NLanguage（按 goroutine ID 存入 sync.Map）
         → service/handler 中 GetI18NLanguage() 取出语言
         → Translate() 按 language + key 查询翻译
         → 请求结束 ClearI18nLanguage() 清理
```

## 翻译查询链路

```
Translate(language, key)
  → 检查"未找到"缓存（5分钟 TTL）→ 命中则直接返回错误
  → 查 Redis Hash (i18n:{service}:{language} → key → value)
  → 值为 "0" 或 miss → 查数据库（blade_i18n_* 表）
  → 普通文本写入 Redis 缓存，富文本不缓存
  → 未找到 → 写入"未找到"缓存，返回错误
```

---

## 数据库表结构

| 表 | 说明 |
|---|---|
| `blade_i18n_language` | 语言定义（如 zh-CN、en-US） |
| `blade_i18n_terminal` | 终端层级（支持父子关系，用于多租户翻译） |
| `blade_i18n_key` | 翻译键（code + data_type：0=纯文本，1=富文本） |
| `blade_i18n_value` | 翻译值（关联 language + key） |

---

## 快速开始

### 依赖

```
github.com/gin-gonic/gin
github.com/redis/go-redis/v9
github.com/pkg/errors
github.com/samber/lo
```

### 1. 构造翻译服务

传入 MySQL 连接（实现 `repository.DBTX` 接口）、Redis 客户端和服务名：

```go
import (
    i18nerr "your.module/i18nerror"
    "your.module/i18nerror/sqlc/repository"
)

service := i18nerr.NewI18NService(db, redisClient, i18nerr.ServiceName("my-app-messages"))
// 内部自动调用 SetI18nService()，之后可通过 GetI18nService() 全局获取
```

> `ServiceName` 用于构造 Redis 缓存 key（`i18n:<serviceName>:<language>`），不同项目应传入不同的值。

### 2. 注册 Gin 中间件

在路由上尽早注册，自动从 `Accept-Language` 头提取语言：

```go
router := gin.Default()
router.Use(i18nerr.InjectLanguage())
```

### 3. 业务层返回 i18n 错误

```go
import "your.module/i18nerror/ierror"

func (s *Service) CreateMaterial(name string) error {
    if name == "" {
        return ierror.NewError(ierror.ParamError)
    }
    return ierror.NewError(ierror.MaterialAlreadyExists)
}
```

### 4. 统一错误响应

在 handler 层识别 `I18nTranslationError` 并翻译：

```go
import (
    i18nerr "your.module/i18nerror"
    "your.module/i18nerror/ierror"
    "github.com/pkg/errors"
)

func HandleError(ctx *gin.Context, err error) {
    var i18nErr ierror.I18nTranslationError
    if errors.As(err, &i18nErr) {
        svc := i18nerr.GetI18nService()
        lang := i18nerr.GetI18NLanguage()
        msg, _ := svc.Translate(ctx.Request.Context(), lang, i18nErr.Error())

        ctx.JSON(http.StatusOK, gin.H{
            "code":    i18nErr.Code,
            "message": msg,
            "success": false,
        })
        return
    }
    // 其它错误处理...
}
```

### 5. Wire 依赖注入

将 `NewI18NService` 放入 `providerSet`，并在路由依赖结构体中声明 `I18NService` 字段，确保 Wire 会构造它：

```go
// wire.go
var providerSet = wire.NewSet(
    wire.Value(i18nerr.ServiceName("my-app-messages")),
    i18nerr.NewI18NService,
)

// router.go
type RouterDeps struct {
    I18NService i18nerr.I18NService  // 必须声明，Wire 才会构造
    // ...其它依赖
}

// injector 中使用
wire.Struct(new(router.RouterDeps), "*")
```

> **注意**：仅把 `NewI18NService` 写进 `providerSet` 不够，Wire 需要有一个类型消费方才会生成构造代码。即使 handler 里只用 `GetI18nService()`，也需要在 `RouterDeps` 中声明 `I18NService` 字段来触发 `NewI18NService` 的执行（内部会调用 `SetI18nService`）。

---

## API 参考

### ierror — 错误码与错误类型

```go
// 错误码定义
type I18nCode struct {
    Code       int
    MessageKey I18nTranslationErrorKey
}

// 错误类型（实现 error 接口，Error() 返回 MessageKey）
type I18nTranslationError struct {
    Code       int
    MessageKey I18nTranslationErrorKey
}

// 构造错误（内部用 errors.WithStack 包装）
func NewError(code I18nCode) error
```

预定义错误码示例：

| 变量 | 业务码 | MessageKey |
|---|---|---|
| `ServerError` | 500 | SERVER_ERROR |
| `ParamError` | 400 | PARAM_ERROR |
| `NoLogin` | 401 | NO_LOGIN |
| `UserNotFound` | 2001 | USER_NOT_FOUND |
| `MaterialNotFound` | 2004 | MATERIAL_NOT_FOUND |
| `SubscriptionExist` | 40021 | SUBSCRIPTION_EXIST |
| `ThirdPartyPaymentError` | 2013 | THIRD_PARTY_PAYMENT_ERROR |

完整列表见 `ierror/errors.go`。

### I18NService — 翻译服务

```go
type ServiceName string  // Redis 缓存 key 中的服务名标识

type I18NService interface {
    Translate(ctx context.Context, language, key string) (string, error)
}

func NewI18NService(dbx repository.DBTX, redisClient *redis.Client, serviceName ServiceName) I18NService
```

### 语言上下文（协程级）

```go
SetI18NLanguage(lang string)   // 设置当前协程语言
GetI18NLanguage() string       // 获取（默认 ""）
ClearI18nLanguage()            // 清除

SetI18nService(I18NService)    // 设置全局服务实例
GetI18nService() I18NService   // 获取全局服务实例
```

### TranslateI18nKey — 快捷翻译

翻译失败时返回 key 本身作为 fallback：

```go
msg := i18nerr.TranslateI18nKey("FREE_TRIAL", logger)
```

---

## 错误处理流程

```
HTTP Request (Accept-Language: en-US)
    │
    ▼
InjectLanguage 中间件 → SetI18NLanguage("en-US")
    │
    ▼
业务层返回 ierror.NewError(ierror.PlanNotFound)
    │
    ▼
Handler 层 HandleError → errors.As 识别 I18nTranslationError
    │
    ▼
GetI18nService().Translate(ctx, "en-US", "PLAN_NOT_FOUND")
    │  Redis 缓存 → 数据库查询 → 回写缓存
    ▼
HTTP 200 {code: 40022, msg: "Plan not found", success: false}
```

> 业务错误统一返回 HTTP 200，通过 body 中的 `code` 区分错误类型（对齐 Java 客户端行为）。

---

## Redis 缓存策略

| 缓存类型 | Key 格式 | 说明 |
|---|---|---|
| 翻译缓存 | `i18n:{serviceName}:{language}` | Hash，field 为 key，value 为翻译文本 |
| 未找到缓存 | `i18n:i18n-no-found-message:{lang}_{key}` | String，TTL 5 分钟，防止缓存穿透 |

- 纯文本（data_type=0）写入 Redis 缓存
- 富文本（data_type=1）不缓存，每次从数据库查询

## 配置项

| 参数 | 类型 | 说明 |
|---|---|---|
| `ServiceName` | `string`（命名类型） | Redis key 中的服务名标识，构造 `NewI18NService` 时传入，用于区分不同项目的缓存 |
| 默认语言 | — | `zh-CN`，Accept-Language 为空时的回退语言 |
