package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

)

var (
	errorLog *log.Logger
	debugLog *log.Logger
	infoLog  *log.Logger
	fatalLog *log.Logger // 添加fatalLog变量
	loggers  = []*log.Logger{errorLog, debugLog}
	mu       sync.Mutex
)

// log levels
const (
	DebugLevel = iota
	InfoLevel
	ErrorLevel
	Disabled
)

const (
	colorRed    = "\033[31m" // 红色
	colorGreen  = "\033[32m" // 绿色
	colorYellow = "\033[33m" // 黄色
	colorBlue   = "\033[34m" // 蓝色
	colorReset  = "\033[0m"  // 重置颜色
)

func InitLogger(logFilePath string) {
	var logOutput io.Writer
	if logFilePath == "" {
		logOutput = os.Stdout
	} else {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logOutput = logFile
	}

	errorLog = log.New(logOutput, colorRed+"[error]"+colorReset, log.LstdFlags|log.Lshortfile)
	debugLog = log.New(logOutput, colorGreen+"[debug]"+colorReset, log.LstdFlags|log.Lshortfile)
	infoLog = log.New(logOutput, colorBlue+"[info]"+colorReset, log.LstdFlags|log.Lshortfile)
	fatalLog = log.New(logOutput, colorRed+"[fatal]"+colorReset, log.LstdFlags|log.Lshortfile) // 初始化fatalLog
}
// 文件写入后续优化可以改为 mmap 写入  不用 write file  效率会更高

func Error(v ...any) {
	logWithPosition(errorLog, fmt.Sprint(v...), colorRed)
}

func Errorf(format string, v ...any) {
	logWithPosition(errorLog, fmt.Sprintf(format, v...), colorRed)
}

func Fatal(v ...any) {
	logWithPosition(fatalLog, fmt.Sprint(v...), colorYellow)
	os.Exit(1)
}

// Fatalf is equivalent to [Printf] followed by a call to [os.Exit](1).
func Fatalf(format string, v ...any) {
	logWithPosition(fatalLog, fmt.Sprintf(format, v...), colorYellow)
	os.Exit(1)
}

func Debug(v ...any) {
	logWithPosition(debugLog, fmt.Sprint(v...), colorGreen)
}

func Debugf(format string, v ...any) {
	logWithPosition(debugLog, fmt.Sprintf(format, v...), colorGreen)
}

func Info(v ...any) {
	logWithPosition(infoLog, fmt.Sprint(v...), colorBlue)
}

func Infof(format string, v ...any) {
	logWithPosition(infoLog, fmt.Sprintf(format, v...), colorBlue)
}

func TestColor(v ...any) {
	logWithPosition(debugLog, fmt.Sprint(v...), "33")
}

// logWithPosition 记录一条带有文件名和行号的日志信息
func logWithPosition(l *log.Logger, msg, colorCode string) {
	// 获取调用者的文件名和行号
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		// 获取函数名
		// fn := runtime.FuncForPC(pc).Name()
		// 提取文件名和行号
		file = file[strings.LastIndex(file, "/")+1:]
	}

	// 格式化日志信息
	message := fmt.Sprintf("%s [%s:%d] %s \033[0m", colorCode, file, line, msg)
	fmt.Println(message)
	if l != nil {
		l.Println(message)
	}
}

/*
灵感来源
package mylogger

import (
	"log"
	"os"
	"runtime"
)

type MyLogger struct {
	*log.Logger
}

// New 创建一个新的 MyLogger 实例
func New(out *os.File, prefix string, flag int) *MyLogger {
	logger := log.New(out, prefix, flag)
	return &MyLogger{logger}
}

// Trace 记录一条 trace 级别的日志信息
func (l *MyLogger) Trace(v ...interface{}) {
	l.logWithPosition("TRACE", v...)
}

// Debug 记录一条 debug 级别的日志信息
func (l *MyLogger) Debug(v ...interface{}) {
	l.logWithPosition("DEBUG", v...)
}

// Info 记录一条 info 级别的日志信息
func (l *MyLogger) Info(v ...interface{}) {
	l.logWithPosition("INFO", v...)
}

// Warn 记录一条 warn 级别的日志信息
func (l *MyLogger) Warn(v ...interface{}) {
	l.logWithPosition("WARN", v...)
}

// Error 记录一条 error 级别的日志信息
func (l *MyLogger) Error(v ...interface{}) {
	l.logWithPosition("ERROR", v...)
}

// logWithPosition 记录一条带有文件名和行号的日志信息
func (l *MyLogger) logWithPosition(level string, v ...interface{}) {
	// 获取调用者的文件名和行号
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		// 获取函数名
		fn := runtime.FuncForPC(pc).Name()
		// 提取文件名和行号
		file = fn[strings.LastIndex(fn, "/")+1:]
	}

	// 格式化日志信息
	message := fmt.Sprintf("[%s] [%s:%d] ", level, file, line)
	for _, arg := range v {
		message += fmt.Sprintf("%v ", arg)
	}
	message += "\n"

	// 输出日志信息
	l.Logger.Println(message)
}

func logWithPosition(msg string) {
	// 获取当前 goroutine 的栈跟踪信息
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	stackTrace := string(stack[:length])

	// 查找栈跟踪信息中的文件名和行号
	var fileName string
	var lineNumber int
	_, _ = fmt.Sscanf(stackTrace, "main.go:%d", &lineNumber) // 假设栈顶是 main.go 文件
	for _, line := range strings.Split(stackTrace, "\n") {
		if strings.Contains(line, "logWithPosition") {
			fileName = strings.TrimSpace(strings.Split(line, "(")[0])
			break
		}
	}

	// 打印日志信息和位置
	log.Printf("%s %s:%d %s", msg, fileName, lineNumber, debug.Stack())
}
*/
