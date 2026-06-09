# 调度链路

## 主流程

1. `main.go` 在 `queue` 模式下注册定时任务
2. `crontab.Init()` 从脚本配置表加载 `type = 2`、当前 `SCRIPT_PARTITION`、`status = 1` 的记录
3. 到点后 `execHandler()` 拉起子进程
4. 子进程走主程序的非 `queue` 分支
5. 非 `queue` 分支调用 `internal/script.Exec()`
6. `internal/script/register.go` 按脚本名分发到具体 handler

## 关键点

- cron 调度和脚本执行不是同一层
- `queue` 模式负责常驻调度
- 单次脚本执行靠子进程参数驱动
