# 执行链路

## 入口

1. `main.go` 在 `queue` 模式下调用 `exec.Init()`
2. `exec.Init()` 从脚本配置表加载 `type = 1`、当前 `SCRIPT_PARTITION`、`status = 1` 的记录
3. `exec/registry.go` 提供 `ExecCmd -> Consumer` 注册映射
4. 命中配置后，按 `Topic` 或 `Key` 启动 Kafka / Redis 消费

## 扩容

1. `progressCronCheck()` 每 3 分钟调用 `CheckKafkaProgress()`
2. Kafka 用 consumer lag 判断，Redis 用 `LLen` 判断
3. 当堆积超过 `ProgressLagLimit` 时，走 `ProgressAdd()`
4. 临时扩出来的消费者带 3 分钟 context timeout

## 关键点

- `Scripts` 的 key 可能是 `Topic`，也可能是 `Key`
- `Progress` 是初始并发
- `MaxProgress` 是最大扩容上限
- `TopicNum` 用来跟踪临时扩容中的消费者数量
