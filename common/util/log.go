package util

import (
	logger "gitlab.altex.ro/plug/go_logger"
	"log"
	"os"
)

type LoggerAware struct {
	logger logger.Logger
}

/**
 * SetLogger
 * @param logger *log.Logger
 */
func (la *LoggerAware) SetLogger(logger logger.Logger) {
	la.logger = logger
}

/**
 * GetLogger - lazy load
 * @return *log.Logger
 */
func (la *LoggerAware) GetLogger() logger.Logger {
	if la.logger == nil {
		la.logger = logger.New(os.Stderr, "", log.LstdFlags)
	}

	return la.logger
}
