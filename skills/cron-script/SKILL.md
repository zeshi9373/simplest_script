---
name: cron-script
description: 适用于当前 CPA 项目中定时脚本的新增、修改与排查，尤其是 crontab/ 调度、internal/script/ 注册分发、单次脚本执行入口和执行日志记录链路。
---

# CPA 定时脚本

当任务发生在 `crontab/`、`internal/script/`、`internal/handler/`，或者用户在排查“定时脚本为什么没执行、没分发、没记日志”时，使用这个 skill。

按需加载参考资料：

- 需要判断这次需求该不该触发本 skill 时，读 `references/triggers.md`
- 需要看调度到执行的完整链路时，读 `references/flow.md`
- 需要核对日志记录和脚本注册约束时，读 `references/logging-and-registration.md`
- 需要新增脚本 handler 时，读 `references/example.md`
- 需要避免常见错误写法时，读 `references/anti-patterns.md`

## 优先查看

- `README.md`
- `main.go`
- `crontab/main.go`
- `internal/script/main.go`
- `internal/script/register.go`
- 相关脚本实现文件，例如 `internal/handler/` 或 `internal/delay_queue/`
- 相关日志模型或配置表模型

## 需要保持的项目事实

- 这个项目有两类运行方式：`queue` 模式常驻运行；非 `queue` 分支执行单次脚本。
- 定时任务本身不是直接执行业务函数，而是由 `crontab` 读取配置后启动子进程，再走 `internal/script` 分发。
- `internal/script/register.go` 里的注册是脚本分发入口，没注册就不会执行。
- 脚本参数入口是字符串，handler 自己负责解析。
- 如果开启日志，执行前后会写入定时任务日志表，`uk` 是串联执行记录的关键值。

## 工作流程

1. 先判断需求属于哪一层：
   - 调度层：`crontab/`
   - 分发层：`internal/script/`
   - 业务实现层：`internal/handler/` 或其他业务目录
2. 新增脚本时，先实现 handler，再在 `internal/script/register.go` 注册脚本名。
3. 核对脚本调用参数，确认传入的是字符串，并保持返回值为 `*crontab.Result`。
4. 如果改动了调度行为，检查配置表读取条件、cron 表达式使用方式和子进程启动参数。
5. 如果改动了日志行为，检查 `uk`、状态值、结果序列化和结束时间更新逻辑。
6. 如果脚本实际执行的是其他模块能力，优先保持“调度与业务分离”，不要把业务逻辑塞回 `crontab/main.go`。

## 常见问题

- 只写了业务 handler，但忘了在 `internal/script/register.go` 注册。
- 修改脚本名后，没有同步配置里的 `exec_cmd`。
- 把 cron 调度和单次执行混为一谈，误以为 `queue` 模式会直接执行脚本函数。
- 忽略 `uk` 和日志更新链路，导致数据库里只有 running 没有结束状态。
- 脚本参数仍然是裸字符串，却按结构体直接使用，缺少解析和错误处理。

## 校验要求

- 对修改过的文件执行格式化：`gofmt -w <修改的文件>`
- 编译校验（优先，最快）：
  - `GOCACHE=$(pwd)/.tmp/gocache go build ./crontab/... ./internal/script/...`
- 最小范围测试：
  - `GOCACHE=$(pwd)/.tmp/gocache go test ./crontab ./internal/script -run TestDoesNotExist`
- 如果改到了具体 handler，再补对应包的编译或测试校验。
- 如果修改影响脚本名或参数协议，明确检查是否会破坏已有配置表数据。

## 不适用场景

- Kafka 或 Redis 消费者工作，位置在 `exec/`
- 常驻后台任务，位置在 `resident/`
- 推广 click/callback 主流程改造，主要发生在 `internal/main_progress/` 或 `exec/promotion/`
