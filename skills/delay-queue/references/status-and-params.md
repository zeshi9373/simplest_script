# 状态与参数

## 状态值

- `1`: 待执行
- `2`: 执行中
- `3`: 已完成
- `4`: 失败

## 参数协议

- `Name`: 任务名称
- `ExecCmd`: 延迟任务处理器名
- `Params`: 传给处理器的字符串参数
- `DelayTime`: 延迟秒数
- `ExecTime`: 指定执行时间；有值时覆盖 `DelayTime`

## 排查点

1. `ExecCmd` 是否在 `delay_queue/main.go` 注册
2. 任务是否成功从 `1` 变到 `2`
3. handler 返回的 `core.DelayQueueResult.Status` 是否符合成功语义
