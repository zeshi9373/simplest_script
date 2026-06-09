# exec

`exec` 是脚本消费执行入口，负责三类事情：

1. 从 `console_script_config` 加载启用中的执行配置
2. 按 Kafka topic 或 Redis list 启动对应消费者
3. 定时检查堆积，按配置动态补充短时消费者

## 目录职责

- `main.go`: 执行配置加载、消费者启动、动态扩容
- `registry.go`: 执行器注册表
- `check_progress.go`: Kafka/Redis 堆积检查
- `promotion/`: 推广类执行器

## promotion 子目录

- `callback.go`: 推广回传消费处理
- `click.go`: 推广点击消费处理
- `batch_writer.go`: 点击/回传共用的批量落库缓冲

## 扩展方式

1. 在 `exec/<domain>` 下实现 `kafkaclient.SyncConsumer`
2. 在 `registry.go` 中注册 `ExecCmd`
3. 在后台配置对应 `exec_cmd`、`topic`/`key`、并发数和扩容阈值
