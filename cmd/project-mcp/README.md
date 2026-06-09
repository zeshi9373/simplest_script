# script project MCP

这是一个基于项目根目录 `README.md` 编写的本地 MCP server，用来把项目运行说明、目录职责和常见用法暴露成工具。

## 提供的工具

- `project_overview`
- `build_and_run_guide`
- `directory_reference`
- `delay_queue_reference`
- `model_generation_reference`
- `logging_reference`

## 启动
 
```bash
go run ./cmd/project-mcp
```

或先编译：

```bash
go build -o ./bin/project-mcp ./cmd/project-mcp
./bin/project-mcp
```

## 设计约束

- 只使用标准库，避免额外依赖
- 通过 stdio 提供 MCP 能力
- 输出内容以当前仓库 `README.md` 为主，并结合少量现有代码结构补充解释
- 不返回 README 中的明文数据库凭据
