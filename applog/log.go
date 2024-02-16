package applog

import (
	"fmt"
	"time"
)

func WithTimeStamp(msg string, values ...interface{}) {
	t := time.Now()
	timeTag := t.Format("15:04:05")
	dateTag := t.Format("2006-01-02")
	formattedMsg := fmt.Sprintf(msg, values...)
	fmt.Printf("[%s/%s] %s\n", timeTag, dateTag, formattedMsg)
}
