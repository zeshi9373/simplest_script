# 不推荐写法

## 1. 混淆脚本注册和延迟任务注册

错误方向：

```go
// internal/script/register.go — 错误地在这里注册延迟任务处理器
HandlerEntry["my_delay_task"] = &myhandler.MyHandler{}

// 结果：延迟队列扫描到这个 ExecCmd 时找不到对应处理器，任务直接失败
```

正确方向：

```go
// delay_queue/main.go — 延迟队列自己的注册表
HandlerEntry["my_delay_task"] = &myhandler.MyHandler{}
```

原因：

- `internal/script/register.go` 是定时脚本的分发表，由子进程调用
- `delay_queue/main.go` 是延迟队列的处理器表，由扫描执行器调用
- 两套分发机制完全独立，互不可见

## 2. 只入队，不检查状态流转

错误方向：

```go
// 入队后直接认为任务会执行
err := delayqueue.NewDelayQueue().Push(params)
if err == nil {
    log.Info("任务已投递")
    return
}
// 从不查 delay_queue_log 表确认状态
```

正确方向：

- Push 成功只代表写入了日志表（状态=1，待执行）
- 实际执行依赖扫描执行器定时调度
- 排查时查 `delay_queue_log`：状态卡在 `1` = 没被扫描到；卡在 `2` = 执行中卡住；`4` = handler 失败

## 3. 随手重命名 `internal/delay_queue/mian.go`

错误方向：

```bash
# 因为看到 mian.go 有拼写错误就直接改
mv internal/delay_queue/mian.go internal/delay_queue/main.go
```

原因：

- 当前项目按这个文件名组织了 package 引用
- 改名后必须同步所有 import 路径，改前先运行 `grep -r "delay_queue/mian" .` 确认影响范围
- 非必要不改，先做功能性修改

## 4. 任务失败时不写回状态

错误方向：

```go
func (h *Demo) Handler(params string) core.DelayQueueResult {
    err := doWork(params)
    if err != nil {
        // 只打日志，不返回失败状态
        log.Error(err)
        return core.DelayQueueResult{Status: true} // 错误：返回了成功
    }
    return core.DelayQueueResult{Status: true}
}
```

正确方向：

```go
func (h *Demo) Handler(params string) core.DelayQueueResult {
    err := doWork(params)
    if err != nil {
        return core.DelayQueueResult{Status: false, Data: err.Error()}
    }
    return core.DelayQueueResult{Status: true, Data: "ok"}
}
```

原因：

- `Status: false` 会让扫描执行器把日志表状态更新为 `4`（失败）
- 如果返回 `Status: true` 掩盖了错误，日志表显示成功，实际数据没处理

## 5. ExecTime 和 DelayTime 同时传导致混淆

错误方向：

```go
// 同时设置两个，以为会叠加延迟
Push(core.DelayQueuePushParams{
    DelayTime: 300,               // 5 分钟后
    ExecTime:  time.Now().Unix(), // 立即执行
})
// 结果：ExecTime 优先，任务立即执行，DelayTime 被忽略
```

原因：

- `ExecTime` 有值时，`DelayTime` 完全被忽略
- 两者只用其一：精确时间用 `ExecTime`，相对延迟用 `DelayTime`
