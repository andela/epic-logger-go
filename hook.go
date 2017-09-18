package epiclogger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ServiceHook struct{}

func NewServiceHook() *ServiceHook {
	return &ServiceHook{}
}

func (hook *ServiceHook) Fire(entry *log.Entry) error {
	podName := os.Getenv("POD_NAME")
	serviceNVersion := strings.Split(podName, "-")
	length := len(serviceNVersion)
	if podName != "" {
		entry.Data["service"] = strings.Join(serviceNVersion[:length-2], "-")
		entry.Data["version"] = serviceNVersion[length-2]
	}
	return nil
}

// Levels enumerates the log levels on which the error should be forwarded to
// bugsnag: everything at or above the "Error" level.
func (hook *ServiceHook) Levels() []log.Level {
	return []log.Level{
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
