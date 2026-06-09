---
name: exec-worker
description: 适用于当前 CPA 项目中 Kafka 或 Redis 消费任务的新增、修改与排查，尤其是 exec/ 目录下的消费者实现、ExecCmd 注册、配置驱动启动、推广回传处理和消息堆积扩容逻辑。
---

# CPA 消费执行器

当任务发生在 `exec/` 目录，或者用户在排查“异步消费者为什么没启动、没扩容、没落库”时，使用这个 skill。

按需加载参考资料：

- 需要判断这次需求该不该触发本 skill 时，读 `references/triggers.md`
- 需要看启动和扩容链路时，读 `references/flow.md`
- 需要核对配置字段语义和排查点时，读 `references/config-and-debug.md`
- 需要新增消费者时，读 `references/example.md`
- 需要避免常见错误写法时，读 `references/anti-patterns.md`

## 优先查看

- `README.md`
- `main.go`
- `exec/README.md`
- `exec/main.go`
- `exec/registry.go`
- `exec/check_progress.go`
- `exec/<domain>/` 下的相关文件
- `internal/` 下相关的 types、services、models

## 需要保持的项目事实

- `queue` 模式才是常驻运行模式，会初始化 crontab、exec 消费者和 resident 任务。
- `exec` 的启动是配置驱动的，只在代码里注册还不够，还依赖脚本配置表。
- 消费源可能是 Kafka topic，也可能是 Redis list key，要先区分。
- 扩容行为依赖 `Progress`、`MaxProgress`、`ProgressLagLimit`、`ProgressAvgMsgcount` 等字段。
- 现有 `ExecCmd` 值和配置表语义默认保持兼容，除非用户明确要求变更。

## 工作流程

1. 先判断目标消费者是 Kafka 还是 Redis list。
2. 找到对应的 `ExecCmd`，确认它已经在 `exec/registry.go` 注册。
3. 核对配置加载和索引方式：
   - topic 消费者走 `Topic`
   - Redis 消费者走 `Key`
4. `exec/` 只放启动、调度、扩容逻辑，业务逻辑放在 `exec/<domain>/` 或 `internal/`。
5. 如果多个 handler 共享批处理或落库逻辑，应提取公共组件，不要通过隐藏全局变量耦合。
6. 如果改了扩容逻辑，要同时检查初始并发数和 lag 驱动的临时扩容路径。

## 常见问题

- 只在代码里注册消费者，但忘了运行时仍依赖脚本配置表。
- 修改了 `ExecCmd`、topic 或 key，导致已有配置静默失效。
- 把业务分支写进 `exec/main.go`，而不是写进具体领域消费者。
- 增加跨文件全局变量，导致一个 handler 依赖另一个 handler 的初始化顺序。
- 一上来就跑 `go test ./...`，而不是先做更快的包级校验。

## 校验要求

- 对修改过的文件执行格式化：`gofmt -w <修改的文件>`
- 编译校验（优先，最快）：
  - `GOCACHE=$(pwd)/.tmp/gocache go build ./exec/...`
- 最小范围测试：
  - `GOCACHE=$(pwd)/.tmp/gocache go test ./exec/... -run TestDoesNotExist`
- 如果还改到了其他业务包，再补对应包的编译或测试校验。

## 不适用场景

- 纯定时脚本工作，位置在 `internal/script` 或 `crontab/`
- 常驻后台任务，位置在 `resident/`
- 延迟队列处理器，位置在 `internal/delay_queue/`
