# 触发提示

## 应触发本 skill 的提问

- “新增一个 Kafka 消费者”
- “为什么 exec 里的消费者没有启动”
- “帮我排查 queue 模式下 Redis list 没消费”
- “给这个 topic 增加堆积扩容逻辑”
- “调整 `ExecCmd` 对应的消费处理器”

## 边界案例

- 如果问题主要是 cron 配置、脚本名分发、执行日志，改用 `cron-script`
- 如果问题主要是推广 click/callback 业务规则，优先用 `business-pipeline`
- 如果问题主要是延迟任务状态流转或入队协议，改用 `delay-queue`

## 首轮输出建议

1. 先说明这是 Kafka 还是 Redis 消费链路
2. 再说明问题在注册、启动、扩容还是 handler 逻辑
3. 最后给出最小验证命令
