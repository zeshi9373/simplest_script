# 延迟队列链路

## 入队

1. 调用方执行 `delayqueue.NewDelayQueue().Push(...)`
2. 根据 `ExecTime` 或 `DelayTime` 计算执行时间
3. 写入 `delay_queue_log`

## 执行

1. `internal/delay_queue/mian.go` 暴露脚本入口
2. 该入口调用 `delayqueue.NewDelayQueue().Handler()`
3. `Handler()` 查询未来一分钟内待执行任务
4. 每条任务按 `ExecCmd` 分发到 `delay_queue/main.go` 注册的处理器
5. 按返回结果写回日志状态和结果

## 关键点

- 入队和执行不是一个接口
- 延迟队列处理器注册表和 `internal/script/register.go` 不是同一个表
