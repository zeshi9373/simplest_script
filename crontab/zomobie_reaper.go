//go:build !windows
// +build !windows

package crontab

import (
	"syscall"
	"time"
)

// 定期清理僵尸进程
func ZombieReaper() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 非阻塞方式等待任意子进程
		var wstatus syscall.WaitStatus
		_, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
		if err != nil {
			// 没有子进程或错误，忽略
		}
	}
}
