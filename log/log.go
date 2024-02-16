package log

import (
	"fmt"
	"time"
)

func WithTimeStamp(msg string) {
	t := time.Now()
	timeTag := t.Format("15:04:05")
	dateTag := t.Format("2006-01-02")
	fmt.Printf("[%s/%s] %s\n", timeTag, dateTag, msg)
}
