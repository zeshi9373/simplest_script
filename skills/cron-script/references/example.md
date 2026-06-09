# 最小样例

## 新增一个脚本 handler

### 步骤

1. 在 `internal/handler/` 或合适的业务目录实现 handler
2. 返回值保持 `*crontab.Result`，任何路径都不能返回 nil
3. 在 `internal/script/register.go` 注册脚本名

### 实现骨架

```go
// internal/handler/demo/demo.go
package demo

import (
    "simplest_script/crontab"
    "encoding/json"
)

type Demo struct{}

type demoParams struct {
    Mode string `json:"mode"`
}

func (h *Demo) Handler(params string) *crontab.Result {
    var p demoParams
    if err := json.Unmarshal([]byte(params), &p); err != nil {
        return &crontab.Result{Status: 1, Data: "invalid params: " + err.Error()}
    }

    // 业务逻辑
    result, err := doWork(p.Mode)
    if err != nil {
        return &crontab.Result{Status: 1, Data: err.Error()}
    }

    return &crontab.Result{Status: 0, Data: result}
}
```

### 注册

```go
// internal/script/register.go 的 InitEntry() 里添加：
HandlerEntry["simplest_script"] = &demo.Demo{}
```

### 配置表要同步

配置表里需要有对应记录（`exec_cmd = "simplest_script"`）才会被调度到，代码注册和配置表缺一不可。

---

## 完整链路验证清单

新增或修改脚本后，按顺序检查：

- [ ] handler 在 `internal/script/register.go` 已注册
- [ ] 脚本名（`exec_cmd`）与配置表记录匹配
- [ ] handler 返回值类型是 `*crontab.Result`，所有路径都有返回值
- [ ] 参数是字符串，handler 内部有解析和错误处理
- [ ] 如果开启了日志，日志表的 running → success/failed 状态流转完整
- [ ] `gofmt -w` 格式化
- [ ] `go build .` 编译通过
