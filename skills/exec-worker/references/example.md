# 最小样例

## 新增一个 Kafka 消费者

### 步骤

1. 在 `exec/<domain>/` 下实现消费者
2. 在 `exec/registry.go` 注册 `ExecCmd`
3. 在配置表补充对应记录（`exec_cmd`、`topic`、`progress` 等字段）

### 实现骨架

```go
// exec/demo/consumer.go
package demo

import "encoding/json"

type DemoConsumer struct{}

type demoMsg struct {
    ID   int64  `json:"id"`
    Data string `json:"data"`
}

func (c *DemoConsumer) Consume(msg string, status *bool) {
    var m demoMsg
    if err := json.Unmarshal([]byte(msg), &m); err != nil {
        // 解析失败：标记成功（避免无限重试），打错误日志
        *status = true
        return
    }

    if err := c.handle(m); err != nil {
        *status = false  // 标记失败，框架会重试或记录
        return
    }

    *status = true
}

func (c *DemoConsumer) handle(m demoMsg) error {
    // 业务逻辑
    return nil
}
```

### 注册（exec/registry.go）

```go
func InitEntry() {
    Entry["demo_consumer"] = &demo.DemoConsumer{}
}
```

### 配置表记录示例

| 字段 | 值 |
|------|-----|
| exec_cmd | demo_consumer |
| type | 1 |
| topic | your.kafka.topic |
| progress | 2 |
| max_progress | 10 |
| progress_lag_limit | 1000 |
| progress_avg_msgcount | 500 |
| status | 1 |

---

## 新增一个 Redis list 消费者

与 Kafka 消费者实现相同，差异只在配置表：

| 字段 | 值 |
|------|-----|
| exec_cmd | demo_redis_consumer |
| type | 1 |
| key | your:redis:list:key  ← 注意：填 key，不填 topic |
| progress | 1 |
| status | 1 |

---

## 排查流程：消费者没有启动

按顺序逐步排查：

1. **ExecCmd 注册**：在 `exec/registry.go` 确认 `Entry["your_cmd"]` 存在
2. **配置表记录**：查询配置表，确认 `exec_cmd` 匹配、`type=1`、`status=1`、分区匹配当前机器
3. **消费源类型**：配置里填的是 `topic` 还是 `key`？Kafka 用 topic，Redis 用 key
4. **扩容参数**：如果消费了但太慢，查 `max_progress` 是否 > `progress`，`progress_lag_limit` 是否合理

---

## 完整链路验证清单

- [ ] `exec/registry.go` 已注册新的 ExecCmd
- [ ] 未修改任何已上线的 ExecCmd、topic 或 key 名称
- [ ] 业务逻辑在 `exec/<domain>/` 或 `internal/`，没有写进 `exec/main.go`
- [ ] `Consume` 的 `*status` 赋值逻辑正确（成功 true，失败 false）
- [ ] `gofmt -w` 格式化
- [ ] `go build ./exec/...` 编译通过
