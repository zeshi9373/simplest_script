package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type initializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
}

type toolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type textContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		payload, err := readMessage(reader)
		if err != nil {
			if err == io.EOF {
				return
			}

			writeResponse(writer, rpcResponse{
				JSONRPC: "2.0",
				Error: &rpcError{
					Code:    -32700,
					Message: fmt.Sprintf("parse error: %v", err),
				},
			})
			return
		}

		var req rpcRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			writeResponse(writer, rpcResponse{
				JSONRPC: "2.0",
				Error: &rpcError{
					Code:    -32700,
					Message: fmt.Sprintf("invalid request: %v", err),
				},
			})
			continue
		}

		if len(req.ID) == 0 {
			handleNotification(req)
			continue
		}

		resp := handleRequest(req)
		writeResponse(writer, resp)
	}
}

func handleNotification(req rpcRequest) {
	switch req.Method {
	case "notifications/initialized":
		return
	default:
		return
	}
}

func handleRequest(req rpcRequest) rpcResponse {
	id := decodeID(req.ID)

	switch req.Method {
	case "initialize":
		var params initializeParams
		_ = json.Unmarshal(req.Params, &params)
		version := params.ProtocolVersion
		if version == "" {
			version = "2024-11-05"
		}

		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"protocolVersion": version,
				"capabilities": map[string]any{
					"tools": map[string]any{
						"listChanged": false,
					},
				},
				"serverInfo": map[string]any{
					"name":    "script-project-mcp",
					"version": "0.1.0",
				},
			},
		}
	case "tools/list":
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"tools": projectTools(),
			},
		}
	case "tools/call":
		var params toolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return errorResponse(id, -32602, fmt.Sprintf("invalid tool call params: %v", err))
		}

		text, err := callTool(params.Name, params.Arguments)
		if err != nil {
			return errorResponse(id, -32602, err.Error())
		}

		return rpcResponse{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"content": []textContent{
					{
						Type: "text",
						Text: text,
					},
				},
			},
		}
	default:
		return errorResponse(id, -32601, "method not found")
	}
}

func errorResponse(id any, code int, message string) rpcResponse {
	return rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &rpcError{
			Code:    code,
			Message: message,
		},
	}
}

func decodeID(raw json.RawMessage) any {
	var id any
	if err := json.Unmarshal(raw, &id); err != nil {
		return string(raw)
	}
	return id
}

func readMessage(reader *bufio.Reader) ([]byte, error) {
	contentLength := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		headerName := strings.TrimSpace(strings.ToLower(parts[0]))
		headerValue := strings.TrimSpace(parts[1])
		if headerName == "content-length" {
			_, err := fmt.Sscanf(headerValue, "%d", &contentLength)
			if err != nil {
				return nil, fmt.Errorf("invalid content-length: %w", err)
			}
		}
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("missing content-length")
	}

	payload := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func writeResponse(writer *bufio.Writer, resp rpcResponse) {
	body, err := json.Marshal(resp)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	buffer.WriteString("Content-Type: application/json\r\n\r\n")
	buffer.Write(body)

	_, _ = writer.Write(buffer.Bytes())
	_ = writer.Flush()
}

func projectTools() []tool {
	return []tool{
		{
			Name:        "project_overview",
			Description: "Summarize the CPA script project structure, runtime modes, and core business flow from the local README.",
			InputSchema: objectSchema(
				map[string]any{
					"focus": map[string]any{
						"type":        "string",
						"description": "Optional focus area: runtime, structure, exec, cron, or delay_queue.",
					},
				},
			),
		},
		{
			Name:        "build_and_run_guide",
			Description: "Generate the environment variables, build command, and start command for a selected runtime mode.",
			InputSchema: objectSchema(
				map[string]any{
					"mode": map[string]any{
						"type":        "string",
						"enum":        []string{"queue", "cron"},
						"description": "Runtime mode. queue starts consumers/resident tasks; cron runs a single script entry.",
					},
					"env": map[string]any{
						"type":        "string",
						"enum":        []string{"dev", "test", "release"},
						"description": "SCRIPT_ENV value.",
					},
					"partition": map[string]any{
						"type":        "string",
						"description": "SCRIPT_PARTITION value such as script1.",
					},
				},
				"mode",
			),
		},
		{
			Name:        "directory_reference",
			Description: "Explain the role of a top-level directory from the README.",
			InputSchema: objectSchema(
				map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Directory name, for example core, exec, internal, resident, crontab, or test.",
					},
				},
				"name",
			),
		},
		{
			Name:        "delay_queue_reference",
			Description: "Return the documented delay queue API shape and explain each field.",
			InputSchema: objectSchema(nil),
		},
		{
			Name:        "model_generation_reference",
			Description: "Explain the generateModel.json structure from the README without exposing concrete database credentials.",
			InputSchema: objectSchema(nil),
		},
		{
			Name:        "logging_reference",
			Description: "Explain where this project writes logs and how script execution records are observed.",
			InputSchema: objectSchema(nil),
		},
	}
}

func callTool(name string, args map[string]interface{}) (string, error) {
	switch name {
	case "project_overview":
		return projectOverview(stringArg(args, "focus")), nil
	case "build_and_run_guide":
		mode := stringArg(args, "mode")
		if mode == "" {
			return "", fmt.Errorf("mode is required")
		}
		return buildAndRunGuide(mode, stringArg(args, "env"), stringArg(args, "partition")), nil
	case "directory_reference":
		dir := stringArg(args, "name")
		if dir == "" {
			return "", fmt.Errorf("name is required")
		}
		return directoryReference(dir)
	case "delay_queue_reference":
		return delayQueueReference(), nil
	case "model_generation_reference":
		return modelGenerationReference(), nil
	case "logging_reference":
		return loggingReference(), nil
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func projectOverview(focus string) string {
	base := strings.TrimSpace(`
这是一个 Go 编写的 CPA 脚本项目，核心职责是运行三类任务：

1. 定时脚本任务
2. Kafka / Redis 队列消费任务
3. 不依赖外部触发器的常驻任务

运行基础约束：
- 通过 SCRIPT_ENV 指定环境，常见值为 dev、test、release
- 通过 SCRIPT_PARTITION 指定脚本分区，例如 script1、script2
- 正式环境按 README 约定先 build，再以 mode 启动

顶层结构重点：
- core: 通用基础设施、配置、日志、中间件、工具和预警
- crontab: 定时任务注册与子进程调度
- exec: Kafka / Redis 异步消费执行框架
- internal: 业务实现、模型、类型和主流程封装
- resident: 常驻后台任务
- test: 测试目录
`)

	switch focus {
	case "runtime":
		return strings.TrimSpace(`
运行模式分两类：

- queue: 启动异步消费、常驻任务、程序监控，并注册定时任务
- cron: 执行单个脚本方法，通常用于被 crontab 子进程拉起

README 明确的启动方式：
1. 设置 SCRIPT_ENV 和 SCRIPT_PARTITION
2. 执行 go build .
3. 运行 ./simplest_script -mode queue
`)
	case "structure":
		return base
	case "exec":
		return strings.TrimSpace(`
exec 目录负责异步消费业务，包含 Kafka 和 Redis 队列消费者。README 说明它依赖 mysql 表 script_config 作为配置来源。

在当前代码里，exec 的职责包括：
- 加载脚本消费配置
- 按 topic 或 redis key 启动消费者
- 检查堆积并动态补充消费者
`)
	case "cron":
		return strings.TrimSpace(`
crontab 目录负责定时任务注册。当前代码会从配置表读取 type = 2 的启用脚本，按 cron 表达式调度，并通过子进程再次调用主程序执行脚本方法。

这意味着 cron 模式本质上是“单次脚本执行器”，queue 模式才是长期运行模式。
`)
	case "delay_queue":
		return delayQueueReference()
	default:
		return base
	}
}

func buildAndRunGuide(mode, env, partition string) string {
	if env == "" {
		env = "dev"
	}
	if partition == "" {
		partition = "script1"
	}

	configPath := "./etc/dev.yaml"
	switch env {
	case "release":
		configPath = "/data/simplest_script/etc/release.yaml"
	case "test":
		configPath = "/data/etc/test.yaml"
	case "dev":
		configPath = "./etc/dev.yaml"
	}

	if mode == "queue" {
		return fmt.Sprintf(strings.TrimSpace(`
推荐启动参数：

export SCRIPT_ENV=%s
export SCRIPT_PARTITION=%s

go build .
./simplest_script -mode queue -f %s

queue 模式会启动：
- 定时任务注册
- Kafka / Redis 消费者
- 常驻任务
- 运行期监控
`), env, partition, configPath)
	}

	return fmt.Sprintf(strings.TrimSpace(`
cron 模式用于单次执行脚本方法。

export SCRIPT_ENV=%s
export SCRIPT_PARTITION=%s

go build .
./simplest_script <uk> <exec_cmd> <params>

注意：
- 代码里的默认 mode 是 cron 分支，也就是非 queue 时执行单次脚本
- 常规情况下，这个模式由 crontab 子进程触发，不是人工长期驻留运行
- 配置文件默认按 SCRIPT_ENV 推导，dev 为 %s
`), env, partition, configPath)
}

func directoryReference(name string) (string, error) {
	refs := map[string]string{
		"core": strings.TrimSpace(`
core 是核心非业务公用库，README 列出的子模块包括：
- conf: 配置解析
- logger: 日志封装
- permission: 数据权限
- svc: 服务中间件注册
- tool: 非业务工具包
- warning: 预警基础封装
`),
		"crontab": strings.TrimSpace(`
crontab 负责定时任务注册与执行调度。当前项目中它会从配置表读取 cron 脚本，写执行日志，并用子进程方式拉起单次脚本执行。
`),
		"exec": strings.TrimSpace(`
exec 负责异步消费业务，包括 Kafka 和 Redis 队列。README 指出它的配置来源是 mysql 表 script_config。
`),
		"expand": strings.TrimSpace(`
expand 负责外部业务接口集成，适合放短信、邮件等第三方能力。
`),
		"internal": strings.TrimSpace(`
internal 是内部业务处理目录，README 列出的子模块包括：
- consts: 常量配置
- delay_queue: 延迟队列业务
- handler: 定时任务业务
- main_progress: 业务主流程封装
- model: 数据库表模型
- script: 定时任务方法配置
- services: 业务逻辑
- types: 结构体
`),
		"resident": strings.TrimSpace(`
resident 放常驻进程任务，README 定义为“不需要外部信息或者触发器”的后台任务。
`),
		"test": strings.TrimSpace(`
test 是测试目录。README 只给出了目录定位，没有进一步约束。
`),
		"etc": strings.TrimSpace(`
etc 存放配置文件。主程序会根据 SCRIPT_ENV 选择 dev、test、release 对应的配置路径。
`),
		"logs": strings.TrimSpace(`
logs 是日志目录。README 说明脚本执行命令可以在 logs/*.log 中查看。
`),
	}

	value, ok := refs[strings.ToLower(name)]
	if !ok {
		return "", fmt.Errorf("unknown directory: %s", name)
	}

	return value, nil
}

func delayQueueReference() string {
	return strings.TrimSpace(`
README 中记录的延迟队列调用方式：

delayqueue.NewDelayQueue().Push(params)

其中 params 类型为 []core.DelayQueuePushParams，字段含义如下：
- name: 任务名称
- exec_cmd: 执行方法名，对应 internal/delay_queue 目录中的处理器
- params: 任务参数字符串
- delay_time: 延迟秒数
- exec_time: 指定执行时间；如果设置了这个字段，会覆盖 delay_time

适用理解：
- 这是“把业务任务延后执行”的统一入口
- 任务实际执行器按 exec_cmd 分发
- 延迟和定时是互斥的，以 exec_time 优先
`)
}

func modelGenerationReference() string {
	return strings.TrimSpace(`
README 中的 generateModel.json 用于生成数据库 model 文件，结构包含：

- module_name: Go 模块名，例如 simplest_script
- output_dir: 生成输出目录，例如 ./internal/model
- db_list: 数据库连接配置数组

db_list 每项的关键字段：
- pkg: 生成后的包名
- link: 数据库连接串
- db_name: 数据库名
- table: 指定表名列表；为空表示全表
- const: 连接常量名

安全建议：
- README 里的 link 示例不应该继续直接复用到新代码或新文档
- 实际使用时应改成环境变量、密钥管理或本地未提交配置
- 这个 MCP 只解释结构，不返回任何明文凭据
`)
}

func loggingReference() string {
	return strings.TrimSpace(`
README 和当前代码反映出的日志方式：

- 常规运行日志写入配置中的 logger.path，对应按日期切分的日志文件
- README 明确提示脚本执行命令可以在 logs/*.log 查看
- crontab 子进程执行时，还会把执行状态和结果写入数据库日志表

排查脚本问题时建议先看两处：
1. 文件日志目录，例如 logs/20260604.log
2. 定时任务执行日志记录表
`)
}

func objectSchema(properties map[string]any, required ...string) map[string]interface{} {
	if properties == nil {
		properties = map[string]any{}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func stringArg(args map[string]interface{}, key string) string {
	if args == nil {
		return ""
	}

	value, ok := args[key]
	if !ok || value == nil {
		return ""
	}

	text, ok := value.(string)
	if !ok {
		return ""
	}

	return strings.TrimSpace(text)
}
