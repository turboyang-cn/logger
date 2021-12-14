package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Output int8

type LevelConfig struct {
	Output Output // 输出
	Prefix string // 前缀
	Flag   int    // 标志
}

type Config struct {
	LogPath string      // 日志文件路径
	Trace   LevelConfig // Trace 级别配置
	Debug   LevelConfig // Debug 级别配置
	Info    LevelConfig // Info 级别配置
	Warning LevelConfig // Warning 级别配置
	Error   LevelConfig // Error 级别配置
}

type Logger struct {
	config      Config
	currentTime time.Time
	logger      *log.Logger
	file        *os.File
	mutex       *sync.Mutex
}

const (
	Stdout Output = 1 << iota
	StdErr
	File
)

var DefaultConfig = Config{
	Trace: LevelConfig{
		Prefix: fmt.Sprintf("%-10s", "Trace"),
		Output: Stdout | File,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	},
	Debug: LevelConfig{
		Prefix: fmt.Sprintf("%-10s", "Debug"),
		Output: Stdout | File,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	},
	Info: LevelConfig{
		Prefix: fmt.Sprintf("%-10s", "Info"),
		Output: Stdout | File,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	},
	Warning: LevelConfig{
		Prefix: fmt.Sprintf("%-10s", "Warning"),
		Output: Stdout | File,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	},
	Error: LevelConfig{
		Prefix: fmt.Sprintf("%-10s", "Error"),
		Output: Stdout | File,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	},
}

var DefaultLogger = NewLogger(DefaultConfig)

func NewLogger(config Config) Logger {
	logger := Logger{
		config: config,
		mutex:  new(sync.Mutex),
	}

	return logger
}

func (logger *Logger) resetLogger(config LevelConfig) {
	writers := make([]io.Writer, 0)

	if config.Output&Stdout == Stdout {
		writers = append(writers, os.Stdout)
	}
	if config.Output&StdErr == StdErr {
		writers = append(writers, os.Stderr)
	}
	if config.Output&File == File {
		if currentTime := time.Now(); currentTime.Format("2006-01-02") != logger.currentTime.Format("2006-01-02") {
			logger.currentTime = currentTime

			directory := fmt.Sprintf("logs/%s", currentTime.Format("2006-01-02"))
			os.MkdirAll(directory, 0666)
			file, _ := os.OpenFile(fmt.Sprintf("%s/log.log", directory), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

			logger.file = file
		}

		writers = append(writers, logger.file)
	}

	writer := io.MultiWriter(writers...)

	logger.logger = log.New(writer, config.Prefix, config.Flag)
}

func (logger *Logger) Trace(message interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.resetLogger(logger.config.Trace)

	logger.logger.Println(message)
}

func (logger *Logger) Debug(message interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.resetLogger(logger.config.Debug)

	logger.logger.Println(message)
}

func (logger *Logger) Info(message interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.resetLogger(logger.config.Info)

	logger.logger.Println(message)
}

func (logger *Logger) Warning(message interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.resetLogger(logger.config.Warning)

	logger.logger.Println(message)
}

func (logger *Logger) Error(message interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.resetLogger(logger.config.Error)

	logger.logger.Println(message)
}
