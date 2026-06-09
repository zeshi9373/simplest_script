# 不推荐写法

## 1. 把脚本执行逻辑写回 `crontab/main.go`

错误方向：

```go
// crontab/main.go — 错误示范
func execHandler(script ScriptConfig) {
    // 直接在调度层实现业务逻辑
    if script.ExecCmd == "sync_data" {
        db.Exec("UPDATE ...")
        redis.Del(...)
    }
}
```

正确方向：

- 调度层只负责读取配置、启动子进程、记录日志
- 业务逻辑通过子进程 → `internal/script.Exec()` → 具体 handler 执行

原因：

- 调度层混入业务后，无法单独测试业务逻辑
- 每次新增脚本都要改调度层，风险范围扩大

## 2. 只新增 handler，不做注册

错误方向：

```go
// internal/handler/sync_data.go — 写了，但没注册
type SyncData struct{}
func (h *SyncData) Handler(params string) *crontab.Result { ... }

// internal/script/register.go — 没有添加这行
// HandlerEntry["sync_data"] = &syncdata.SyncData{}  // 漏了！
```

原因：

- 项目靠注册表分发，`internal/script/register.go` 里没有就永远不会执行
- 不会有任何报错，只是静默地不执行

## 3. 破坏子进程参数协议

错误方向：

```go
// 调度层原来按顺序传参
cmd := exec.Command(binPath, uk, execCmd, params)

// 误改成命名参数或调换顺序
cmd := exec.Command(binPath, "--uk="+uk, "--cmd="+execCmd)
```

原因：

- 主程序非 `queue` 分支按位置解析参数：`os.Args[1]` = uk，`os.Args[2]` = exec_cmd，`os.Args[3]` = params
- 改顺序或格式后，已在配置表里的所有脚本全部失效，且不报错

## 4. 日志链路不完整（只有 running，没有结束状态）

错误方向：

```go
func (h *Demo) Handler(params string) *crontab.Result {
    doSomething()
    // 没有返回值，或 panic 被吞掉
}
```

原因：

- 主程序在 handler 执行后调用 `crontab.UpdateCrontabLog()` 写结束状态
- 如果 handler panic 或返回 nil，日志表里会一直停在 running
- 正确做法：任何异常路径都要 recover 并返回错误 Result

## 5. 在 handler 里直接使用裸字符串参数

错误方向：

```go
func (h *Demo) Handler(params string) *crontab.Result {
    // 直接把 params 当 JSON key 用，没有解析
    if params == "mode=full" { ... }
}
```

正确方向：

```go
type DemoParams struct {
    Mode string `json:"mode"`
}

func (h *Demo) Handler(params string) *crontab.Result {
    var p DemoParams
    if err := json.Unmarshal([]byte(params), &p); err != nil {
        return &crontab.Result{Status: 1, Data: "invalid params: " + err.Error()}
    }
    // 使用 p.Mode
}
```
