# core

项目基础设施层，提供配置加载、数据库/Redis/Kafka/RabbitMQ/ES 客户端、工具函数、监控和告警能力。业务代码通过 import 本包使用，不应在此层写业务逻辑。

---

## 目录结构

```
core/
├── consts.go           # 全局常量（状态码、DB/Redis 标识、过期时间）
├── delay_queue.go      # 延迟队列参数与结果类型定义
├── env.go              # 环境变量工具（StatusIsEnv）
├── response.go         # HTTP 响应结构体与构造函数
├── conf/
│   ├── h.go            # 配置结构体定义（Config、RedisConfig …）
│   └── m.go            # MustLoad：加载 YAML 配置文件，失败则 Fatal
├── logger/
│   └── logger.go       # 自定义结构化日志（JSON/文本、异步写入、按大小轮转）
├── monitor/
│   └── server.go       # 运行时指标采集（内存、GC、协程），定时推送飞书
├── svc/
│   ├── service.go      # MySQL（NewDb）和 Redis（NewRedis）客户端单例
│   ├── elasticsearch.go# ES 类型化客户端封装（CRUD + 分页搜索）
│   ├── rabbitmq.go     # RabbitMQ 连接池与收发封装
│   ├── kafkaclient/
│   │   ├── consumer.go # Kafka 消费者（SyncConsumer 接口 + 重连）
│   │   └── producer.go # Kafka 生产者连接池（GetWriter / PutWriter）
│   └── redislist/
│       └── main.go     # Redis list 消费者（BRPop 驱动，失败自动重推）
├── tool/
│   ├── common.go       # Random、RandString、GetTodayZeroTime
│   ├── crypto.go       # MD5、HmacSha256、AES-CBC 加解密、AES-ECB 加密、SHA256
│   ├── http.go         # HTTP 客户端封装（共享连接池，支持 GET/POST/PostByForm）
│   ├── jwt.go          # JWT HS256 生成与解析
│   ├── password.go     # bcrypt 密码加密与校验
│   ├── response.go     # （空文件，保留占位）
│   ├── slice.go        # IsInSlice[T comparable] 泛型切片查找
│   ├── strings.go      # StringsLimitLength（按字节截断）
│   ├── time_range.go   # ParseTimeRange、DateStringToTime、时间布局常量
│   ├── url.go          # IsURLEncoded、ParseURLEx、FormatParas（排序参数字符串）
│   ├── util.go         # ConvertStringIdsToInt、ConvertInt64IdsToString、LimitStringLength、ParseOsType
│   ├── uuid.go         # Uuid()（基于 google/uuid）
│   ├── crypto/
│   │   └── aes_cbc.go  # 另一套 AES-CBC 封装（自动补齐 key 至 32 字节）
│   └── flowlimit/
│       └── limit.go    # 基于 Redis 的滑动窗口限流（FlowLimit）
├── pool/
│   └── pool.go         # 并发限制协程池（Go/TryGo/Wait，支持错误收集、context 取消、panic 恢复）
└── warning/
    └── main.go         # 基于 Redis 的计数告警（Warning，支持多周期阈值）
```

---

## 常量（consts.go）

### 响应码

| 常量 | 值 | 说明 |
|------|----|------|
| `CodeOK` | `0` | 成功 |
| `CodeFail` | `-1` | 通用失败 |
| `CodeLoginFail` | `-10002` | 登录失败 |

### 数据库标识（传入 `svc.NewDb`）

| 常量 | 值 |
|------|----|
| `DBBusiness` | `"business"` |
| `DBConsole`  | `"console"` |

### Redis 标识（传入 `svc.NewRedis`）

| 常量 | 值 |
|------|----|
| `RDSDefault` | `"default"` |
| `RDSData` | `"data"` |

### 过期时间

| 常量 | 值 |
|------|----|
| `ExpireTimeSecond10` | 10s |
| `ExpireTimeSecond30` | 30s |
| `ExpireTimeMinute` | 1min |
| `ExpireTimeHour` | 1h |
| `ExpireTimeHour3` | 3h |
| `ExpireTimeDay` | 24h |
| `ExpireTimeDay3` | 72h |

---

## 配置（conf）

### 加载配置

```go
var c conf.Config
conf.MustLoad("./etc/dev.yaml", &c)
// 失败时 Fatal，不会继续执行
```

加载后可通过全局变量 `conf.Conf` 在任意位置读取配置，但推荐通过 `svc` 层封装访问。

### Config 主要字段

```yaml
Name: simplest_script
Mode: queue
Mysql:
  Business: "user:pass@tcp(host:3306)/db?..."
Redis:
  Default:
    Addr: "127.0.0.1:6379"
    Pass: ""
    Db: 0
Kafka:
  Default:
    Brokers: "broker1:9092,broker2:9092"
    MaxIdle: 10
Logger:
  Path: "./logs"
  MaxSize: 100   # MB
  MaxBackups: 10
  MaxAge: 7      # 天
```

---

## 服务客户端（svc）

### MySQL

```go
db := svc.NewDb(core.DBBusiness)
// 返回 *gorm.DB，非 Release 环境自动开启 Debug 模式
// 支持最多 3 次自动重连
```

### Redis

```go
rdb := svc.NewRedis(core.RDSDefault)
rdb := svc.NewRedis(core.RDSData)
// 返回 *redis.Client，支持最多 3 次自动重连
```

### Kafka 消费者

```go
// 实现 SyncConsumer 接口
type MyConsumer struct{}

func (c *MyConsumer) Consume(msg string, status *bool) {
    // 处理消息
    *status = true  // true=提交 offset，false=不提交
}

kafkaclient.NewConsumer(core.FlagDefault).
    ConsumerHandlerMessage(ctx, "my-topic", "my-group", &MyConsumer{})
// 出错后自动重启消费协程
```

### Kafka 生产者

```go
producer := kafkaclient.NewProducer(core.FlagDefault)
writer := producer.GetKafkaWriterClient("my-topic")
// 使用完后归还连接池
producer.PutKafkaWriterClient("my-topic", writer)

// 程序退出时关闭所有连接
kafkaclient.CloseKafkaWriterClient()
```

### Redis list 消费者

```go
// 使用同一个 SyncConsumer 接口
redislist.RedisListConsumer(ctx, "my:list:key", &MyConsumer{})
// 消费失败（status=false）自动将消息重新推回队头
// BRPop 阻塞等待，无消息时不占用 CPU
```

### Elasticsearch

```go
es := svc.NewElasticsearch()

// 写入
es.Create(ctx, "index-name", "doc-id", doc)

// 分页搜索
result, err := es.Search(ctx, "index-name", query, page, pageSize, sortOptions...)
// result.Total / result.Data / result.Page / result.PageSize
```

### RabbitMQ

```go
mq := svc.NewRabbitMQ("queue-name", 5) // 5 个 channel 连接池
mq.SetNeedAck(true).Publish("message body")

mq.Consume(func(body string) bool {
    // 返回 true=Ack，false=Reject
    return true
})
```

---

## 延迟队列类型（delay_queue.go）

```go
// 入队参数
type DelayQueuePushParams struct {
    Name      string // 任务名称（仅用于展示）
    ExecCmd   string // 处理器名，对应 delay_queue/main.go 的注册 key
    Params    string // 传给处理器的 JSON 字符串
    DelayTime int64  // 延迟秒数（ExecTime 有值时忽略）
    ExecTime  int    // 指定执行时间（Unix 秒，优先级高于 DelayTime）
}

// 处理器返回值
type DelayQueueResult struct {
    Status bool `json:"status"` // true=成功，false=失败
    Data   any  `json:"data"`
}
```

---

## 环境工具（env.go）

```go
// 判断当前环境是否包含指定 bit
core.StatusIsEnv(core.EnvReleaseValue) // 仅 release 环境返回 true
```

| 常量 | 值 | 环境变量 `SCRIPT_ENV` |
|------|----|------------------|
| `EnvDev` / `EnvDevValue` | `"dev"` / `1` | dev |
| `EnvTest` / `EnvTestValue` | `"test"` / `2` | test |
| `EnvRelease` / `EnvReleaseValue` | `"release"` / `4` | release |
| `EnvPre` / `EnvPreValue` | `"pre"` / `8` | pre |

---

## 响应构造（response.go）

```go
core.Success(requestId, "成功", data)
core.Fail(requestId, "参数错误", nil)
core.LoginFail(requestId, "token 已过期", nil)
// requestId 为空时自动生成（时间戳 + UUID）
```

---

## 工具函数（tool）

### HTTP 客户端

```go
h := tool.NewHttp("https://example.com/api", 5*time.Second)

// JSON POST
body, err := h.Post(map[string]string{"X-Token": "abc"}, jsonBytes)

// GET（自动合并 query 参数）
body, err := h.Get(nil, map[string]string{"page": "1"})

// 表单 POST
body, err := h.PostByForm(nil, map[string]string{"key": "val"})
```

所有请求共享同一个连接池（100 个空闲连接），TLS 跳过证书验证。超时通过 context 控制。

### 加密

```go
tool.Md5("text")
tool.HmacSha256ToHex("key", "data")
tool.Sha256ToHex("text")

// AES-CBC（随机 IV，base64 输出）
cipher, err := tool.AesEncrypt(plainBytes, keyBytes, false)
plain, err  := tool.AesDecrypt(cipherBase64, keyBytes, false)

// AES-ECB（base64 key）
tool.AesEcbEncrypt("plaintext", base64Key)
```

### JWT

```go
token, err := tool.GetJwtToken(secretKey, iat, expireSeconds, uid)
uid, err    := tool.GetJwtUid(tokenString)
```

### 时间

```go
// 常量
tool.DateLayout     // "2006-01-02"
tool.DatetimeLayout // "2006-01-02 15:04:05"

// 解析时间范围（上海时区）
startTs, endTs, err := tool.ParseTimeRange(`{"start_date":"2024-01-01","end_date":"2024-01-31"}`, 72)

// 今日零点时间戳（上海时区）
tool.GetTodayZeroTime()

// 字符串转 time.Time
tool.DateStringToTime("2024-01-01 10:00:00")
```

### 切片

```go
tool.IsInSlice([]int{1, 2, 3}, 2)       // true
tool.IsInSlice([]string{"a", "b"}, "c") // false
```

### 字符串

```go
tool.StringsLimitLength("hello", 3)      // "hel"（按字节截断）
tool.LimitStringLength("你好世界", 2)     // "你好"（按 rune 截断，中文安全）
```

### 限流（FlowLimit）

```go
// 配置：1 分钟内最多 100 次，1 小时内最多 1000 次
limit := flowlimit.NewFlowLimit().
    Make("promotion_callback").
    SetPeriod(100, 60).
    SetPeriod(1000, 3600)

limit.Add(1)       // 计数 +1
ok := limit.Check() // false 表示超限
```

### 告警计数（Warning）

```go
w := warning.NewWarning("task_fail").
    SetPeriod([]warning.PeriodItem{
        {Limit: 10, Duration: 300},  // 5 分钟内不超过 10 次
    })

w.Add(1)                      // 计数 +1
ok, duration := w.Check()     // false + 触发的周期（秒）
nonce := w.GetNonce(300)      // 该周期的 nonce（用于幂等去重告警）
count := w.GetCount(300)      // 当前计数
```

---

## 协程池（pool）

控制同一批任务的最大并发数，并统一收集错误。

### 基础用法

```go
p := pool.New(10) // 最多 10 个协程同时运行

for _, item := range list {
    item := item
    p.Go(func() error {
        return process(item)
    })
}

if err := p.Wait(); err != nil {
    // errors.Join 合并的所有任务错误
    log.Printf("部分任务失败: %v", err)
}
```

### 进阶用法

```go
// context 取消 + 自定义 panic 处理
p := pool.New(10,
    pool.WithContext(ctx),
    pool.WithPanicHandler(func(r any) {
        log.Printf("任务 panic: %v", r)
    }),
)

// 非阻塞提交：队列满时立即返回 false，不阻塞
if !p.TryGo(func() error {
    return doSomething()
}) {
    // 当前并发已满，降级处理或跳过
}

// 主动取消，后续 Go / TryGo 返回 false
p.Cancel()

p.Wait()
```

### API 一览

| 方法 | 说明 |
|------|------|
| `New(n, opts...)` | 创建并发上限为 n 的池 |
| `Go(f func() error) bool` | 提交任务，满员时阻塞；context 取消时返回 false |
| `TryGo(f func() error) bool` | 非阻塞提交，满员或 context 取消时立即返回 false |
| `Wait() error` | 等待全部任务完成，返回合并后的错误 |
| `Cancel()` | 取消 context，阻止后续任务提交 |
| `Errors() []error` | 随时读取已收集的错误列表 |

**选项：**

| 选项 | 说明 |
|------|------|
| `WithContext(ctx)` | 外部 context 取消时自动停止接受任务 |
| `WithPanicHandler(func(r any))` | 自定义 panic 处理；panic 同时也被收集进 error |

> panic 会被自动转换为 error 纳入 `Wait()` 的返回值，不会导致进程崩溃。

---

## 自定义日志（logger）

```go
log := logger.NewLogger("promotion") // 日志写入 {Logger.Path}/promotion/
log.Info("处理完成", logger.Fields{"uk": "xxx", "cost": 12})
log.Error("回传失败", logger.Fields{"err": err.Error()})

// 特性：
// - 异步写入，缓冲区 1024 条
// - 每 200ms 或满 100 条刷盘一次
// - 单文件最大 100MB，自动轮转
// - JSON 格式输出
```

---

## 监控（monitor）

每 5 分钟采集一次运行时指标，通过飞书机器人推送报告。报告包含：

- 当前协程数、CPU 核心数、GOMAXPROCS
- 堆内存使用率（HeapInuse / HeapSys）
- GC 次数、最近 GC 暂停时长

```go
m := monitor.NewEnhancedMonitor()
// 启动时调用，程序退出时：
m.Stop()
```
