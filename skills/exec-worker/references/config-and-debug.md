# 配置与排查

## 配置表关键字段

| 字段 | 说明 |
|------|------|
| `ExecCmd` | 代码注册名，必须与 `exec/registry.go` 中的 key 完全一致 |
| `Topic` | Kafka 消费主题（Kafka 消费者专用，Redis 消费者留空） |
| `Key` | Redis list key（Redis 消费者专用，Kafka 消费者留空） |
| `GroupId` | Kafka consumer group ID |
| `Progress` | 初始消费者并发数 |
| `MaxProgress` | 消费者扩容上限（必须 > Progress 才能触发扩容） |
| `ProgressLagLimit` | 触发扩容的堆积条数阈值 |
| `ProgressAvgMsgcount` | 每增加一个临时消费者对应的消息量参考值 |

## 排查步骤

### 消费者没有启动

1. 确认 `ExecCmd` 在 `exec/registry.go` 有注册
2. 查配置表：`type = 1`、`status = 1`、`SCRIPT_PARTITION` 与当前机器匹配
3. 确认消费源类型：Kafka 消费者填 `Topic`，Redis 消费者填 `Key`（填错则静默不消费）
4. 确认 `main.go` 以 `queue` 模式启动（非 queue 模式不启动消费者）

### 消费者启动了但不扩容

1. `MaxProgress` 是否 > `Progress`（等于时永远不扩容）
2. `ProgressLagLimit` 是否设置合理（设为 0 会频繁扩容；设太大则永远不触发）
3. `CheckKafkaProgress()` / `CheckRedisProgress()` 执行间隔是否正常（每 3 分钟一次）
4. 临时扩容的消费者带 3 分钟 context timeout，到期自动退出

### 消息消费了但没有落库

1. 查 `Consume` 里的 `*status` 赋值逻辑：`false` 代表处理失败
2. 确认业务逻辑在 `exec/<domain>/` 或 `internal/` 里，不在 `exec/main.go`
3. 查是否有全局变量初始化顺序依赖问题（A consumer 依赖 B consumer 的全局状态）

## 常见配置错误

| 现象 | 可能原因 |
|------|----------|
| 消费者完全不启动 | ExecCmd 未注册；配置表 status≠1；type≠1；分区不匹配 |
| 消费者启动但不消费 | Kafka 消费者填了 Key（应填 Topic）；Redis 消费者填了 Topic（应填 Key） |
| 消费了但不扩容 | MaxProgress ≤ Progress；ProgressLagLimit 设太大 |
| 部署后老消费者消失 | ExecCmd 改名了，配置表没同步 |
