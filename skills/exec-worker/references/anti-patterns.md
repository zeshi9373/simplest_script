# 不推荐写法

## 1. 把业务逻辑塞进 `exec/main.go`

错误方向：

```go
// exec/main.go — 错误示范
func Init() {
    for _, script := range scripts {
        if script.ExecCmd == "promotion_callback" {
            // 直接在这里写业务判断
            if script.Topic == "xxx" {
                go startSpecialConsumer()
            }
        }
    }
}
```

正确方向：

- `exec/main.go` 只负责按 `ExecCmd` 查注册表、按 `Topic`/`Key` 启动消费者
- 业务判断落在 `exec/<domain>/` 下的具体消费者文件或 `internal/` 里

原因：

- 入口混入业务逻辑后，每次修改都要改核心启动流程，风险高
- 无法单独测试业务逻辑

## 2. 变更已上线的 `ExecCmd`、topic 或 key

错误方向：

```go
// exec/registry.go — 直接改名
// 原来
Entry["promotion_callback"] = &PromotionCallbackConsumer{}
// 改成
Entry["promotion_callback_v2"] = &PromotionCallbackConsumer{}
```

原因：

- 配置表里的记录仍然是 `exec_cmd = "promotion_callback"`
- 改名后消费者注册存在但找不到，配置表里的消费者静默不启动，不报错
- 必须同步更新配置表，或采用双注册过渡

## 3. 用隐藏全局变量串多个 handler

错误方向：

```go
// exec/click/handler.go
var sharedBatchWriter *BatchWriter  // 全局变量，由 click handler 初始化

// exec/callback/handler.go
func (c *CallbackConsumer) Consume(msg string, status *bool) {
    // 依赖 click handler 先初始化了 sharedBatchWriter
    click.sharedBatchWriter.Add(...)
}
```

正确方向：

```go
// 提取公共组件，通过依赖注入共享
type BatchWriter struct { ... }

func NewBatchWriter() *BatchWriter { ... }

// 两个 consumer 各自持有或共享同一个实例（通过 svc 或初始化时注入）
```

原因：

- 初始化顺序依赖是隐式的，后来的开发者很难发现
- 一旦启动顺序变化，就会 nil pointer panic 或静默使用错误数据

## 4. 扩容参数设置不合理

错误方向：

```go
// 配置表设置了错误的扩容参数
// MaxProgress = 1  （和 Progress 一样，永远不会扩容）
// ProgressLagLimit = 0  （lag 为 0 时就触发扩容，一直扩容）
```

排查方向：

- `Progress`: 初始消费者数（通常 1~3）
- `MaxProgress`: 扩容上限，必须 > `Progress` 才能扩容
- `ProgressLagLimit`: 触发扩容的堆积条数阈值，设为 0 或过小会导致频繁扩容
- `ProgressAvgMsgcount`: 每增加一个消费者对应的消息量参考值

## 5. Kafka 和 Redis 消费者配置混填

错误方向：

```
# Redis 消费者配置里填了 Topic（应填 Key）
exec_cmd: redis_consumer
topic: some_topic   # 错误，Redis 消费者走 Key，不走 Topic
```

原因：

- `Topic` 用于 Kafka 消费者定位
- `Key` 用于 Redis list 消费者定位
- 填错字段后消费者找不到数据源，静默不消费
