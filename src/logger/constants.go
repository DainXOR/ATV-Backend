package logger

import (
	"fmt"
	"log"
)

type logLevel int

const ( // Log levels
	LEVEL_DEBUG   logLevel = 1 << iota // 0b0000_00000001
	LEVEL_INFO                         // 0b0000_00000010
	LEVEL_WARNING                      // 0b0000_00000100
	LEVEL_ERROR                        // 0b0000_00001000
	LEVEL_FATAL                        // 0b0000_00010000

	LEVEL_DEPRECATE         logLevel = (1 << 8)                        // 0b0001_00000000
	LEVEL_DEPRECATE_WARNING          = LEVEL_DEPRECATE | LEVEL_WARNING // 0b0001_00000100
	LEVEL_DEPRECATE_ERROR            = LEVEL_DEPRECATE | LEVEL_ERROR   // 0b0001_00001000
	LEVEL_DEPRECATE_FATAL            = LEVEL_DEPRECATE | LEVEL_FATAL   // 0b0001_00010000
	LEVEL_LAVA              logLevel = (1 << 9)                        // 0b0010_00000000
	LEVEL_LAVA_HOT                   = LEVEL_LAVA | LEVEL_DEBUG        // 0b0010_00000001
	LEVEL_LAVA_COLD                  = LEVEL_LAVA | LEVEL_WARNING      // 0b0010_00000100
	LEVEL_LAVA_DRY                   = LEVEL_LAVA | LEVEL_ERROR        // 0b0010_00001000

	LEVEL_ALL           = LEVEL_DEBUG | LEVEL_INFO | LEVEL_WARNING | LEVEL_ERROR | LEVEL_FATAL
	LEVEL_NONE logLevel = 0
)

func (l logLevel) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_WARNING:
		return "WARNING"
	case LEVEL_ERROR:
		return "ERROR"
	case LEVEL_FATAL:
		return "FATAL"
	case LEVEL_ALL:
		return "ALL"
	case LEVEL_NONE:
		return "NONE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", l)
	}
}

const ( // File log constants
	LOG_PATH      = "../artifacts/"
	LOG_FILE      = "logs.log"
	LOG_FULL_PATH = LOG_PATH + LOG_FILE
)

const ( // Default logger settings
	DEFAULT_LOGS_TO_FILE         = false
	DEFAULT_LOGS_TO_CONSOLE      = true
	DEFAULT_COLOR_LOGGING        = false
	DEFAULT_LOG_LEVEL            = LEVEL_ALL
	ENABLE_LOG_ATTEMPTS_MESSAGES = true
	DEFAULT_MAX_LOG_ATTEMPTS     = 15
	DEFAULT_WARNING_LOG_ATTEMPTS = 10

	DEFAULT_LOGGER_FLAGS = log.Ldate | log.Ltime

	DEFAULT_APP_VERSION = "0.1.0"
)
