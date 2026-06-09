package monitor

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"simplest_script/core/tool"
	"simplest_script/expand/feishu"
	"strings"
	"sync"
	"time"
)

var startTime = time.Now()

type EnhancedMonitor struct {
	stopChan chan struct{}
	wg       sync.WaitGroup
}

type MetricsSnapshot struct {
	Timestamp     time.Time
	Goroutines    int
	Sys           uint64
	HeapAlloc     uint64
	HeapSys       uint64
	HeapIdle      uint64
	HeapInuse     uint64
	StackInuse    uint64
	GCPercent     int
	NumGC         uint32
	PauseTotalNs  uint64
	LastGCPauseNs uint64
	CPUCount      int
	GOMAXPROCS    int
	NumCgoCall    int64
}

func NewEnhancedMonitor() *EnhancedMonitor {
	m := &EnhancedMonitor{
		stopChan: make(chan struct{}),
	}

	m.wg.Add(1)
	go m.startEnhancedMonitoring()

	return m
}

func (m *EnhancedMonitor) startEnhancedMonitoring() {
	defer m.wg.Done()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	m.collectAndReport()

	for {
		select {
		case <-ticker.C:
			m.collectAndReport()
		case <-m.stopChan:
			return
		}
	}
}

func (m *EnhancedMonitor) collectAndReport() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	snapshot := MetricsSnapshot{
		Timestamp:    time.Now(),
		Goroutines:   runtime.NumGoroutine(),
		HeapAlloc:    memStats.HeapAlloc,
		HeapSys:      memStats.HeapSys,
		Sys:          memStats.Sys,
		HeapIdle:     memStats.HeapIdle,
		HeapInuse:    memStats.HeapInuse,
		StackInuse:   memStats.StackInuse,
		NumGC:        memStats.NumGC,
		PauseTotalNs: memStats.PauseTotalNs,
		CPUCount:     runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
		NumCgoCall:   runtime.NumCgoCall(),
	}

	gcPercent := debug.SetGCPercent(-1)
	debug.SetGCPercent(gcPercent)
	snapshot.GCPercent = gcPercent

	if memStats.NumGC > 0 {
		snapshot.LastGCPauseNs = memStats.PauseNs[(memStats.NumGC+255)%256]
	}

	m.printReport(snapshot)
}

func (m *EnhancedMonitor) printReport(snapshot MetricsSnapshot) {
	monitorDetail := strings.Builder{}
	monitorDetail.WriteString(fmt.Sprintf("运行环境: %s\n", os.Getenv("SCRIPT_ENV")))
	monitorDetail.WriteString(fmt.Sprintf("运行机器: %s\n", os.Getenv("SCRIPT_PARTITION")))
	monitorDetail.WriteString(fmt.Sprintf("监控时间: %s\n", snapshot.Timestamp.Format(tool.DatetimeLayout)))
	monitorDetail.WriteString(fmt.Sprintf("运行时间: %v\n", time.Since(startTime)))
	monitorDetail.WriteString("\n1. 协程统计:\n")
	monitorDetail.WriteString(fmt.Sprintf("  当前协程数: %d\n", snapshot.Goroutines))
	monitorDetail.WriteString(fmt.Sprintf("  CPU核心数: %d\n  GOMAXPROCS: %d\n", snapshot.CPUCount, snapshot.GOMAXPROCS))
	monitorDetail.WriteString(fmt.Sprintf("  总内存: %.2f MB\n", float64(snapshot.Sys)/1024/1024))
	monitorDetail.WriteString("\n2. 内存使用:\n")
	monitorDetail.WriteString(fmt.Sprintf("  HeapSys(从OS申请的总堆内存): %.2f MB\n", float64(snapshot.HeapSys)/1024/1024))
	monitorDetail.WriteString(fmt.Sprintf("  HeapAlloc(活跃对象内存): %.2f MB\n", float64(snapshot.HeapAlloc)/1024/1024))
	monitorDetail.WriteString(fmt.Sprintf("  HeapInuse(使用中的堆区域): %.2f MB\n", float64(snapshot.HeapInuse)/1024/1024))
	monitorDetail.WriteString(fmt.Sprintf("  HeapIdle(空闲的堆区域): %.2f MB\n", float64(snapshot.HeapIdle)/1024/1024))
	monitorDetail.WriteString(fmt.Sprintf("  内存使用率: %.1f%%\n", float64(snapshot.HeapInuse)/float64(snapshot.HeapSys)*100))
	monitorDetail.WriteString("\n3. GC统计:\n")
	monitorDetail.WriteString(fmt.Sprintf("  GC次数: %d\n", snapshot.NumGC))

	if snapshot.LastGCPauseNs > 0 {
		monitorDetail.WriteString(fmt.Sprintf("  最近GC暂停: %.2f ms\n", float64(snapshot.LastGCPauseNs)/1000000))
	}

	monitorDetail.WriteString(fmt.Sprintf("  GC总暂停时间: %.2f ms\n", float64(snapshot.PauseTotalNs)/1000000))

	feishu.NewSendMessage("07ff68b9-5e9c-496a-b696-b803499d3e39").Send("text", monitorDetail.String())
}

func (m *EnhancedMonitor) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}
