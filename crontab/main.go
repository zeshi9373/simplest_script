package crontab

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"simplest_script/core"
	"simplest_script/core/conf"
	"simplest_script/core/tool"
	"simplest_script/internal/model/console"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron"
)

type CrontabExec struct {
	Cron    string `json:"cron"`
	Name    string `json:"name"`
	ExecCmd string `json:"exec_cmd"`
	Params  string `json:"params"`
	Status  int    `json:"status"`
	IsLog   int    `json:"is_log"`
}

type Result struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func Init() {
	// 创建 cron 实例（默认支持分钟级）
	c := cron.New()

	paritition := os.Getenv("SCRIPT_PARTITION")
	byteStream, err := os.ReadFile("./crontab_" + paritition + ".json")

	if err != nil {
		panic("read config file error")
	}

	crontabs := make([]CrontabExec, 0)
	if err := json.Unmarshal(byteStream, &crontabs); err != nil {
		panic("unmarshal config file error")
	}

	if len(crontabs) > 0 {
		go ZombieReaper()
	}

	for _, v := range crontabs {
		if core.StatusIsEnv(v.Status) {
			fmt.Printf("crontab script %v \n", v)
			c.AddFunc(v.Cron, func() {
				go execHandler(v.Name, v.ExecCmd, v.Params, v.IsLog)
			})
		}
	}

	// 启动 cron
	c.Start()
	defer c.Stop() // 程序退出时停止

	// 保持程序运行
	select {}
}

func execHandler(name string, execCmd string, params string, isLog int) {
	var uk string

	if isLog > 0 {
		uk = tool.Uuid()
		model := console.CrontabLog{
			Name:       name,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			Status:     "running",
			Result:     "",
			ExecCmd:    execCmd,
			Params:     params,
			Partition:  os.Getenv("SCRIPT_PARTITION"),
			Uk:         uk,
			StartTime:  int(time.Now().UnixMilli()),
		}
		console.NewCrontabLogModel().Create(&model)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		args := []string{"/B", conf.Conf.ExecCmd, uk, execCmd, params}
		cmd = exec.Command("cmd", "/c", "start")
		cmd.Args = append(cmd.Args, args...)
	} else {
		cmd = exec.Command("nohup", conf.Conf.ExecCmd, uk, execCmd, params)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	cmd.Process.Release()
}

func UpdateCrontabLog(uk string, status string, result any) {
	if len(uk) == 0 {
		return
	}

	log := console.CrontabLog{}

	console.NewCrontabLogModel().Where("uk = ? and start_time > ?", uk, time.Now().Add(-24*time.Hour).UnixMilli()).First(&log)

	if log.Id > 0 {
		endTime := time.Now().UnixMilli()
		costTime := endTime - int64(log.StartTime)
		rs, _ := json.Marshal(result)
		console.NewCrontabLogModel().Where("id = ?", log.Id).Updates(map[string]any{
			"pid":         os.Getpid(),
			"status":      status,
			"result":      string(rs),
			"update_time": time.Now(),
			"end_time":    int(time.Now().UnixMilli()),
			"cost_time":   int(costTime),
		})
	} else {
		lg := map[string]any{
			"uk":     uk,
			"status": status,
			"result": result,
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		}

		lgs, _ := json.Marshal(lg)
		hlog.Info("updateCrontabLogError: " + string(lgs))
	}
}
