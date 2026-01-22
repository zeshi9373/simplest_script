package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"simplest_script/core/conf"
	"strings"
	"sync"
	"time"
)

var loggerClient = make(map[string]*Logger)
var loggerMx sync.Mutex

// 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// 日志配置
type Config struct {
	LogDir      string   // 日志目录
	CustomDir   string   // 自定义日志目录
	FileName    string   // 日志文件名
	MaxFileSize int64    // 最大文件大小（字节）
	MaxBackups  int      // 最大备份文件数
	MaxAge      int      // 最大保留天数
	Compress    bool     // 是否压缩备份文件
	LogLevel    LogLevel // 日志级别
	Async       bool     // 是否异步写入
	BufferSize  int      // 缓冲区大小（异步时使用）
	JSONFormat  bool     // 是否使用JSON格式
	WithCaller  bool     // 是否记录调用者信息
}

// 日志条目
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Caller    string `json:"caller,omitempty"`
	File      string `json:"file,omitempty"`
	Line      int    `json:"line,omitempty"`
	Fields    Fields `json:"fields,omitempty"`
}

type Fields map[string]interface{}

// Logger 主结构
type Logger struct {
	config     *Config
	file       *os.File
	writer     *bufio.Writer
	mu         sync.RWMutex
	queue      chan *LogEntry
	stopCh     chan struct{}
	wg         sync.WaitGroup
	fields     Fields
	callerSkip int // 调用者跳过的层级
}

func NewLogger(customDir string) *Logger {
	config := &Config{
		LogDir:      conf.Conf.Logger.Path + "/" + customDir,
		CustomDir:   customDir,
		FileName:    "app-" + time.Now().Format("20060102") + ".log",
		MaxFileSize: 1024 * 1024 * 100, // 100MB
		MaxBackups:  10,
		MaxAge:      7,
		Compress:    true,
		LogLevel:    INFO,
		Async:       true,
		BufferSize:  1024,
		JSONFormat:  true,
	}

	return createLogger(config)
}

// createLogger 创建新的日志记录器
func createLogger(config *Config) *Logger {
	defer loggerMx.Unlock()
	loggerMx.Lock()

	if loggerClient[config.CustomDir] != nil {
		return loggerClient[config.CustomDir]
	}

	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil
	}

	l := &Logger{
		config:     config,
		queue:      make(chan *LogEntry, config.BufferSize),
		stopCh:     make(chan struct{}),
		fields:     make(Fields),
		callerSkip: 3, // 默认跳过3层调用栈
	}

	// 打开日志文件
	if err := l.openFile(); err != nil {
		return nil
	}

	// 启动异步写入协程
	if config.Async {
		l.wg.Add(1)
		go l.asyncWriter()
	}

	loggerClient[config.CustomDir] = l

	return l
}

// 打开日志文件
func (l *Logger) openFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 关闭旧文件
	if l.file != nil {
		l.writer.Flush()
		l.file.Close()
	}

	filePath := filepath.Join(l.config.LogDir, l.config.FileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	l.file = file
	l.writer = bufio.NewWriterSize(file, 32*1024) // 32KB缓冲区
	return nil
}

// 检查并轮转日志文件
func (l *Logger) rotateIfNeeded() error {
	info, err := l.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() < l.config.MaxFileSize {
		return nil
	}

	// 轮转文件
	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(l.config.LogDir, l.config.FileName)
	newPath := filepath.Join(l.config.LogDir,
		fmt.Sprintf("%s.%s", l.config.FileName, timestamp))

	// 关闭当前文件
	l.writer.Flush()
	l.file.Close()

	// 重命名文件
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	// 重新打开文件
	return l.openFile()
}

// 异步写入协程
func (l *Logger) asyncWriter() {
	// 关闭清空队列
	defer func() {
		for {
			select {
			case entry := <-l.queue:
				l.writeEntry(entry)
			default:
				l.writer.Flush()
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case entry := <-l.queue:
			l.writeEntry(entry)
		case <-ticker.C:
			// 清空队列
			for {
				select {
				case entry := <-l.queue:
					l.writeEntry(entry)
				default:
					l.writer.Flush()
					return
				}
			}
		}
	}
}

// 写入日志条目
func (l *Logger) writeEntry(entry *LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查是否需要轮转
	if err := l.rotateIfNeeded(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}

	var logLine string
	if l.config.JSONFormat {
		data, _ := json.Marshal(entry)
		logLine = string(data) + "\n"
	} else {
		callerInfo := ""
		if entry.Caller != "" {
			callerInfo = fmt.Sprintf(" [%s]", entry.Caller)
		}
		fieldsStr := ""
		if len(entry.Fields) > 0 {
			fieldsStr = " " + l.formatFields(entry.Fields)
		}
		logLine = fmt.Sprintf("%s [%s]%s %s%s\n",
			entry.Timestamp, entry.Level, callerInfo, entry.Message, fieldsStr)
	}

	if _, err := l.writer.WriteString(logLine); err != nil {
		fmt.Fprintf(os.Stderr, "写入日志失败: %v\n", err)
	}
}

// 格式化字段
func (l *Logger) formatFields(fields Fields) string {
	var parts []string
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, " ")
}

// 记录日志
func (l *Logger) log(level LogLevel, msg string, fields Fields) {
	if level < l.config.LogLevel {
		return
	}

	entry := &LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Level:     levelNames[level],
		Message:   msg,
		Fields:    l.mergeFields(fields),
	}

	// 添加调用者信息
	if l.config.WithCaller {
		if pc, file, line, ok := runtime.Caller(l.callerSkip); ok {
			entry.File = filepath.Base(file)
			entry.Line = line
			funcName := runtime.FuncForPC(pc).Name()
			entry.Caller = fmt.Sprintf("%s:%d", funcName, line)
		}
	}

	if l.config.Async {
		select {
		case l.queue <- entry:
		default:
			// 缓冲区满，同步写入
			l.writeEntry(entry)
		}
	} else {
		l.writeEntry(entry)
	}
}

// 合并字段
func (l *Logger) mergeFields(fields Fields) Fields {
	if len(l.fields) == 0 && len(fields) == 0 {
		return nil
	}

	merged := make(Fields)
	for k, v := range l.fields {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}
	return merged
}

// 公共日志方法
func (l *Logger) Debug(msg string, fields ...Fields) {
	l.log(DEBUG, msg, l.getFields(fields))
}

func (l *Logger) Info(msg string, fields ...Fields) {
	l.log(INFO, msg, l.getFields(fields))
}

func (l *Logger) Warn(msg string, fields ...Fields) {
	l.log(WARN, msg, l.getFields(fields))
}

func (l *Logger) Error(msg string, fields ...Fields) {
	l.log(ERROR, msg, l.getFields(fields))
}

func (l *Logger) Fatal(msg string, fields ...Fields) {
	l.log(FATAL, msg, l.getFields(fields))
	os.Exit(1)
}

func (l *Logger) getFields(fields []Fields) Fields {
	if len(fields) > 0 {
		return fields[0]
	}
	return nil
}

// 设置字段
func (l *Logger) WithFields(fields Fields) *Logger {
	l.fields = l.mergeFields(fields)
	return l
}

// 设置调用者跳过层级
func (l *Logger) WithCallerSkip(skip int) *Logger {
	l.callerSkip = skip
	return l
}

// 同步刷新
func (l *Logger) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.writer.Flush()
}

// 关闭日志记录器
func (l *Logger) Close() error {
	if l.config.Async {
		close(l.stopCh)
		l.wg.Wait()
	}
	return l.Sync()
}
