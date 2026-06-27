package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	logFile *os.File
	logOnce sync.Once
)

// Log writes a timestamped message to /tmp/peruzzi.log.
// It is safe to call from multiple goroutines/packages.
func Log(format string, args ...interface{}) {
	logOnce.Do(func() {
		var err error
		logFile, err = os.OpenFile("/tmp/peruzzi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return
		}
	})
	if logFile == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	logFile.WriteString(time.Now().Format("15:04:05.000") + " " + msg + "\n")
	logFile.Sync()
}
