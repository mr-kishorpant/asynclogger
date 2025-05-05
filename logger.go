package async

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type logEntry struct {
	Level   LogLevel
	Message string
	Time    time.Time
}
type Logger struct {
	out       io.Writer
	logChan   chan logEntry
	wg        sync.WaitGroup
	quit      chan struct{}
	onceClose sync.Once
}

// [LogLevel, Logger struct, constants remain unchanged...]

var (
	defaultLogger  *Logger
	onceInit       sync.Once
	defaultBufSize = 100
)

func GetDateBasedLogFile() (*os.File, error) {
	dir := "storage/logs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	filename := filepath.Join(dir, "log-"+time.Now().Format("2006-01-02")+".log")
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func InitWithBufferSize(bufferSize int) error {
	file, err := GetDateBasedLogFile()
	if err != nil {
		return err
	}

	logger := &Logger{
		out:     file,
		logChan: make(chan logEntry, bufferSize),
		quit:    make(chan struct{}),
	}
	defaultLogger = logger
	logger.start()
	listenForShutdown() // <== Hook into OS signals
	return nil
}

func Init() error {
	return InitWithBufferSize(defaultBufSize)
}

func (l *Logger) start() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		encoder := json.NewEncoder(l.out) // stream encoder (writes each log on its own line)
		for {
			select {
			case entry := <-l.logChan:
				encoder.Encode(entry) // one JSON object per line
			case <-l.quit:
				for len(l.logChan) > 0 {
					entry := <-l.logChan
					encoder.Encode(entry)
				}
				return
			}
		}
	}()
}

func Shutdown() {
	if defaultLogger != nil {
		defaultLogger.onceClose.Do(func() {
			close(defaultLogger.quit)
			defaultLogger.wg.Wait()
		})
	}
}

func ensureInit() {
	onceInit.Do(func() {
		if err := Init(); err != nil {
			fmt.Println("Logger failed to initialize:", err)
		}
	})
}

func log(level LogLevel, msg string, args ...interface{}) {
	ensureInit()
	if defaultLogger == nil {
		fmt.Println("Logger is not initialized")
		return
	}

	// Concatenate args with the message
	var sb strings.Builder
	sb.WriteString(msg)
	for _, arg := range args {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(arg))
	}

	entry := logEntry{
		Level:   LogLevel(level),
		Time:    time.Now(),
		Message: sb.String(),
	}

	select {
	case defaultLogger.logChan <- entry:
	default:
		fmt.Println("Log buffer full. Dropping log:", entry)
	}
}

// Updated public API
func Info(msg string, args ...interface{})  { log(INFO, msg, args...) }
func Warn(msg string, args ...interface{})  { log(WARN, msg, args...) }
func Error(msg string, args ...interface{}) { log(ERROR, msg, args...) }

func listenForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		Shutdown()
		os.Exit(0)
	}()
}
