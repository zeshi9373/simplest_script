# Script 服务

## 环境配置

正式环境需要设置以下环境变量：

```bash
export SCRIPT_ENV=release        # 环境：dev / test / release
export SCRIPT_PARTITION=script1  # 脚本所在分区，默认 script1，后续机器递增 script2、script3 ...
```

## 构建与运行

```bash
go build .

# queue 模式（常驻运行，启动 crontab + exec 消费者 + resident 任务）
./simplest_script -mode queue

# 单次脚本执行
./simplest_script "" exec_cmd params
```

执行命令可在日志文件 `logs/*.log` 中查看。

## 目录结构

```
/cmd                   命令行入口
/core                  核心非业务公用库
  ~/conf               解析配置文件
  ~/logger             日志记录自定义包
  ~/permission         数据权限
  ~/svc                服务中间件注册
  ~/tool               工具包（非业务）
  ~/warning            预警封装基础包
/crontab               定时任务调度（非业务）
/delay_queue           延迟队列扫描执行器
/etc                   配置文件
/exec                  异步消费业务（Kafka 及 Redis 队列，依赖 script_config 配置表）
  ~/promotion          推广 click / callback 消费入口
/expand                外部业务接口
/internal              内部业务处理
  ~/consts             常量配置
  ~/delay_queue        延迟队列任务处理器入口
  ~/handler            定时任务业务实现
  ~/main_progress      业务主流程封装（公用）
  ~/model              数据库表模型
  ~/script             定时任务注册与分发
  ~/services           业务逻辑
  ~/types              结构体定义
/logs                  日志目录
/resident              常驻进程任务（不需要外部触发器）
/skills                AI 辅助开发技能包（见下方说明）
/test                  测试目录
```

## 生成数据库 Model 文件

配置文件：`generateModel.json`

```json
{
    "module_name": "simplest_script",
    "output_dir": "./internal/model",
    "db_list": [
        {
            "pkg": "console",
            "link": "<dsn>",
            "db_name": "db",
            "table": [],
            "const": "DBConsole"
        },
        {
            "pkg": "business",
            "link": "<dsn>",
            "db_name": "db",
            "table": [],
            "const": "DBBusiness"
        },
    ]
}
```

运行生成工具：

```bash
./generate_model_darwin   # macOS
./generate_model_linux    # Linux
./generate_model.exe      # windows
```

## 延迟队列用法

```go
delayqueue.NewDelayQueue().Push(params)  // params []core.DelayQueuePushParams

type DelayQueuePushParams struct {
    Name      string `json:"name"`       // 任务名称
    ExecCmd   string `json:"exec_cmd"`   // 处理器名称（delay_queue/main.go 中注册）
    Params    string `json:"params"`     // 任务参数（字符串）
    DelayTime int64  `json:"delay_time"` // 延迟秒数
    ExecTime  int    `json:"exec_time"`  // 指定执行时间（优先级高于 DelayTime）
}
```

`delay_queue_log` 表状态值：`1` 待执行 / `2` 执行中 / `3` 已完成 / `4` 失败

---

## Skills —— AI 辅助开发技能包

`skills/` 目录包含针对本项目各业务模块的 AI 开发辅助规范，在 AI 编码时按需加载对应 skill，确保生成的代码符合项目约定。

### cron-script — 定时脚本

**适用场景：** 新增/修改/排查定时脚本，涉及 `crontab/`、`internal/script/`、`internal/handler/` 的调度、注册分发、执行日志链路。

关键约定：
- 定时任务通过 `crontab` 读取配置后启动子进程，再由 `internal/script` 分发，**不是直接调用业务函数**。
- 新增脚本须先实现 handler，再在 `internal/script/register.go` 注册脚本名。
- 脚本参数为字符串，handler 自行解析；返回值为 `*crontab.Result`。
- 日志链路关键值是 `uk`，需贯穿执行前后写入。

### exec-worker — 消费执行器

**适用场景：** 新增/修改/排查 Kafka 或 Redis 消费任务，涉及 `exec/` 目录下的消费者实现、ExecCmd 注册、配置驱动启动、扩容逻辑。

关键约定：
- 启动依赖**脚本配置表**，只在代码里注册 ExecCmd 还不够。
- 扩容行为由 `Progress`、`MaxProgress`、`ProgressLagLimit`、`ProgressAvgMsgcount` 等字段驱动。
- `exec/` 只放启动、调度、扩容逻辑，业务逻辑放 `exec/<domain>/` 或 `internal/`。

### delay-queue — 延迟队列

**适用场景：** 新增/修改/排查延迟投递任务，涉及 `delay_queue/`、`internal/delay_queue/`、`delay_queue_log` 状态流转和入队参数协议。

关键约定：
- 入队（`Push`）与执行（`Handler`）是两层职责，不要混写。
- `ExecCmd` 对应 `delay_queue/main.go` 的处理器注册，**不是** `internal/script/register.go` 的脚本名。
- `ExecTime` 优先级高于 `DelayTime`。
- `internal/delay_queue/mian.go`（拼写有误）是现有引用名，**不能随意重命名**。

### resident-task — 常驻后台任务

**适用场景：** 新增/修改/排查常驻后台任务，涉及 `resident/` 目录下的 Handler 实现、按机器分区注册、优雅退出。

关键约定：
- 常驻任务**只在 `queue` 模式下启动**。
- `InitEntry()` 按 `SCRIPT_PARTITION` 决定哪台机器运行哪些任务，注册须在正确分区的 `case` 下。
- `Handler()` 必须包含无限循环，并在每轮循环检查 `svc.KillSignal`。
- 优雅退出宽限期为 **8 秒**，超时强退；退出路径需保证在 8 秒内完成。
- 初始化逻辑（连接、预加载）放在循环之前，**不要放进 `InitEntry()`**。
