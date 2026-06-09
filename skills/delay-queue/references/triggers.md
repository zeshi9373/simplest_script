# 触发提示

## 应触发本 skill 的提问

- “新增一个延迟任务”
- “为什么 delay_queue 没执行”
- “帮我加一个延迟任务处理器”
- “`delay_queue_log` 一直停在待执行”
- “`ExecTime` 和 `DelayTime` 应该怎么用”

## 边界案例

- 如果问题主要是普通 cron 调度和脚本执行，改用 `cron-script`
- 如果问题主要是 Kafka/Redis 消费者，改用 `exec-worker`
- 如果问题主要是推广业务回传和扣量，改用 `business-pipeline`

## 首轮输出建议

1. 先判断问题在入队层、扫描层还是处理器层
2. 再确认 `ExecCmd` 注册和日志状态流转
3. 最后说明最小复现或验证方式
