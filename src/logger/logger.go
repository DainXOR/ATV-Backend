package logger

import (
	"dainxor/atv/utils"
	"regexp"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type logLevel = int

const (
	LEVEL_DEBUG   logLevel = 1 << iota // 0b00001
	LEVEL_INFO                         // 0b00010
	LEVEL_WARNING                      // 0b00100
	LEVEL_ERROR                        // 0b01000
	LEVEL_FATAL                        // 0b10000

	LEVEL_ALL           = LEVEL_DEBUG | LEVEL_INFO | LEVEL_WARNING | LEVEL_ERROR | LEVEL_FATAL
	LEVEL_NONE logLevel = 0 // 0b00000
)

func hasLogLevel(options logLevel, level logLevel) bool {
	return options&level == level
}

const (
	LOG_PATH      = "../artifacts/"
	LOG_FILE      = "logs.log"
	LOG_FULL_PATH = LOG_PATH + LOG_FILE

	DEFAULT_LOGS_TO_FILE         = false     // Default to not logging to file
	DEFAULT_LOGS_TO_CONSOLE      = true      // Default to logging to console
	DEFAULT_LOG_LEVEL            = LEVEL_ALL // Default logging level
	ENABLE_LOG_ATTEMPTS_MESSAGES = true      // Enable warning log attempts messages
	DEFAULT_MAX_LOG_ATTEMPTS     = 15        // Default maximum log attempts before panic
	DEFAULT_WARNING_LOG_ATTEMPTS = 10        // Default maximum log attempts before warning

	DEFAULT_APP_VERSION = "0.1.0" // Default application version

	TXT_BLACK   = "30m"
	TXT_RED     = "31m"
	TXT_GREEN   = "32m"
	TXT_YELLOW  = "33m"
	TXT_BLUE    = "34m"
	TXT_MAGENTA = "35m"
	TXT_CYAN    = "36m"
	TXT_WHITE   = "37m"

	BG_BLACK   = "40"
	BG_RED     = "41"
	BG_GREEN   = "42"
	BG_YELLOW  = "43"
	BG_BLUE    = "44"
	BG_MAGENTA = "45"
	BG_CYAN    = "46"
	BG_WHITE   = "47"

	CLR_START = "\033["
	CLR_RESET = "\033[0m"

	CLR_DEBUG = CLR_START + BG_GREEN + ";" + TXT_BLACK  // Green background with black text for debug messages
	CLR_INFO  = CLR_START + BG_CYAN + ";" + TXT_BLACK   // Cyan background with black text for info messages
	CLR_WARN  = CLR_START + BG_YELLOW + ";" + TXT_BLACK // Yellow background with black text for warning messages
	CLR_ERROR = CLR_START + BG_RED + ";" + TXT_BLACK    // Red background with black text for error messages
	CLR_FATAL = CLR_START + BG_RED + ";" + TXT_WHITE    // Red background with white text for fatal messages

	CLR_DEPRECATE   = CLR_START + BG_MAGENTA + ";" + TXT_WHITE // Magenta background with white text for deprecation messages
	CLR_DEPR_REASON = CLR_START + BG_YELLOW + ";" + TXT_WHITE  // Yellow background with white text for deprecation reason messages

	CLR_LAVA       = CLR_START + BG_WHITE + ";" + TXT_BLACK  // White background with black text for lava messages
	CLR_COLD_LAVA  = CLR_START + BG_YELLOW + ";" + TXT_BLACK // Yellow background with black text for cold lava messages
	CLR_DRIED_LAVA = CLR_START + BG_RED + ";" + TXT_BLACK    // Red background with black text for dried lava messages

	CLR_FILE = CLR_START + BG_BLUE + ";" + TXT_WHITE // Blue background with white text for file paths
)

type dnxLogger struct {
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger

	LogToFile    bool
	LogToConsole bool
	LogLevels    logLevel
	logAttempts  int

	appVersion      string
	appVersionMajor uint64
	appVersionMinor uint64
	appVersionPatch uint64
}

var dnxLoggerInstance *dnxLogger

// Returns the singleton instance of dnxLogger, initializing it if necessary.
// This function should be used to access the logger throughout the internal package.
// It abstracts the initialization logic and provides a single point of access to the logger instance.
// It ensures that the logger is initialized only once, and provides a consistent interface for logging.
func get() *dnxLogger {
	if dnxLoggerInstance == nil {
		Init()
	}
	return dnxLoggerInstance
}
func colorTxt(txt string, textColor string, bgColor string) string {
	return CLR_START + bgColor + ";" + textColor + txt + CLR_RESET
}
func colorWith(txt string, colorString string) string {
	return colorString + txt + CLR_RESET
}

func calcMajorVersion(version string) uint64 {
	extractedStr := utils.Extract("", version, ".")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func calcMinorVersion(version string) uint64 {
	extractedStr := utils.Extract(".", version, ".")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func calcPatchVersion(version string) uint64 {
	extractedStr := utils.Extract(".", version, ".")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func compareVersions(ver1, ver2 string) int8 {
	major1 := calcMajorVersion(ver1)
	minor1 := calcMinorVersion(ver1)
	patch1 := calcPatchVersion(ver1)

	major2 := calcMajorVersion(ver2)
	minor2 := calcMinorVersion(ver2)
	patch2 := calcPatchVersion(ver2)

	return compareVersionsNum(major1, minor1, patch1, major2, minor2, patch2)
}
func compareVersionsNum(major1, minor1, patch1, major2, minor2, patch2 uint64) int8 {
	if major1 > major2 {
		return 1
	} else if major1 < major2 {
		return -1
	}

	if minor1 > minor2 {
		return 1
	} else if minor1 < minor2 {
		return -1
	}

	if patch1 > patch2 {
		return 1
	} else if patch1 < patch2 {
		return -1
	}

	return 0
}
func addToVersion(version string, major, minor, patch uint64) string {
	majorVersion := calcMajorVersion(version) + major
	minorVersion := calcMinorVersion(version) + minor
	patchVersion := calcPatchVersion(version) + patch

	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion)
}

func Init() {
	dnxLoggerInstance = &dnxLogger{
		LogToFile:       DEFAULT_LOGS_TO_FILE,
		LogToConsole:    DEFAULT_LOGS_TO_CONSOLE,
		LogLevels:       DEFAULT_LOG_LEVEL,
		logAttempts:     0,
		appVersion:      DEFAULT_APP_VERSION,
		appVersionMajor: calcMajorVersion(DEFAULT_APP_VERSION),
		appVersionMinor: calcMinorVersion(DEFAULT_APP_VERSION),
		appVersionPatch: calcPatchVersion(DEFAULT_APP_VERSION),

		DebugLogger:   log.New(os.Stdout, "|"+colorWith(" DEBUG ", CLR_DEBUG)+"| ", log.LstdFlags),
		InfoLogger:    log.New(os.Stdout, "|"+colorWith(" INFO ", CLR_INFO)+"| ", log.LstdFlags),
		WarningLogger: log.New(os.Stdout, "|"+colorWith(" WARNING ", CLR_WARN)+"| ", log.LstdFlags),
		ErrorLogger:   log.New(os.Stderr, "|"+colorWith(" ERROR ", CLR_ERROR)+"| ", log.LstdFlags),
		FatalLogger:   log.New(os.Stderr, "|"+colorWith(" FATAL ", CLR_FATAL)+"| ", log.LstdFlags),
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
		Info("Setting minimum log level to ", minLogLevel)
		SetMinLogLevel(LogLevelValue(minLogLevel))
	}
	if existDisableLevels {
		levels := strings.Split(disableLevels, "|")
		options := LEVEL_NONE
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
			Warning("Defaulting to console logging: ", DEFAULT_LOGS_TO_CONSOLE)
			SetLogToConsole(DEFAULT_LOGS_TO_CONSOLE)
		} else {
			SetLogToConsole(b)
		}
	}
	if existLogFile {
		b, err := strconv.ParseBool(logFile)
		if err != nil {
			Warning("Failed to parse DNX_LOG_FILE value")
			Warning("Defaulting to file logging: ", DEFAULT_LOGS_TO_FILE)
			SetLogToFile(DEFAULT_LOGS_TO_FILE)
		} else {
			SetLogToFile(b)
		}
	}

	Info("Logger environment variables loaded")
}

func tryCreateLogFile() bool {
	if _, err := os.Stat(LOG_PATH); os.IsNotExist(err) {
		logWarning(true, "Log directory does not exist")
		logInfo(true, "Attempting to create log directory at ", LOG_PATH)
		err := os.MkdirAll(filepath.Dir(LOG_PATH), 0755)

		if err != nil {
			logError(true, "Failed to create log directory: ", err)
			logWarning(true, "Defaulting to no file logging")
			SetLogToFile(false)
			return false

		} else {
			logInfo(true, "Log directory created at ", LOG_PATH)
		}
	}

	if _, err := os.Stat(LOG_FULL_PATH); os.IsNotExist(err) {
		file, err := os.Create(LOG_FULL_PATH)

		if err != nil {
			logError(true, "Failed to create log file: ", err)
			logWarning(true, "Defaulting to no file logging")
			SetLogToFile(false)
			return false
		} else {
			file.Close()
			logInfo(true, "Log file created at ", LOG_FULL_PATH)
			return true
		}
	}

	Info("Log file already exists at ", LOG_FULL_PATH)
	return true
}

func registerLogAttempt(forceNoFileWrite bool) bool {
	get().logAttempts++

	if get().logAttempts > DEFAULT_WARNING_LOG_ATTEMPTS && get().logAttempts < DEFAULT_MAX_LOG_ATTEMPTS {
		if ENABLE_LOG_ATTEMPTS_MESSAGES {
			logWarning(forceNoFileWrite, "Too many log attempts, will panic if this continues")
		}
		return false

	} else if get().logAttempts >= DEFAULT_MAX_LOG_ATTEMPTS {
		if ENABLE_LOG_ATTEMPTS_MESSAGES {
			logError(forceNoFileWrite, "Too many log attempts, will panic now")
			logFatal(forceNoFileWrite, "Too many log attempts, this is likely a bug in the logger, please report it")
		}

		return false
	}
	return true
}
func resetLogAttempts(forceNoFileWrite bool) bool {
	if get().logAttempts == 1 {
		get().logAttempts = 0
		return true

	} else if get().logAttempts > 1 {
		get().logAttempts = 0
		logInfo(forceNoFileWrite, "Log attempts reset")
		return true
	}

	return false
}

func LogsToFile() bool {
	return get().LogToFile
}
func SetLogToFile(value bool) {
	logInfo(!value, "File logging set to ", value)
	get().LogToFile = value
}

func LogsToConsole() bool {
	return get().LogToConsole
}
func SetLogToConsole(value bool) {
	Info("Console logging set to ", value)
	get().LogToConsole = value
}

func SetAppVersion(version string) {
	regex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	if !regex.MatchString(version) {
		Warning("Invalid app version format")
		return
	}

	get().appVersion = version
	get().appVersionMajor = calcMajorVersion(version)
	get().appVersionMinor = calcMinorVersion(version)
	get().appVersionPatch = calcPatchVersion(version)
	Info("App version set to: ", version)
}
func AppVersion() string {
	return get().appVersion
}

func LogLevels() logLevel {
	return get().LogLevels
}
func LogLevelsHas(option logLevel) bool {
	return LogLevels()&option == option
}
func SetLogLevels(options logLevel) bool {
	if options < LEVEL_NONE || options > LEVEL_ALL {
		Warning("Invalid logging options")
		return false
	} else if options == LEVEL_ALL {
		Info("Logging options set to ALL")
		get().LogLevels = LEVEL_ALL
		return true
	} else if options == LEVEL_NONE {
		Info("Logging options set to NONE")
		get().LogLevels = LEVEL_NONE
		return true
	}

	var msg string

	if LogLevelsHas(LEVEL_DEBUG) {
		msg += "| DEBUG |"
	}
	if LogLevelsHas(LEVEL_INFO) {
		msg += "| INFO |"
	}
	if LogLevelsHas(LEVEL_WARNING) {
		msg += "| WARNING |"
	}
	if LogLevelsHas(LEVEL_ERROR) {
		msg += "| ERROR |"
	}
	if LogLevelsHas(LEVEL_FATAL) {
		msg += "| FATAL |"
	}

	Info("Logging options set to: ", msg)
	get().LogLevels = options
	return true
}
func EnableLogOptions(options logLevel) bool {
	if options < LEVEL_NONE || options > LEVEL_ALL {
		Warning("Invalid logging option")
		return false
	} else if options == LEVEL_ALL {
		Info("Enabled all logging options")
		SetLogLevels(LEVEL_ALL)
		return true
	} else if options == LEVEL_NONE {
		Info("Disabled all logging options")
		SetLogLevels(LEVEL_NONE)
		return true
	}

	var msg string

	if hasLogLevel(options, LEVEL_DEBUG) {
		msg += "| DEBUG |"
	}
	if hasLogLevel(options, LEVEL_INFO) {
		msg += "| INFO |"
	}
	if hasLogLevel(options, LEVEL_WARNING) {
		msg += "| WARNING |"
	}
	if hasLogLevel(options, LEVEL_ERROR) {
		msg += "| ERROR |"
	}
	if hasLogLevel(options, LEVEL_FATAL) {
		msg += "| FATAL |"
	}

	Info("Enabled logging options: ", msg)
	get().LogLevels |= options
	return true
}
func DisableLogOptions(options logLevel) bool {
	if options < LEVEL_NONE || options > LEVEL_ALL {
		Warning("Invalid logging option")
		return false
	}

	var msg string

	if hasLogLevel(options, LEVEL_DEBUG) {
		msg += "| DEBUG |"
	}
	if hasLogLevel(options, LEVEL_INFO) {
		msg += "| INFO |"
	}
	if hasLogLevel(options, LEVEL_WARNING) {
		msg += "| WARNING |"
	}
	if hasLogLevel(options, LEVEL_ERROR) {
		msg += "| ERROR |"
	}
	if hasLogLevel(options, LEVEL_FATAL) {
		msg += "| FATAL |"
	}

	Info("Disabled logging options: ", msg)
	get().LogLevels &= ^options
	return true
}
func SetMinLogLevel(level logLevel) bool {
	if level < LEVEL_NONE || level > LEVEL_ALL {
		Warning("Invalid logging level")
		return false
	}

	var msg string

	switch level {
	case LEVEL_ALL:
		fallthrough
	case LEVEL_DEBUG:
		msg = "DEBUG"

	case LEVEL_INFO:
		msg = "INFO"
	case LEVEL_WARNING:
		msg = "WARNING"
	case LEVEL_ERROR:
		msg = "ERROR"
	case LEVEL_FATAL:
		msg = "FATAL"
	case LEVEL_NONE:
		msg = "NONE"
	}

	Info("Minimum logging level set to: ", msg)
	SetLogLevels(level)
	return true
}
func LogLevelValue(levelName string) logLevel {
	value := LogLevels()

	switch levelName {
	case "DEBUG":
		value = LEVEL_DEBUG
	case "INFO":
		value = LEVEL_INFO
	case "WARNING":
		value = LEVEL_WARNING
	case "ERROR":
		value = LEVEL_ERROR
	case "FATAL":
		value = LEVEL_FATAL
	case "ALL":
		value = LEVEL_ALL
	case "NONE":
		value = LEVEL_NONE
	default:
		Warning("Invalid logging level")
	}

	return value
}
func canLogWith(logger *log.Logger) bool {
	if LogLevelsHas(LEVEL_ALL) {
		return true
	} else if LogLevelsHas(LEVEL_NONE) {
		return false
	}

	if logger == get().DebugLogger {
		return LogLevelsHas(LEVEL_DEBUG)
	} else if logger == get().InfoLogger {
		return LogLevelsHas(LEVEL_INFO)
	} else if logger == get().WarningLogger {
		return LogLevelsHas(LEVEL_WARNING)
	} else if logger == get().ErrorLogger {
		return LogLevelsHas(LEVEL_ERROR)
	} else if logger == get().FatalLogger {
		return LogLevelsHas(LEVEL_FATAL)
	}

	return false
}

func writeToFile(logger *log.Logger, prefix string, v ...any) bool {
	file, err := os.OpenFile(LOG_FULL_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		logError(true, "Failed to open log file")
		return false
	}
	defer file.Close()

	stringValues := utils.AsStrings(v)
	stringValues = utils.Podate(stringValues, "[ ]")
	trimmedArgs := strings.Join(stringValues, " ")
	trimmedArgs = strings.Trim(trimmedArgs, "[]")

	originalOutput := logger.Writer()
	logger.SetOutput(file)
	logger.Println(prefix, trimmedArgs)
	logger.SetOutput(originalOutput)
	return true
}

// Private function that handles the actual logging.
// TO DO: Add option to modify depth of the call stack to log origin
func logWith(logger *log.Logger, forceNoFileWrite bool, v ...any) {
	if !canLogWith(logger) {
		return
	}

	registerLogAttempt(forceNoFileWrite)

	orignalPrefix := logger.Prefix()
	extraPrefix := utils.CallOrigin(4)
	extraPrefix = colorWith(extraPrefix, CLR_FILE)
	extraPrefix += ":"

	stringValues := utils.AsStrings(v)
	stringValues = utils.Podate(stringValues, "[ ]")
	trimmedArgs := strings.Join(stringValues, " ")
	trimmedArgs = strings.Trim(trimmedArgs, "[]")

	if LogsToConsole() {
		logger.Println(extraPrefix, trimmedArgs)
	}

	// This might be short-circuited, but idk if it is, so I will leave it like this
	if !forceNoFileWrite && LogsToFile() {
		if !writeToFile(logger, extraPrefix, trimmedArgs) {
			SetLogToFile(false)
		}
	}

	logger.SetPrefix(orignalPrefix) // Reset the prefix to the original one
	resetLogAttempts(forceNoFileWrite)
}

// Private functions for logging at different levels
// These functions are used internally and should not be called directly by users.
// Use forceNoFileWrite to prevent writing to the log file, useful for testing or when file logging is disabled or not possible.
func logDebug(forceNoFileWrite bool, v ...any) {
	logWith(get().DebugLogger, forceNoFileWrite, v...)
}
func logInfo(forceNoFileWrite bool, v ...any) {
	logWith(get().InfoLogger, forceNoFileWrite, v...)
}
func logWarning(forceNoFileWrite bool, v ...any) {
	logWith(get().WarningLogger, forceNoFileWrite, v...)
}
func logError(forceNoFileWrite bool, v ...any) {
	logWith(get().ErrorLogger, forceNoFileWrite, v...)
}
func logFatal(forceNoFileWrite bool, v ...any) {
	logWith(get().FatalLogger, forceNoFileWrite, v...)
	panic(fmt.Sprint(v...))
}

// Public functions for logging at different levels
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
}

func Deprecate(deprecatedVersion string, removalVersion string, v ...any) (bool, error) {
	args := utils.Join(v, " ")
	deprecateTxt := colorWith(" DEPRECATED ", CLR_DEPRECATE)
	reasonTxt := colorWith(" REASON : ", CLR_DEPR_REASON)

	if compareVersions(AppVersion(), addToVersion(removalVersion, 0, 1, 0)) > 0 {
		Fatal(deprecateTxt+": This feature has been removed in version", removalVersion)
		Fatal(reasonTxt+":", args)

	} else if compareVersions(AppVersion(), removalVersion) >= 0 {
		Error(deprecateTxt+": This feature has been removed in version ", removalVersion)
		Error(reasonTxt+":", args)
		return false, fmt.Errorf("feature removed in version %s. %s", removalVersion, args)

	} else if compareVersions(AppVersion(), deprecatedVersion) >= 0 && compareVersions(AppVersion(), removalVersion) < 0 {
		Warning(deprecateTxt+": This feature is marked for removal in version ", removalVersion)
		Warning(reasonTxt+":", args)
		return false, fmt.Errorf("feature deprecated in version %s, will be removed in version %s. %s", deprecatedVersion, removalVersion, args)

	} else if compareVersions(AppVersion(), deprecatedVersion) == 0 {
		Warning(deprecateTxt + ": This feature will be removed in future versions")
		Warning(reasonTxt+":", args)
		return true, fmt.Errorf("feature deprecated in version %s. %s", deprecatedVersion, args)
	}
	return true, nil
}
func DeprecateMsg(deprecatedVersion string, removalVersion string, v ...any) string {
	_, err := Deprecate(deprecatedVersion, removalVersion, v...)
	return err.Error()
}

type volcano struct {
	version string
}

func Lava(version string, v ...any) volcano {
	args := utils.Join(v, " ")
	lavaTxt := colorWith(" LAVA ", CLR_LAVA)
	coldLavaTxt := colorWith(" COLD LAVA ", CLR_COLD_LAVA)
	driedLavaTxt := colorWith(" DRIED LAVA ", CLR_DRIED_LAVA)

	if compareVersions(AppVersion(), version) == 0 {
		Warning(lavaTxt+": Running code that should be removed, cleaned up or refactored. ", args)
	} else if compareVersions(AppVersion(), addToVersion(version, 0, 2, 0)) >= 0 {
		Warning(coldLavaTxt+": This code must be refactored. ", args)
	} else if compareVersions(AppVersion(), addToVersion(version, 0, 4, 0)) > 0 {
		Error(driedLavaTxt+": This code should not be running as is, it is likely a bug. ", args)
	}

	return volcano{
		version: version,
	}
}
func (v *volcano) LavaStart() {
	lavaTxt := colorWith(" LAVA ", CLR_LAVA)
	coldLavaTxt := colorWith(" COLD LAVA ", CLR_COLD_LAVA)
	driedLavaTxt := colorWith(" DRIED LAVA ", CLR_DRIED_LAVA)

	if AppVersion() == v.version {
		Warning(lavaTxt + ": Start of flow")
	} else if compareVersions(AppVersion(), addToVersion(v.version, 0, 2, 0)) >= 0 {
		Warning(coldLavaTxt + ": Start of flow")
	} else if compareVersions(AppVersion(), addToVersion(v.version, 0, 4, 0)) > 0 {
		Error(driedLavaTxt + ": Start of lava cast")
	}
}
func (v *volcano) LavaEnd() {
	lavaTxt := colorWith(" LAVA ", CLR_LAVA)
	coldLavaTxt := colorWith(" COLD LAVA ", CLR_COLD_LAVA)
	driedLavaTxt := colorWith(" DRIED LAVA ", CLR_DRIED_LAVA)

	if AppVersion() == v.version {
		Warning(lavaTxt + ": End of flow")
	} else if compareVersions(AppVersion(), addToVersion(v.version, 0, 2, 0)) >= 0 {
		Warning(coldLavaTxt + ": End of flow")
	} else if compareVersions(AppVersion(), addToVersion(v.version, 0, 4, 0)) > 0 {
		Error(driedLavaTxt + ": End of lava cast")
	}
}
