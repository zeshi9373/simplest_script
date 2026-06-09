# 触发提示

## 应触发本 skill 的提问

- “新增一个定时脚本”
- “这个 cron 为什么没执行”
- “帮我加一个 internal/script 的 handler”
- “为什么数据库里只有 running 没有 success”
- “脚本参数怎么从 crontab 传到 handler”

## 边界案例

- 如果问题主要是 Kafka/Redis 消费者，改用 `exec-worker`
- 如果问题主要是延迟队列扫描或 `delay_queue_log`，改用 `delay-queue`
- 如果问题主要是推广回传、扣量、去重，优先用 `business-pipeline`

## 首轮输出建议

1. 先区分调度层、分发层、业务层
2. 再确认脚本名注册和参数协议是否被改动
3. 最后说明日志链路是否完整
