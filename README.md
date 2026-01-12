# 项目介绍
simplest_script是一个集成定时任务、异步队列消费及延时队列消费于一体的分布式脚本go项目，集成了一些简易的基础工具，可以快速上手集成到商业项目中

# 设计思想及目的
目前市面上没有一个go项目集成异步队列消费、定时任务及延迟队列消费于一体的分布式脚本项目，此项目就是解决这种需求目的

定时任务设计思想：定时向系统发送一个后台运营的脚本，记录这个脚本运行的记录（crontab_log）数据表内，这种设计即使代码更新发版都不影响正在执行中的任务

异步队列消费设计思想：项目目前支持kafka及redis list队列异步消费，可以配置初始消费者数量（progress），最大消费者数量（max_progress），会有一个定时器（每3分钟）去检查消息堆积的数量，根据配置的消息堆积的阈值（progress_lag_limit）及每个消费者平均处理的消息数（progress_avg_msgcount）去增加消费者数量，但是不会超过设置的最大消费者数量

延时队列设计思想：每次创建一个延时队列处理任务都会记录在delay_queue_log表中，每分钟有一个定时任务去捞取未来60秒需要执行的数据，并行开协程去处理（协程中会有延时处理，根据执行时间exec_time判断延迟多久）

# 正式环境需要设置

```
export SCRIPT_ENV=release # 环境 dev test release
export SCRIPT_PARTITION=script1 # 脚本所在分区 默认为 script1 后续机器编号 递增script2 script3 ...
```

## 目录结构
```
   /core  核心非业务公用库
    ~~/conf  解析配置文件
	~~/logger 日志记录自定义文件夹包
	~~/svc 服务中间件注册
	~~/tool 工具包（非业务）
	~~/warning 预警封装基础包
	/crontab 定时任务注册（非业务）
	/etc 配置文件
	/exec 异步消费业务（包括kafka及redis队列，配置文件queue_script1.json，queue_script2.json）
	/expand 外部业务接口
	/internal 内部业务处理
	~~/consts  常量配置
	~~/delay_queue 延迟队列业务
	~~/handler  定时任务业务
	~~/model  数据库表
	~~/script 定时任务方法配置
	~~/services 业务逻辑
	~~/types 结构体
	/logs 日志目录
	/test 测试目录
```
## 各类用法说明

####异步消费者配置 queue_script1.json  
```
[
    {
        "name": "kafka测试",     // 名称
        "exec_cmd": "kafka_test",   // 脚本名称  在文件夹exec文件夹内
        "topic": "kafka_test_topic",  // topic , kafka 必填
        "group_id": "kafkaTest", // 消费者, kafka 必填
        "progress": 2,  // 默认协程数
        "max_progress": 10,  // 最大协程数，kafka默认分区数
        "progress_lag_limit": 20000,  // 消息堆积阈值
        "progress_avg_msgcount": 10000, // 协程平均消息数，超时5分钟
        "status": 4  // 位运算 1: dev 2: test 4: release
    },
    {
        "name": "redis测试",     // 名称
        "exec_cmd": "redis_test",   // 脚本名称  在文件夹exec内
        "key": "redis_test",  // redis list key, redis 必填
        "progress": 2,  // 默认协程数
        "max_progress": 10,  // 最大协程数
        "progress_lag_limit": 5000,  // 消息堆积阈值
        "progress_avg_msgcount": 1000, // 协程平均消息数，超时5分钟
        "status": 4  // 位运算 1: dev 2: test 4: release
    }
]
```

### 定时脚本配置 crontab_script1.json
```
[
    {
        "cron": "0 * * * * *",     // 执行配置 秒 分钟 小时 日 月 周
        "name": "测试定时任务",  // 名称
        "exec_cmd": "crontab_track",  // 命令  在internal/handler文件夹内
        "params": "{\"name\":\"测试\"}",  // 参数
        "status": 4,  // 位运算 1: dev 2: test 4: release
        "is_log": 0  // 是否记录日志 1: 记录 0: 不记录
    }
]
```

脚本执行命令可以在日志文件里面查看 logs/*.log


### 生成数据库model文件配置
```
generateModel.json
{
    "module_name":"simplest_script",  // 项目模块名称
    "output_dir":"./internal/model", // 输出目录
    "db_list":[
        {
            "pkg":"console", // 包名
            "link":"root:123456@tcp(127.0.0.1:3306)/console?charset=utf8mb4&parseTime=True&loc=Local", // 连接
            "db_name":"console", // 数据库名
            "table":[], // 指定表名，为空表示所有表
            "const":"DBConsole" // 连接常量
        }
    ]
}
```

### 延时队列适用方法
```
delayqueue.NewDelayQueue().Push(params)   // params []core.DelayQueuePushParams

type DelayQueuePushParams struct {
	Name      string `json:"name"` // 名称
	ExecCmd   string `json:"exec_cmd"` // 执行命令方法 在internal/delay_queue文件夹内
	Params    string `json:"params"`  // 参数
	DelayTime int64  `json:"delay_time"` // 延迟秒数
	ExecTime  int    `json:"exec_time"`  // 执行时间，有设置就忽略DelayTime
}
```

### 日志记录
```
logger.NewLogger("testLog").Info("test log")  // testLog是项目logs目录下的子目录，由于日志写入是异步每秒钟刷入硬盘，如果程序执行时间短，建议执行末尾加上一个time.Sleep()
```