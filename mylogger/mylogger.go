package mylogger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kohinigeee/mylog/clog"
)

// -----------------初期化関数-----------------

var (
	iniialAllowMode allowModeT = modeAllAllow
	// iniialAllowMode allowModeT = modeManualAllow
	loggerManager *LoggerManager
	mainLoggerID  loggerIDType = "MainLogger"
)

// -----------------ログID-----------------
const (
	LogIdConstructGreedy loggerIDType = "ConstructGreedy"
)

func init() {
	const logFileName = ""
	loggerManager = newLoggerManager(iniialAllowMode, logFileName)

	_ = loggerManager.GetLogger(mainLoggerID)

	if loggerManager.allowMode == modeManualAllow {
		initialAllowdLoggerIDs := []loggerIDType{
			// mainLoggerID,
			LogIdConstructGreedy,
		}

		for _, id := range initialAllowdLoggerIDs {
			loggerManager.AddAllowID(id)
		}
	}
}

func GetLogger(id loggerIDType) *LoggerItem {
	return loggerManager.GetLogger(id)
}

func L() *LoggerItem {
	return GetLogger(mainLoggerID)
}

const (
	defaultLogLevel slog.Level = slog.LevelDebug
)

func makeLogger(logLevel *slog.LevelVar, w io.Writer) *slog.Logger {
	handler, err := clog.NewCustomTextHandler(w,
		clog.WithHandlerOption(&slog.HandlerOptions{
			Level: logLevel,
		}))

	if err != nil {
		panic(err)
	}

	logger := slog.New(handler)
	return logger
}

type loggerIDType string

// -----------------LoggerItem-----------------
type LoggerItem struct {
	id       loggerIDType
	logLevel *slog.LevelVar
	logger   *slog.Logger
	isShown  bool
}

type logType int

const (
	logTypeDebug logType = iota
	logTypeInfo
	logTypeWarn
	logTypeError
)

func newLoggerItem(id loggerIDType, writer io.Writer) *LoggerItem {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	return &LoggerItem{
		id:       id,
		logLevel: logLevel,
		logger:   makeLogger(logLevel, writer),
		isShown:  true,
	}
}

func (li *LoggerItem) log(ltype logType, msg string, args ...any) {
	if !li.isShown {
		return
	}

	msg = fmt.Sprintf("[ID:%s] %s", li.id, msg)

	switch ltype {
	case logTypeDebug:
		li.logger.Debug(msg, args...)
	case logTypeInfo:
		li.logger.Info(msg, args...)
	case logTypeWarn:
		li.logger.Warn(msg, args...)
	case logTypeError:
		li.logger.Error(msg, args...)
	}
}

func (li *LoggerItem) Info(msg string, args ...any) {
	li.log(logTypeInfo, msg, args...)
}

func (li *LoggerItem) Debug(msg string, args ...any) {
	li.log(logTypeDebug, msg, args...)
}

func (li *LoggerItem) Warn(msg string, args ...any) {
	li.log(logTypeWarn, msg, args...)
}

func (li *LoggerItem) Error(msg string, args ...any) {
	li.log(logTypeError, msg, args...)
}

func (li *LoggerItem) SetLevel(level slog.Level) {
	li.logLevel.Set(level)
}

// -----------------LoggerManager-----------------

type allowModeT int

const (
	modeAllAllow    allowModeT = iota
	modeManualAllow            //初期化時に許可するLoggerIDを指定する
)

type LoggerManager struct {
	loggerMap     map[loggerIDType]*LoggerItem
	allowedLogger map[loggerIDType]any
	allowMode     allowModeT
	fileName      string //ログ出力先ファイル名
	writer        io.Writer
}

func newLoggerManager(allowMode allowModeT, logFname string) *LoggerManager {
	lm := &LoggerManager{
		loggerMap:     make(map[loggerIDType]*LoggerItem),
		allowedLogger: map[loggerIDType]any{},
		allowMode:     allowMode,
		fileName:      logFname,
	}

	if logFname == "" {
		lm.writer = os.Stdout
	} else {
		fileDir := "./Log"
		fpath := filepath.Join(fileDir, logFname)
		file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		lm.writer = file
	}
	return lm
}

func (lm *LoggerManager) IsShown(id loggerIDType) bool {
	if _, ok := lm.allowedLogger[id]; ok {
		return true
	}

	return false
}

// マニュアルモードでのみ使用可
func (lm *LoggerManager) AddAllowID(id loggerIDType) {
	if lm.allowMode != modeManualAllow {
		return
	}

	lm.allowedLogger[id] = nil
}

func (lm *LoggerManager) GetLogger(id loggerIDType) *LoggerItem {
	if LoggerItem, ok := lm.loggerMap[id]; ok {
		LoggerItem.isShown = lm.IsShown(id)
		return LoggerItem
	}

	LoggerItem := newLoggerItem(id, lm.writer)
	lm.loggerMap[id] = LoggerItem

	switch lm.allowMode {
	case modeAllAllow:
		lm.allowedLogger[id] = nil
	default:
	}

	LoggerItem.isShown = lm.IsShown(id)
	return LoggerItem
}

func (lm *LoggerManager) SetLevel(id loggerIDType, level slog.Level) {
	if item, ok := lm.loggerMap[id]; ok {
		item.logLevel.Set(level)
	} else {
		panic(fmt.Sprintf("Logger ID [%s] is not found : SetLevel", id))
	}
}

func Close() {
	if loggerManager.fileName != "" {
		loggerManager.writer.(io.WriteCloser).Close()
	}
}

func ReInitializeManager() {
	Close()
	// loggerManager = newLoggerManager(iniialAllowMode, global.LogFileName)
}
