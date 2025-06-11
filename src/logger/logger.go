package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	DEBUG   = 0b00001
	INFO    = 0b00010
	WARNING = 0b00100
	ERROR   = 0b01000
	FATAL   = 0b10000

	ALL  = DEBUG | INFO | WARNING | ERROR | FATAL
	NONE = 0
)

const (
	LOG_PATH      = "../artifacts"
	LOG_FILE      = "/logs.log"
	LOG_FULL_PATH = LOG_PATH + LOG_FILE
)

type dnxLogger struct {
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger

	LogToFile    bool
	LogToConsole bool
	LogOptions   int
	logAttempts  int
}

var dnxLoggerInstance *dnxLogger

func Init() {
	dnxLoggerInstance = &dnxLogger{
		LogToFile:    true,
		LogToConsole: true,
		LogOptions:   ALL,
		logAttempts:  0,

		DebugLogger:   log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		InfoLogger:    log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		WarningLogger: log.New(os.Stdout, "[WARNING] ", log.LstdFlags),
		ErrorLogger:   log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		FatalLogger:   log.New(os.Stderr, "[FATAL] ", log.LstdFlags),
	}

	// Create the path if it doesn't exist
	tryCreateLogFile()

	Info("Logger initialized")
}

func EnvInit() {
	Init()
	Info("Loading environment variables for logger")
	minLogLevel, existMinLevel := os.LookupEnv("DNX_LOG_MIN_LEVEL")
	disableLevels, existDisableLevels := os.LookupEnv("DNX_LOG_DISABLE_LEVELS")
	logConsole, existLogConsole := os.LookupEnv("DNX_LOG_CONSOLE")
	logFile, existLogFile := os.LookupEnv("DNX_LOG_FILE")

	if existMinLevel {
		Info("Setting minimum log level to", minLogLevel)
		SetMinLogLevel(LogLevelValue(minLogLevel))
	}
	if existDisableLevels {
		levels := strings.Split(disableLevels, "|")
		options := 0
		Info("Disabling log levels:")

		for _, level := range levels {
			level = strings.TrimSpace(level)
			options |= LogLevelValue(level)
			Info(" - ", level)
		}

		DisableLogOptions(options)
	}
	if existLogConsole {
		b, err := strconv.ParseBool(logConsole)
		if err != nil {
			Warning("Failed to parse DNX_LOG_CONSOLE value")
			SetLogToConsole(true)
			Warning("Defaulting to console logging enabled")
		} else {
			SetLogToConsole(b)
		}
	}
	if existLogFile {
		b, err := strconv.ParseBool(logFile)
		if err != nil {
			Warning("Failed to parse DNX_LOG_FILE value")
			SetLogToFile(true)
			Warning("Defaulting to file logging enabled")
		} else {
			SetLogToFile(b)
		}
	}

	Info("Logger environment variables loaded")
}

func registerLogAttempt() bool {
	dnxLoggerInstance.logAttempts++
	if dnxLoggerInstance.logAttempts > 10 && dnxLoggerInstance.logAttempts < 15 {
		Warning("Too many log attempts, will panic if this continues")
		return false
	} else if dnxLoggerInstance.logAttempts >= 15 {
		Error("Too many log attempts, will panic now")
		Fatal("Too many log attempts, this is likely a bug in the logger, please report it")
		return false
	}
	return true
}
func resetLogAttempts() {
	if dnxLoggerInstance.logAttempts > 0 {
		dnxLoggerInstance.logAttempts = 0
		Info("Log attempts counter reset")
	}
}

func tryCreateLogFile() {
	if _, err := os.Stat(LOG_PATH); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(LOG_PATH), 0755)

		if err != nil {
			SetLogToFile(false)
			Error("Failed to create log directory:", err)
			Warning("Defaulting to console logging only")
			return

		} else {
			Info("Log directory created at", LOG_PATH)
		}
	}

	if _, err := os.Stat(LOG_FULL_PATH); os.IsNotExist(err) {
		file, err := os.Create(LOG_FULL_PATH)

		if err != nil {
			SetLogToFile(false)
			logError(false, "Failed to create log file:", err)
			Warning("Defaulting to console logging only")
			return
		}
		defer file.Close()
	}

	Info("Log file created at", LOG_FULL_PATH)
}
func LogsToFile() bool {
	return dnxLoggerInstance.LogToFile
}
func SetLogToFile(value bool) {
	dnxLoggerInstance.LogToFile = value
	Info("File logging set to", value)
}

func LogsToConsole() bool {
	return dnxLoggerInstance.LogToConsole
}
func SetLogToConsole(value bool) {
	Info("Console logging set to", value)
	dnxLoggerInstance.LogToConsole = value
}

func LogOptions() int {
	return dnxLoggerInstance.LogOptions
}
func LogOptionsHas(option int) bool {
	return LogOptions()&option == option
}
func SetLogOptions(options int) {
	if options < NONE || options > ALL {
		Warning("Invalid logging options")
		return
	} else if options == ALL {
		Info("Logging options set to ALL")
		return
	} else if options == NONE {
		Info("Logging options set to NONE")
		return
	}

	msg := "Logging options set to: "

	if LogOptionsHas(DEBUG) {
		msg += "| DEBUG |"
	}
	if LogOptionsHas(INFO) {
		msg += "| INFO |"
	}
	if LogOptionsHas(WARNING) {
		msg += "| WARNING |"
	}
	if LogOptionsHas(ERROR) {
		msg += "| ERROR |"
	}
	if LogOptionsHas(FATAL) {
		msg += "| FATAL |"
	}

	Info(msg)

	dnxLoggerInstance.LogOptions = options
}
func EnableLogOptions(options int) {
	if options < NONE || options > ALL {
		Warning("Invalid logging option")
		return
	}

	var msg string

	if options&DEBUG == DEBUG {
		msg += "| DEBUG |"
	}
	if options&INFO == INFO {
		msg += "| INFO |"
	}
	if options&WARNING == WARNING {
		msg += "| WARNING |"
	}
	if options&ERROR == ERROR {
		msg += "| ERROR |"
	}
	if options&FATAL == FATAL {
		msg += "| FATAL |"
	}

	Info("Enabled logging options: ", msg)
	dnxLoggerInstance.LogOptions |= options
}
func DisableLogOptions(options int) {
	if options < NONE || options > ALL {
		Warning("Invalid logging option")
		return
	}

	var msg string

	if options&DEBUG == DEBUG {
		msg += "| DEBUG |"
	}
	if options&INFO == INFO {
		msg += "| INFO |"
	}
	if options&WARNING == WARNING {
		msg += "| WARNING |"
	}
	if options&ERROR == ERROR {
		msg += "| ERROR |"
	}
	if options&FATAL == FATAL {
		msg += "| FATAL |"
	}

	Info("Disabled logging options: ", msg)
	dnxLoggerInstance.LogOptions &= ^options
}
func SetMinLogLevel(level int) {
	if level < NONE || level > ALL {
		Warning("Invalid logging level")
		return
	}

	var msg string

	switch level {
	case DEBUG:
		msg += "| DEBUG |"
		fallthrough
	case INFO:
		msg += "| INFO |"
		fallthrough
	case WARNING:
		msg += "| WARNING |"
		fallthrough
	case ERROR:
		msg += "| ERROR |"
		fallthrough
	case FATAL:
		msg += "| FATAL |"
		fallthrough
	case ALL:
		msg = "| ALL |"
	case NONE:
		msg = "| NONE |"
	}

	Info("Minimum logging level set to: ", msg)
	SetLogOptions(level)
}
func LogLevelValue(level string) int {
	value := LogOptions()

	switch level {
	case "DEBUG":
		value = DEBUG
	case "INFO":
		value = INFO
	case "WARNING":
		value = WARNING
	case "ERROR":
		value = ERROR
	case "FATAL":
		value = FATAL
	case "ALL":
		value = ALL
	case "NONE":
		value = NONE
	default:
		Warning("Invalid logging level")
	}

	return value
}

func canLogWith(logger *log.Logger) bool {
	if LogOptionsHas(ALL) {
		return true
	} else if LogOptionsHas(NONE) {
		return false
	}

	if logger == dnxLoggerInstance.DebugLogger {
		return LogOptionsHas(DEBUG)
	} else if logger == dnxLoggerInstance.InfoLogger {
		return LogOptionsHas(INFO)
	} else if logger == dnxLoggerInstance.WarningLogger {
		return LogOptionsHas(WARNING)
	} else if logger == dnxLoggerInstance.ErrorLogger {
		return LogOptionsHas(ERROR)
	} else if logger == dnxLoggerInstance.FatalLogger {
		return LogOptionsHas(FATAL)
	}

	return false
}

func writeToFile(prefix string, v ...any) {
	file, err := os.OpenFile(LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		SetLogToFile(false)
		logError(false, "Failed to open log file")
		return
	}
	defer file.Close()

	logger := log.New(file, prefix, log.LstdFlags|log.Lshortfile)
	logger.Println(v...)
}

func logWith(logger *log.Logger, ForceWriteFile bool, v ...any) {
	if !canLogWith(logger) {
		return
	}

	registerLogAttempt()

	orignalPrefix := logger.Prefix()
	extraPrefix := ""
	_, file, line, ok := runtime.Caller(3)
	if ok {
		splitPath := strings.Split(file, "/")
		file = splitPath[len(splitPath)-1] // Get the last part of the path
		extraPrefix = fmt.Sprintf("%s:%d: ", file, line)
	} else {
		extraPrefix = "UnknownFile:0: "
	}

	trimmedArgs := strings.Trim(fmt.Sprint(v...), "[]")

	if dnxLoggerInstance.LogToConsole {
		logger.Println(extraPrefix, trimmedArgs)
	}
	if ForceWriteFile || dnxLoggerInstance.LogToFile {
		writeToFile(extraPrefix, trimmedArgs)
	}

	logger.SetPrefix(orignalPrefix) // Reset the prefix to the original one
	resetLogAttempts()
}

func logDebug(writeFile bool, v ...any) {
	logWith(dnxLoggerInstance.DebugLogger, writeFile, v...)
}
func logInfo(writeFile bool, v ...any) {
	logWith(dnxLoggerInstance.InfoLogger, writeFile, v...)
}
func logWarning(writeFile bool, v ...any) {
	logWith(dnxLoggerInstance.WarningLogger, writeFile, v...)
}
func logError(writeFile bool, v ...any) {
	logWith(dnxLoggerInstance.ErrorLogger, writeFile, v...)
}
func logFatal(writeFile bool, v ...any) {
	logWith(dnxLoggerInstance.FatalLogger, writeFile, v...)
	os.Exit(1)
}

func Debug(v ...any) {
	logDebug(false, v...)
}
func Info(v ...any) {
	logInfo(false, v...)
}
func Warning(v ...any) {
	logWarning(false, v...)
}
func Error(v ...any) {
	logError(false, v...)
}
func Fatal(v ...any) {
	logFatal(false, v...)
	os.Exit(1)
}
