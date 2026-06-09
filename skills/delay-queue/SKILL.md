---
name: delay-queue
description: 适用于当前 CPA 项目中延迟队列任务的新增、修改与排查，尤其是 delay_queue/ 执行器、internal/delay_queue/ 脚本入口、delay_queue_log 日志表驱动执行和延迟任务入队参数协议。
---

# CPA 延迟队列

当任务涉及延迟投递、定时执行、`delay_queue` 脚本、`delay_queue_log` 状态流转或 `DelayQueuePushParams` 参数协议时，使用这个 skill。

按需加载参考资料：

- 需要判断这次需求该不该触发本 skill 时，读 `references/triggers.md`
- 需要看入队到执行的链路时，读 `references/flow.md`
- 需要核对状态流转和参数协议时，读 `references/status-and-params.md`
- 需要新增延迟任务处理器时，读 `references/example.md`
- 需要避免常见错误写法时，读 `references/anti-patterns.md`

## 优先查看

- `README.md`
- `core/delay_queue.go`
- `delay_queue/main.go`
- `delay_queue/delay_queue.go`
- `internal/delay_queue/mian.go`
- `internal/script/register.go`
- `internal/model/console/delay_queue_log.go`

## 需要保持的项目事实

- 延迟队列的“入队”和“执行”是两层职责：
  - `delay_queue.NewDelayQueue().Push(...)` 负责写入日志表
  - `delay_queue.NewDelayQueue().Handler()` 负责扫描待执行任务并分发执行
- `internal/delay_queue/mian.go` 是脚本分发入口的一部分，虽然文件名有拼写问题，但当前项目就是这样引用的，不能随手改名。
- 延迟任务的实际执行记录依赖 `delay_queue_log` 表，状态值含义是：
  - `1` 待执行
  - `2` 执行中
  - `3` 已完成
  - `4` 失败
- `ExecCmd` 对应的是 `delay_queue/main.go` 里的处理器注册，不是 `internal/script/register.go` 的脚本名。
- `ExecTime` 优先级高于 `DelayTime`，两者同时存在时以 `ExecTime` 为准。

## 工作流程

1. 先判断需求属于哪一层：
   - 入队协议：`core/delay_queue.go` 和调用方
   - 扫描执行器：`delay_queue/delay_queue.go`
   - 任务处理器注册：`delay_queue/main.go`
   - 脚本入口：`internal/delay_queue/mian.go`
2. 新增延迟任务处理器时，先实现 `delay_queue/main.go` 约定的 `Handler(params string) core.DelayQueueResult`。
3. 确认新的 `ExecCmd` 已在 `delay_queue/main.go` 注册，否则任务会落到失败状态。
4. 如果改动了入队逻辑，检查：
   - `Name`
   - `ExecCmd`
   - `Params`
   - `DelayTime`
   - `ExecTime`
5. 如果改动了执行逻辑，检查任务状态流转是否完整，至少覆盖待执行、执行中、成功、失败四种状态。
6. 如果需要从定时脚本触发延迟队列，保持“脚本入口”和“队列执行器”分离，不要把扫描逻辑塞回其他脚本里。

## 常见问题

- 只写了入队逻辑，但没有注册对应的 `ExecCmd` 处理器。
- 把 `internal/script/register.go` 里的脚本名误认为延迟队列处理器名。
- 修改了 `ExecTime` 计算逻辑，却忘了 `ExecTime` 本来就覆盖 `DelayTime`。
- 任务失败时没有写回状态和结果，导致日志表里卡在执行中。
- 看到 `internal/delay_queue/mian.go` 文件名拼写不对就直接重命名，结果把现有引用打断。

## 校验要求

- 对修改过的文件执行格式化：`gofmt -w <修改的文件>`
- 编译校验（优先，最快）：
  - `GOCACHE=$(pwd)/.tmp/gocache go build ./delay_queue/... ./internal/delay_queue/...`
- 最小范围测试：
  - `GOCACHE=$(pwd)/.tmp/gocache go test ./delay_queue ./internal/delay_queue -run TestDoesNotExist`
- 如果还改到了脚本分发，再补：
  - `GOCACHE=$(pwd)/.tmp/gocache go test ./internal/script -run TestDoesNotExist`
- 如果改到了状态流转或模型字段，明确检查是否会破坏已有日志表数据兼容性。


## 不适用场景

- Kafka 或 Redis 消费者工作，位置在 `exec/`
- 普通定时脚本调度问题，主要在 `crontab/` 和 `internal/script/`
- 推广 click/callback 主业务链路，主要在 `exec/promotion/` 和 `internal/main_progress/`
