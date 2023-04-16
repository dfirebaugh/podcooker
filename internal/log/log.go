package log

import (
	"fmt"
	"os"
)

type FileLogger struct {
	OutputFile string
}

func (l FileLogger) Println(v ...interface{}) {
	f, err := os.OpenFile(l.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("failed to open log file %s: %v", l.OutputFile, err)
		return
	}
	defer f.Close()

	fmt.Fprintln(f, v...)
}
