# 最小样例

## 新增一个延迟队列处理器

### 步骤

1. 实现处理器（签名固定：`Handler(params string) core.DelayQueueResult`）
2. 在 `delay_queue/main.go` 的注册表里添加 `ExecCmd`
3. 在调用方使用 `Push()` 写入任务

### 实现骨架

```go
// internal/handler/demo_delay/demo.go
package demodelay

import (
    "simplest_script/core"
    "encoding/json"
)

type Demo struct{}

type demoParams struct {
    UserID int64  `json:"user_id"`
    Action string `json:"action"`
}

func (h *Demo) Handler(params string) core.DelayQueueResult {
    var p demoParams
    if err := json.Unmarshal([]byte(params), &p); err != nil {
        return core.DelayQueueResult{Status: false, Data: "invalid params: " + err.Error()}
    }

    // 业务逻辑
    if err := processAction(p.UserID, p.Action); err != nil {
        return core.DelayQueueResult{Status: false, Data: err.Error()}
    }

    return core.DelayQueueResult{Status: true, Data: "ok"}
}
```

### 注册（delay_queue/main.go）

```go
func InitEntry() {
    // 注意：这里是延迟队列的注册表，不是 internal/script/register.go
    HandlerEntry["demo_delay_task"] = &demodelay.Demo{}
}
```

### 入队（调用方）

```go
import "simplest_script/core"

// 相对延迟：60 秒后执行
delayqueue.NewDelayQueue().Push(core.DelayQueuePushParams{
    Name:      "demo 延迟任务",
    ExecCmd:   "demo_delay_task",      // 必须与注册表里的 key 一致
    Params:    `{"user_id":123,"action":"notify"}`,
    DelayTime: 60,
})

// 精确时间：2 分钟后执行（ExecTime 优先于 DelayTime）
delayqueue.NewDelayQueue().Push(core.DelayQueuePushParams{
    Name:    "demo 精确时间任务",
    ExecCmd: "demo_delay_task",
    Params:  `{"user_id":123,"action":"notify"}`,
    ExecTime: time.Now().Add(2 * time.Minute).Unix(),
})
```

---

## 排查流程：任务没有执行

1. **确认入队成功**：查 `delay_queue_log`，看是否有状态=1 的记录
2. **确认 ExecCmd 注册**：在 `delay_queue/main.go` 找对应的 key
3. **确认执行时间**：`exec_time` 字段是否已过，还是还在未来
4. **查看状态流转**：
   - 状态=1（待执行）：扫描执行器还没扫到，或执行时间未到
   - 状态=2（执行中）：handler 卡住了，检查是否有死循环或长耗时操作
   - 状态=4（失败）：查 `result` 字段看错误信息
5. **确认扫描脚本在运行**：`internal/delay_queue/mian.go` 对应的定时脚本是否在配置表里启用

---

## 完整链路验证清单

- [ ] handler 在 `delay_queue/main.go` 已注册（不是 internal/script/register.go）
- [ ] `ExecCmd` 名称与注册表 key 完全一致
- [ ] handler 失败路径返回 `Status: false`
- [ ] `ExecTime` 和 `DelayTime` 只用了其中一个（避免混用）
- [ ] `gofmt -w` 格式化
- [ ] `go build ./delay_queue/...` 编译通过
