package logger

import (
	"dainxor/atv/utils"
	"regexp"

	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var ( // Default color functions for logging
	colorAsDebug   = colorApplier(TXT_BLACK, BG_GREEN)
	colorAsInfo    = colorApplier(TXT_BLACK, BG_CYAN)
	colorAsWarning = colorApplier(TXT_BLACK, BG_YELLOW)
	colorAsError   = colorApplier(TXT_BLACK, BG_RED)
	colorAsFatal   = colorApplier(TXT_WHITE, BG_RED)

	colorAsDeprecate  = colorApplier(TXT_WHITE, BG_MAGENTA)
	colorAsDeprReason = colorApplier(TXT_WHITE, BG_YELLOW)

	colorAsLava      = colorApplier(TXT_BLACK, BG_WHITE)
	colorAsColdLava  = colorApplier(TXT_BLACK, BG_YELLOW)
	colorAsDriedLava = colorApplier(TXT_BLACK, BG_RED)

	colorAsFile = colorApplier(TXT_WHITE, BG_BLUE)
)

func colorApplier(textColor, bgColor AnsiCode) func(txt string) string {
	if !LogsWithColor() {
		return func(txt string) string {
			return txt // If color logging is disabled, return the text as is
		}
	}

	return func(txt string) string {
		return txt //CLR_START + bgColor + ";" + textColor + txt + CLR_RESET
	}
}

type dnxLogger struct {
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger

	LogToFile    bool
	LogToConsole bool
	ColorLogs    bool
	LogLevels    logLevel
	logAttempts  int

	appVersion      string
	appVersionMajor uint64
	appVersionMinor uint64
	appVersionPatch uint64

	usingDefaults bool // Indicates if the logger is using default settings
}

var dnxLoggerInstance *dnxLogger

func init() {
	defaultInit()
}
func defaultInit() {
	dnxLoggerInstance = &dnxLogger{
		//		LogToFile:       DEFAULT_LOGS_TO_FILE,
		//		LogToConsole:    DEFAULT_LOGS_TO_CONSOLE,
		//		ColorLogs:       DEFAULT_COLOR_LOGGING,
		//		LogLevels:       DEFAULT_LOG_LEVEL,
		//		logAttempts:     0,
		//		appVersion:      DEFAULT_APP_VERSION,
		//		appVersionMajor: majorVersionOf(DEFAULT_APP_VERSION),
		//		appVersionMinor: minorVersionOf(DEFAULT_APP_VERSION),
		//		appVersionPatch: patchVersionOf(DEFAULT_APP_VERSION),
		usingDefaults: true,

		DebugLogger:   log.New(os.Stdout, "| DEBUG | ", log.LstdFlags),
		InfoLogger:    log.New(os.Stdout, "| INFO | ", log.LstdFlags),
		WarningLogger: log.New(os.Stdout, "| WARNING | ", log.LstdFlags),
		ErrorLogger:   log.New(os.Stderr, "| ERROR | ", log.LstdFlags),
		FatalLogger:   log.New(os.Stderr, "| FATAL | ", log.LstdFlags),
	}

	envInit() // Initialize environment variables for logger

	if LogsToFile() {
		tryCreateLogFile() // Create the path if it doesn't exist, else it will set log to console only
	}

	Debug("Logger initialized")
}
func envInit() {
	Debug("Loading environment variables for logger")

	minLogLevel, existMinLevel := os.LookupEnv("DNX_LOG_MIN_LEVEL")
	disableLevels, existDisableLevels := os.LookupEnv("DNX_LOG_DISABLE_LEVELS")
	logConsole, existLogConsole := os.LookupEnv("DNX_LOG_CONSOLE")
	logFile, existLogFile := os.LookupEnv("DNX_LOG_FILE")
	logWithColor, existLogWithColor := os.LookupEnv("DNX_LOG_WITH_COLOR")

	if existMinLevel {
		level, err := Level.Get(minLogLevel)
		if err != nil {
			Warning("Failed to parse DNX_LOG_MIN_LEVEL value, using default level: ", currentLogLevels())
			level = Level.All()
		}

		Debug("Setting minimum log level to ", minLogLevel)
		get().usingDefaults = !SetMinLogLevel(level) // If any environment variable is set, we are not using defaults
	} else {
		Debug("DNX_LOG_MIN_LEVEL not set, using default level: ", currentLogLevels())
	}
	if existDisableLevels {
		levels := strings.Split(disableLevels, "|")
		//options := LEVEL_NONE
		Debug("Disabling log levels:")

		for _, level := range levels {
			level = strings.TrimSpace(level)
			//options |= LogLevelValue(level)
			Debug(" - ", level)
		}

		//get().usingDefaults = !DisableLogOptions(options) && get().usingDefaults // If any environment variable is set, we are not using defaults
	} else {
		Debug("DNX_LOG_DISABLE_LEVELS not set, keeping current log levels: ", currentLogLevels())
	}
	if existLogConsole {
		b, err := strconv.ParseBool(logConsole)
		if err != nil {
			Warning("Failed to parse DNX_LOG_CONSOLE value")
			//Warning("Defaulting to console logging: ", DEFAULT_LOGS_TO_CONSOLE)
			//SetLogToConsole(DEFAULT_LOGS_TO_CONSOLE)
		} else {
			SetLogToConsole(b)
			get().usingDefaults = false // If any environment variable is set, we are not using defaults
		}
	} else {
		//Debug("DNX_LOG_CONSOLE not set, using default value: ", DEFAULT_LOGS_TO_CONSOLE)
	}
	if existLogFile {
		b, err := strconv.ParseBool(logFile)
		if err != nil {
			Warning("Failed to parse DNX_LOG_FILE value")
			//Warning("Defaulting to file logging: ", DEFAULT_LOGS_TO_FILE)
			//SetLogToFile(DEFAULT_LOGS_TO_FILE)
		} else {
			SetLogToFile(b)
			get().usingDefaults = false // If any environment variable is set, we are not using defaults
		}
	} else {
		//Debug("DNX_LOG_FILE not set, using default value: ", DEFAULT_LOGS_TO_FILE)
	}
	if existLogWithColor {
		b, err := strconv.ParseBool(logWithColor)
		if err != nil {
			Warning("Failed to parse DNX_LOG_WITH_COLOR value")
			//Warning("Defaulting to color logging: ", DEFAULT_COLOR_LOGGING)
			//SetLogWithColor(DEFAULT_COLOR_LOGGING)
		} else {
			SetLogWithColor(b)
			get().usingDefaults = false // If any environment variable is set, we are not using defaults
		}
	} else {
		//Debug("DNX_LOG_WITH_COLOR not set, using default value: ", DEFAULT_COLOR_LOGGING)
		get().DebugLogger.SetPrefix("|" + Ansi.Debug().Apply(" DEBUG ") + "| ")
		get().InfoLogger.SetPrefix("|" + Ansi.Info().Apply(" INFO ") + "| ")
		get().WarningLogger.SetPrefix("|" + Ansi.Warn().Apply(" WARNING ") + "| ")
		get().ErrorLogger.SetPrefix("|" + Ansi.Error().Apply(" ERROR ") + "| ")
		get().FatalLogger.SetPrefix("|" + Ansi.Fatal().Apply(" FATAL ") + "| ")
	}

	Debug("Logger environment variables loaded")
}

func ReloadEnv() {
	envInit()
}

func hasLogLevel(options logLevel, level logLevel) bool {
	return options.Has(level)
}

// Returns the singleton instance of dnxLogger, initializing it if necessary.
// This function should be used to access the logger throughout the internal package.
// It abstracts the initialization logic and provides a single point of access to the logger instance.
// It ensures that the logger is initialized only once, and provides a consistent interface for logging.
func get() *dnxLogger {
	if dnxLoggerInstance == nil {
		defaultInit()
	}
	return dnxLoggerInstance
}
func UsingDefaults() bool {
	return get().usingDefaults
}
func colorTxt(txt string, textColor string, bgColor string) string {
	if !LogsWithColor() {
		return txt // If color logging is disabled, return the text as is
	}

	return txt //CLR_START + bgColor + ";" + textColor + txt + CLR_RESET
}
func colorWith(txt string, colorString string) string {
	if !LogsWithColor() {
		return txt // If color logging is disabled, return the text as is
	}

	return txt //colorString + txt + CLR_RESET
}

func majorVersionOf(version string) uint64 {
	extractedStr := strings.Split(version, ".")[0]
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func minorVersionOf(version string) uint64 {
	extractedStr := strings.Split(version, ".")[1]
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func patchVersionOf(version string) uint64 {
	extractedStr := strings.Split(version, ".")[2]
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func compareVersions(ver1, ver2 string) int8 {
	major1 := majorVersionOf(ver1)
	minor1 := minorVersionOf(ver1)
	patch1 := patchVersionOf(ver1)

	major2 := majorVersionOf(ver2)
	minor2 := minorVersionOf(ver2)
	patch2 := patchVersionOf(ver2)

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
	majorVersion := majorVersionOf(version) + major
	minorVersion := minorVersionOf(version) + minor
	patchVersion := patchVersionOf(version) + patch

	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion)
}

func tryCreateLogFile() bool {
	/*
		if _, err := os.Stat(LOG_PATH); os.IsNotExist(err) {
			logWarning(true, "Log directory does not exist")
			logDebug(true, "Attempting to create log directory at ", LOG_PATH)
			err := os.MkdirAll(filepath.Dir(LOG_PATH), 0755)

			if err != nil {
				logError(true, "Failed to create log directory: ", err)
				logWarning(true, "Defaulting to no file logging")
				SetLogToFile(false)
				return false

			} else {
				logDebug(true, "Log directory created at ", LOG_PATH)
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
				logDebug(true, "Log file created at ", LOG_FULL_PATH)
				return true
			}
		}

		Debug("Log file already exists at ", LOG_FULL_PATH)
	*/
	return true
}

func registerLogAttempt(forceNoFileWrite bool) bool {
	get().logAttempts++
	/*
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
		}*/
	return true
}
func resetLogAttempts(forceNoFileWrite bool) bool {
	if get().logAttempts == 1 {
		get().logAttempts = 0
		return true

	} else if get().logAttempts > 1 {
		get().logAttempts = 0
		logDebug(forceNoFileWrite, "Log attempts reset (", get().logAttempts, ")")
		return true
	}

	return false
}

func LogsToFile() bool {
	return get().LogToFile
}
func SetLogToFile(value bool) {
	logDebug(!value, "File logging set to ", value)
	get().LogToFile = value
}

func LogsToConsole() bool {
	return get().LogToConsole
}
func SetLogToConsole(value bool) {
	Debug("Console logging set to ", value)
	get().LogToConsole = value
}

func LogsWithColor() bool {
	return get().ColorLogs
}
func SetLogWithColor(value bool) {
	if value {
		Debug("Color logging enabled")
	} else {
		Debug("Color logging disabled")
	}

	get().ColorLogs = value
}

func SetAppVersion(version string) {
	regex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	if !regex.MatchString(version) {
		Warning("Invalid app version format")
		return
	}

	get().appVersion = version
	get().appVersionMajor = majorVersionOf(version)
	get().appVersionMinor = minorVersionOf(version)
	get().appVersionPatch = patchVersionOf(version)
	Debug("App version set to: ", version)
}

// Returns the application version is supposed to be the current
func AppVersion() string {
	return get().appVersion
}

func currentLogLevels() string {
	return "| " + LogLevels().Name() + " |"
}

// Returns the current logging levels loaded as bitmask
func LogLevels() logLevel {
	return get().LogLevels
}

func SetLogLevels(options logLevel) bool {
	msg := "Logging options set to:"

	if !options.IsValid() {
		Warning("Invalid logging options")
		return false
	} else if options.Is(Level.All()) {
		Debug(msg, Level.All().Name())
		get().LogLevels = Level.All()
		return true
	} else if options.Is(Level.None()) {
		Debug(msg, Level.None().Name())
		get().LogLevels = Level.None()
		return true
	}

	if LogLevels().Has(Level.Debug()) {
		msg += "| DEBUG |"
	}
	if LogLevels().Has(Level.Info()) {
		msg += "| INFO |"
	}
	if LogLevels().Has(Level.Warning()) {
		msg += "| WARNING |"
	}
	if LogLevels().Has(Level.Error()) {
		msg += "| ERROR |"
	}
	if LogLevels().Has(Level.Fatal()) {
		msg += "| FATAL |"
	}

	Debug("Logging options set to: ", msg)
	get().LogLevels = options
	return true
}
func EnableLogOptions(options logLevel) bool {
	if !options.IsValid() {
		Warning("Invalid logging option")
		return false
	} else if options.Is(Level.All()) {
		Debug("Enabled all logging options")
		SetLogLevels(Level.All())
		return true
	} else if options.Is(Level.None()) {
		Debug("Disabled all logging options")
		SetLogLevels(Level.None())
		return true
	}

	var msg string

	if options.Has(Level.Debug()) {
		msg += "| DEBUG |"
	}
	if options.Has(Level.Info()) {
		msg += "| INFO |"
	}
	if options.Has(Level.Warning()) {
		msg += "| WARNING |"
	}
	if options.Has(Level.Error()) {
		msg += "| ERROR |"
	}
	if options.Has(Level.Fatal()) {
		msg += "| FATAL |"
	}

	Debug("Enabled logging options: ", msg)
	get().LogLevels = LogLevels().And(options)
	return true
}
func DisableLogOptions(options logLevel) bool {
	if !options.IsValid() {
		Warning("Invalid logging option")
		return false
	}

	var msg string

	if options.Has(Level.Debug()) {
		msg += "| DEBUG |"
	}
	if options.Has(Level.Info()) {
		msg += "| INFO |"
	}
	if options.Has(Level.Warning()) {
		msg += "| WARNING |"
	}
	if options.Has(Level.Error()) {
		msg += "| ERROR |"
	}
	if options.Has(Level.Fatal()) {
		msg += "| FATAL |"
	}

	Debug("Disabled logging options: ", msg)
	get().LogLevels = LogLevels().Not(options)
	return true
}
func SetMinLogLevel(level logLevel) bool {
	if !level.IsValid() {
		Warning("Invalid logging level")
		return false
	}

	Debug("Minimum logging level set to: ", level.Name())
	SetLogLevels(level)
	return true
}
func canLogWith(logger *log.Logger) bool {
	if LogLevels().Has(Level.All()) {
		return true
	} else if LogLevels().Has(Level.None()) {
		return false
	}

	if logger == get().DebugLogger {
		return LogLevels().Has(Level.Debug())
	} else if logger == get().InfoLogger {
		return LogLevels().Has(Level.Info())
	} else if logger == get().WarningLogger {
		return LogLevels().Has(Level.Warning())
	} else if logger == get().ErrorLogger {
		return LogLevels().Has(Level.Error())
	} else if logger == get().FatalLogger {
		return LogLevels().Has(Level.Fatal())
	}

	return false
}

func writeToFile(logger *log.Logger, prefix string, v ...any) bool {
	file, err := os.OpenFile("get()", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

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
	internalLogWith(logger, forceNoFileWrite, 1, v...)
}
func internalLogWith(logger *log.Logger, forceNoFileWrite bool, extraTraceDepth int, v ...any) {
	if !canLogWith(logger) {
		return
	}

	registerLogAttempt(forceNoFileWrite)

	orignalPrefix := logger.Prefix()
	filePrefix, _ := utils.CallOrigin(4 + extraTraceDepth)
	filePrefix = Ansi.File().Apply(filePrefix)
	filePrefix += ":"

	stringValues := utils.AsStrings(v)
	stringValues = utils.Podate(stringValues, "[ ]")
	trimmedArgs := strings.Join(stringValues, " ")
	trimmedArgs = strings.Trim(trimmedArgs, "[]")

	if LogsToConsole() {
		logger.Println(filePrefix, trimmedArgs)
	}

	// This might be short-circuited, but idk if it is, so I will leave it like this
	if !forceNoFileWrite && LogsToFile() {
		if !writeToFile(logger, filePrefix, trimmedArgs) {
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

func iLogDebug(forceNoFileWrite bool, extraFileDepth int, v ...any) {
	internalLogWith(get().DebugLogger, forceNoFileWrite, extraFileDepth, v...)
}
func iLogInfo(forceNoFileWrite bool, extraFileDepth int, v ...any) {
	internalLogWith(get().InfoLogger, forceNoFileWrite, extraFileDepth, v...)
}
func iLogWarning(forceNoFileWrite bool, extraFileDepth int, v ...any) {
	internalLogWith(get().WarningLogger, forceNoFileWrite, extraFileDepth, v...)
}
func iLogError(forceNoFileWrite bool, extraFileDepth int, v ...any) {
	internalLogWith(get().ErrorLogger, forceNoFileWrite, extraFileDepth, v...)
}
func iLogFatal(forceNoFileWrite bool, extraFileDepth int, v ...any) {
	internalLogWith(get().FatalLogger, forceNoFileWrite, extraFileDepth, v...)
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
	deprecateTxt := Ansi.Deprecate().Apply(" DEPRECATED:")
	reasonTxt := Ansi.DeprecateReason().Apply(" REASON:")

	if compareVersions(AppVersion(), addToVersion(removalVersion, 0, 1, 0)) >= 0 {
		iLogFatal(false, 1, deprecateTxt, "This feature has been removed in version", removalVersion)
		iLogFatal(false, 1, reasonTxt, args)

	} else if compareVersions(AppVersion(), removalVersion) >= 0 {
		iLogError(false, 1, deprecateTxt+" This feature has been removed in version ", removalVersion)
		iLogError(false, 1, reasonTxt, args)
		return false, fmt.Errorf("feature removed in version %s. %s", removalVersion, args)

	} else if compareVersions(AppVersion(), deprecatedVersion) >= 0 && compareVersions(AppVersion(), removalVersion) < 0 {
		iLogWarning(false, 1, deprecateTxt, "This feature marked for removal in version", removalVersion)
		iLogWarning(false, 1, reasonTxt, args)
		return false, fmt.Errorf("feature deprecated in version %s, will be removed in version %s. %s", deprecatedVersion, removalVersion, args)

	} else if compareVersions(AppVersion(), deprecatedVersion) == 0 {
		iLogWarning(false, 1, deprecateTxt, "This feature will be removed in future versions")
		iLogWarning(false, 1, reasonTxt, args)
		return true, fmt.Errorf("feature deprecated in version %s. %s", deprecatedVersion, args)
	}
	return true, nil
}
func DeprecateMsg(deprecatedVersion string, removalVersion string, v ...any) string {
	_, err := Deprecate(deprecatedVersion, removalVersion, v...)
	return err.Error()
}

type volcano struct {
	currentVersion string
	coldVersion    string
	driedVersion   string

	lavaTxt      string
	coldLavaTxt  string
	driedLavaTxt string
}

func Lava(version string, v ...any) volcano {
	args := utils.Join(v, " ")
	lavaTxt := Ansi.Lava().Apply(" LAVA:")
	coldLavaTxt := Ansi.ColdLava().Apply(" COLD LAVA:")
	driedLavaTxt := Ansi.DriedLava().Apply(" DRIED LAVA:")

	coldVersion := addToVersion(version, 0, 0, 2)
	driedVersion := addToVersion(version, 0, 0, 4)

	if compareVersions(AppVersion(), version) == 0 {
		iLogWarning(false, 1, lavaTxt, "Running code that should be removed, cleaned up or refactored.", args)
	} else if compareVersions(AppVersion(), coldVersion) >= 0 {
		iLogWarning(false, 1, coldLavaTxt, "This code must be refactored.", args)
	} else if compareVersions(AppVersion(), driedVersion) > 0 {
		iLogError(false, 1, driedLavaTxt, "This code should not be running as is, it is likely a bug.", args)
	}

	return volcano{
		currentVersion: version,
		coldVersion:    coldVersion,
		driedVersion:   driedVersion,

		lavaTxt:      lavaTxt,
		coldLavaTxt:  coldLavaTxt,
		driedLavaTxt: driedLavaTxt,
	}
}
func (v *volcano) LavaStart() {
	if AppVersion() == v.currentVersion {
		iLogWarning(false, 1, v.lavaTxt, "Start of flow")
	} else if compareVersions(AppVersion(), v.coldVersion) >= 0 {
		iLogWarning(false, 1, v.coldLavaTxt, "Start of flow")
	} else if compareVersions(AppVersion(), v.driedVersion) > 0 {
		iLogError(false, 1, v.driedLavaTxt, "Start of lava cast")
	}
}
func (v *volcano) LavaEnd() {
	if AppVersion() == v.currentVersion {
		iLogWarning(false, 1, v.lavaTxt, "End of flow")
	} else if compareVersions(AppVersion(), v.coldVersion) >= 0 {
		iLogWarning(false, 1, v.coldLavaTxt, "End of flow")
	} else if compareVersions(AppVersion(), v.driedVersion) > 0 {
		iLogError(false, 1, v.driedLavaTxt, "End of lava cast")
	}
}
