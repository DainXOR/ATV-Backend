package logger

import (
	"dainxor/atv/types"
	"dainxor/atv/utils"
	"errors"
	"strconv"
	"strings"

	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

type dnxLogger struct {
	configurations

	logAttempts   uint8 // Number of log attempts, used to prevent infinite loops in logging
	usingDefaults bool  // Indicates if the logger is using default settings

	abnormalWriters []writerAndFormatter // Special writers used when normal logging is failing
}

var dnxGlobalLogger *dnxLogger

func init() {
	dnxGlobalLogger = newDefaultLogger()
}

func newDefaultLogger() *dnxLogger {
	logger := &dnxLogger{}
	fm := SimpleFormatter.New()

	logger = &dnxLogger{
		configurations: NewConfigs(),
		logAttempts:    0,
		usingDefaults:  true,

		abnormalWriters: []writerAndFormatter{
			{writer: ConsoleWriter.NewLine().New(), formatter: &fm},
		},
	}

	envLoad(logger) // Initialize environment variables for logger

	tryCreateLogFile() // Create the path if it doesn't exist, else it will set log to console only

	logger.Debug("Logger initialized")

	return logger
}
func cloneLogger(original *dnxLogger) *dnxLogger {
	logger := &dnxLogger{
		configurations:  original.configurations,
		logAttempts:     0,
		usingDefaults:   false,
		abnormalWriters: make([]writerAndFormatter, len(original.abnormalWriters)),
	}

	copy(logger.abnormalWriters, original.abnormalWriters)
	logger.Debug("Logger cloned")
	return logger
}
func envLoad(logger *dnxLogger) {
	logger.Debug("Loading environment variables for logger")

	//minLogLevel, existMinLevel := os.LookupEnv("DNX_LOG_MIN_LEVEL")
	//disableLevels, existDisableLevels := os.LookupEnv("DNX_LOG_DISABLE_LEVELS")
	//logConsole, existLogConsole := os.LookupEnv("DNX_LOG_CONSOLE")
	//logFile, existLogFile := os.LookupEnv("DNX_LOG_FILE")
	//logWithColor, existLogWithColor := os.LookupEnv("DNX_LOG_WITH_COLOR")

	logger.Debug("Logger environment variables loaded")
}

// Returns the singleton instance of dnxLogger, initializing it if necessary.
// This function should be used to access the logger throughout the internal package.
// It abstracts the initialization logic and provides a single point of access to the logger instance.
// It ensures that the logger is initialized only once, and provides a consistent interface for logging.
func get() *dnxLogger {
	if dnxGlobalLogger == nil {
		dnxGlobalLogger = newDefaultLogger()
	}
	return dnxGlobalLogger
}

func GetNew() *dnxLogger {
	return cloneLogger(get())
}
func GetNewDefault() *dnxLogger {
	return newDefaultLogger()
}
func GetNewWithEnv() *dnxLogger {
	logger := GetNewDefault()
	envLoad(logger)
	return logger
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
func panicOrError(msg string, condition bool) error {
	if condition {
		panic(msg)
	} else {
		return errors.New(msg)
	}
}
func (i *dnxLogger) logLevelsNames() string {
	contained := Level.ContainedIn(i.LogLevels())
	names := utils.Map(contained, logLevel.Name)
	return "| " + strings.Join(names, " | ") + " |"
}

func (i *dnxLogger) registerLogAttempt() bool {
	i.logAttempts++

	if i.logAttempts > i.warningLogAttempts && i.logAttempts < i.maxLogAttempts {
		i.internalAbnormalWrite(NewRecord("Too many log attempts, will panic if this continues"))
		return false

	} else if i.logAttempts >= i.maxLogAttempts {
		i.internalAbnormalWrite(NewRecord("Too many log attempts, will panic now"))
		i.internalAbnormalWrite(NewRecord("Too many log attempts, this is likely a bug in the logger, please report it"))

		return false
	}
	return true
}
func (i *dnxLogger) resetLogAttempts() bool {
	if i.logAttempts == 1 {
		i.logAttempts = 0
		return true

	} else if i.logAttempts > 1 {
		i.logAttempts = 0
		i.internalAbnormalWrite(NewRecord("Resetting log attempts after too many attempts"))
		return true
	}

	return false
}

func (i *dnxLogger) SetVersion(version types.Version) {
	i.appVersion = version
	i.Debug("App version changed to: ", version.String())
}
func SetVersion(version types.Version) {
	get().SetVersion(version)
}

// Returns the application version used as reference for deprecation warnings.
func (i *dnxLogger) AppVersion() types.Version {
	return i.appVersion
}
func AppVersion() types.Version {
	return get().AppVersion()
}

// Returns the current logging levels loaded as bitmask
func (i *dnxLogger) LogLevels() logLevel {
	return i.logLevels
}

func (i *dnxLogger) SetLogLevels(options logLevel) bool {
	msg := "Logging options set to:"

	if !options.IsValid() {
		i.Warning("Invalid logging options")
		return false
	} else if options.Is(Level.All()) {
		i.Debug(msg, Level.All().Name())
		i.logLevels = Level.All()
		return true
	} else if options.Is(Level.None()) {
		i.Debug(msg, Level.None().Name())
		i.logLevels = Level.None()
		return true
	}

	if i.LogLevels().Has(Level.Debug()) {
		msg += "| DEBUG |"
	}
	if i.LogLevels().Has(Level.Info()) {
		msg += "| INFO |"
	}
	if i.LogLevels().Has(Level.Warning()) {
		msg += "| WARNING |"
	}
	if i.LogLevels().Has(Level.Error()) {
		msg += "| ERROR |"
	}
	if i.LogLevels().Has(Level.Fatal()) {
		msg += "| FATAL |"
	}

	i.Debug("Logging options set to: ", msg)
	i.logLevels = options
	return true
}
func (i *dnxLogger) EnableLogOptions(options logLevel) bool {
	if !options.IsValid() {
		i.Warning("Invalid logging option")
		return false
	} else if options.Is(Level.All()) {
		i.Debug("Enabled all logging options")
		i.SetLogLevels(Level.All())
		return true
	} else if options.Is(Level.None()) {
		i.Debug("Disabled all logging options")
		i.SetLogLevels(Level.None())
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

	i.Debug("Enabled logging options: ", msg)
	i.logLevels.Set(i.LogLevels().And(options))
	return true
}
func (i *dnxLogger) DisableLogOptions(options logLevel) bool {
	if !options.IsValid() {
		i.Warning("Invalid logging option")
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

	i.Debug("Disabled logging options: ", msg)
	i.logLevels.Set(i.LogLevels().Not(options))
	return true
}
func (i *dnxLogger) SetMinLogLevel(level logLevel) bool {
	if !level.IsValid() {
		i.Warning("Invalid logging level")
		return false
	}

	i.Debug("Minimum logging level set to: ", level.Name())
	i.SetLogLevels(level)
	return true
}
func (i *dnxLogger) canLog(level logLevel) bool {
	return i.LogLevels().Has(level)
}

func (i *dnxLogger) internalWrite(record Record) {
	if !i.canLog(record.LogLevel) {
		return
	}

	i.registerLogAttempt()
	removeList := make([]int, 0, len(i.writers))

	for p, pair := range i.writers {
		w := pair.writer
		f := *pair.formatter

		if w == nil || f == nil {
			record := NewRecord("Invalid writer or formatter in writer list, marking for removal")
			record.LogLevel = Level.Error()
			i.internalAbnormalWrite(record)
			removeList = append(removeList, p)
			continue
		}

		if str, err := f.Format(&record); err == nil {
			if err = w.Write(str); err != nil {
				record := NewRecord("Error during write: " + err.Error() + ", marking writer for removal")
				record.LogLevel = Level.Error()
				i.internalAbnormalWrite(record)
				removeList = append(removeList, p)
			}

		} else {
			panic("Error formatting record for abnormal write: " + err.Error())
		}
	}

	for _, index := range removeList {
		i.RemoveWriter(index)
	}

	i.resetLogAttempts()
}
func (i *dnxLogger) internalAbnormalWrite(record Record) error {
	if !i.canLog(record.LogLevel) {
		return nil
	}

	for _, wf := range get().abnormalWriters {
		if wf.writer == nil || wf.formatter == nil {
			return panicOrError("Invalid writer or formatter in abnormal write", i.canPanicOnAbnormalWrite)
		}

		w, f := wf.writer, (*wf.formatter)
		if str, err := f.Format(&record); err == nil {
			if err = w.Write(str); err != nil {
				return panicOrError("Error during abnormal write: "+err.Error(), i.canPanicOnAbnormalWrite)
			}

		} else {
			return panicOrError("Error formatting record for abnormal write: "+err.Error(), i.canPanicOnAbnormalWrite)
		}
	}

	return nil
}

// Public functions for logging at different levels
func generateRecord(level logLevel, v ...any) Record {
	msg := ""
	extraParts := []types.SPair[string]{types.NewSPair("_internal_call_origin_offset", "2")}

	for _, pair := range utils.ValuesOfType[types.SPair[string]](v) {
		if internal.CallOriginOffset().Check(pair.First) {
			offset, _ := strconv.Atoi(pair.Second)
			extraParts[0].Second = fmt.Sprint(offset + 2)
			continue
		} else if internal.FormatString().Check(pair.First) {
			msg = fmt.Sprintf(pair.Second, utils.ExcludeOfType[types.SPair[string]](v)...)
			continue
		}

		extraParts = append(extraParts, pair)
	}

	if msg == "" {
		msg = fmt.Sprint(utils.ExcludeOfType[types.SPair[string]](v)...)
	}

	record := NewRecord(msg, extraParts...)
	record.LogLevel = level
	return record
}

func (i *dnxLogger) iWrite(l logLevel, v ...any) {
	record := generateRecord(l, v...)
	i.internalWrite(record)
}
func (i *dnxLogger) iWritef(l logLevel, format string, v ...any) {
	extraParts := []any{types.NewSPair(internal.FormatString().String(), format)}
	extraParts = append(extraParts, v...)

	record := generateRecord(l, extraParts...)
	i.internalWrite(record)
}
func (i *dnxLogger) iWriter(r Record) {
	i.internalWrite(r)
}

func (i *dnxLogger) Debug(v ...any) {
	i.iWrite(Level.Debug(), v...)
}
func (i *dnxLogger) Debugf(format string, v ...any) {
	i.iWritef(Level.Debug(), format, v...)
}
func (i *dnxLogger) Debugr(r Record) {
	r.LogLevel = Level.Debug()
	i.iWriter(r)
}
func Debug(v ...any) {
	get().Debug(v...)
}
func Debugf(format string, v ...any) {
	get().Debugf(format, v...)
}
func Debugr(r Record) {
	get().Debugr(r)
}

func (i *dnxLogger) Info(v ...any) {
	i.iWrite(Level.Info(), v...)
}
func (i *dnxLogger) Infof(format string, v ...any) {
	i.iWritef(Level.Info(), format, v...)
}
func (i *dnxLogger) Infor(r Record) {
	r.LogLevel = Level.Info()
	i.iWriter(r)
}
func Info(v ...any) {
	get().Info(v...)
}
func Infof(format string, v ...any) {
	get().Infof(format, v...)
}
func Infor(r Record) {
	get().Infor(r)
}

func (i *dnxLogger) Warning(v ...any) {
	i.iWrite(Level.Warning(), v...)
}
func (i *dnxLogger) Warningf(format string, v ...any) {
	i.iWritef(Level.Warning(), format, v...)
}
func (i *dnxLogger) Warningr(r Record) {
	r.LogLevel = Level.Warning()
	i.iWriter(r)
}
func Warning(v ...any) {
	get().Warning(v...)
}
func Warningf(format string, v ...any) {
	get().Warningf(format, v...)
}
func Warningr(r Record) {
	get().Warningr(r)
}

func (i *dnxLogger) Error(v ...any) {
	i.iWrite(Level.Error(), v...)
}
func (i *dnxLogger) Errorf(format string, v ...any) {
	i.iWritef(Level.Error(), format, v...)
}
func (i *dnxLogger) Errorr(r Record) {
	r.LogLevel = Level.Error()
	i.iWriter(r)
}
func Error(v ...any) {
	get().Error(v...)
}
func Errorf(format string, v ...any) {
	get().Errorf(format, v...)
}
func Errorr(r Record) {
	get().Errorr(r)
}

func (i *dnxLogger) Fatal(v ...any) {
	i.iWrite(Level.Fatal(), v...)
}
func (i *dnxLogger) Fatalf(format string, v ...any) {
	i.iWritef(Level.Fatal(), format, v...)
}
func (i *dnxLogger) Fatalr(r Record) {
	r.LogLevel = Level.Fatal()
	i.iWriter(r)
}
func Fatal(v ...any) {
	get().Fatal(v...)
}
func Fatalf(format string, v ...any) {
	get().Fatalf(format, v...)
}
func Fatalr(r Record) {
	get().Fatalr(r)
}

func (i *dnxLogger) Deprecate(deprecatedVersion types.Version, removalVersion types.Version, v ...any) (bool, error) {
	args := utils.Join(v, " ")

	if i.AppVersion().GreaterEq(removalVersion.Plus(types.V("0.1.0"))) {
		i.iWrite(Level.DeprecateFatal(), "This feature has been removed in version:", removalVersion)
		i.iWrite(Level.DeprecateFatal(), "Reason:", args)
		i.iWrite(Level.DeprecateFatal(), "Unrecoverable deprecation error, application will terminate")
		panic(TerminationCode.Deprecated())

	} else if i.AppVersion().GreaterEq(removalVersion) {
		i.iWrite(Level.DeprecateError(), "This feature has been removed in version:", removalVersion)
		i.iWrite(Level.DeprecateError(), "Reason:", args)
		return false, fmt.Errorf("feature removed in version %s. %s", removalVersion, args)

	} else if i.AppVersion().GreaterEq(deprecatedVersion) && i.AppVersion().LessThan(removalVersion) {
		i.iWrite(Level.DeprecateWarning(), "This feature is marked for removal in version:", removalVersion)
		i.iWrite(Level.DeprecateWarning(), "Reason:", args)
		return false, fmt.Errorf("feature deprecated in version %s, will be removed in version %s. %s", deprecatedVersion, removalVersion, args)

	} else if i.AppVersion().GreaterEq(deprecatedVersion) {
		i.iWrite(Level.Deprecate(), "This feature will be removed in future versions")
		i.iWrite(Level.Deprecate(), "Reason:", args)
		return true, fmt.Errorf("feature deprecated in version %s. %s", deprecatedVersion, args)
	}
	return true, nil
}
func (i *dnxLogger) DeprecateMsg(deprecatedVersion types.Version, removalVersion types.Version, v ...any) string {
	_, err := i.Deprecate(deprecatedVersion, removalVersion, v...)
	return err.Error()
}
func Deprecate(deprecatedVersion types.Version, removalVersion types.Version, v ...any) (bool, error) {
	return get().Deprecate(deprecatedVersion, removalVersion, v...)
}
func DeprecateMsg(deprecatedVersion types.Version, removalVersion types.Version, v ...any) string {
	return get().DeprecateMsg(deprecatedVersion, removalVersion, v...)
}

type volcano struct {
	logger         *dnxLogger
	initialVersion types.Version
	coldVersion    types.Version
	driedVersion   types.Version
}

func (i *dnxLogger) Lava(version types.Version, v ...any) volcano {
	args := utils.Join(v, " ")

	coldVersion := version.Plus(types.V("0.1.0"))
	driedVersion := version.Plus(types.V("0.2.0"))

	if i.AppVersion().GreaterEq(driedVersion) {
		i.iWrite(Level.LavaDry(), "Running code that should be removed, cleaned up or refactored.")
		i.iWrite(Level.LavaDry(), "Reason:", args)
		//iLogWarning(false, 1, lavaTxt, "Running code that should be removed, cleaned up or refactored.", args)
	} else if i.AppVersion().GreaterEq(coldVersion) {
		i.iWrite(Level.LavaCold(), "This code must be refactored.")
		i.iWrite(Level.LavaCold(), "Reason:", args)
		//iLogWarning(false, 1, coldLavaTxt, "This code must be refactored.", args)
	} else if i.AppVersion().GreaterEq(version) {
		i.iWrite(Level.LavaHot(), "This code should not be running as is, it is likely a bug.")
		i.iWrite(Level.LavaHot(), "Reason:", args)
		//iLogError(false, 1, driedLavaTxt, "This code should not be running as is, it is likely a bug.", args)
	}

	return volcano{
		logger:         i,
		initialVersion: version,
		coldVersion:    coldVersion,
		driedVersion:   driedVersion,
	}
}
func Lava(version types.Version, v ...any) volcano {
	return get().Lava(version, v...)
}
func (v *volcano) LavaStart() {
	if v.logger.AppVersion().GreaterEq(v.driedVersion) {
		v.logger.iWrite(Level.LavaDry(), "Start of lava cast")
	} else if v.logger.AppVersion().GreaterEq(v.coldVersion) {
		v.logger.iWrite(Level.LavaCold(), "Start of flow")
	} else if v.logger.AppVersion().GreaterEq(v.initialVersion) {
		v.logger.iWrite(Level.LavaHot(), "Start of flow")
	}
}
func (v *volcano) LavaEnd() {
	if v.logger.AppVersion().GreaterEq(v.driedVersion) {
		v.logger.iWrite(Level.LavaDry(), "End of lava cast")
	} else if v.logger.AppVersion().GreaterEq(v.coldVersion) {
		v.logger.iWrite(Level.LavaCold(), "End of flow")
	} else if v.logger.AppVersion().GreaterEq(v.initialVersion) {
		v.logger.iWrite(Level.LavaHot(), "End of flow")
	}
}
