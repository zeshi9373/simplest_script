# 日志与注册

## 注册点

- `internal/script/register.go` 是脚本名到 handler 的注册表
- 改脚本名时，配置里的 `exec_cmd` 也要同步

## 日志链路

- 开启日志时，`crontab.execHandler()` 会先写一条 running 记录
- 主程序非 `queue` 分支执行完 `script.Exec()` 后调用 `crontab.UpdateCrontabLog()`
- `uk` 用来把开始记录和结束记录串起来

## 排查点

1. 配置里是否开启日志
2. 子进程参数顺序是否还是 `uk exec_cmd params`
3. handler 返回值是否仍为 `*crontab.Result`
