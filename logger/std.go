// +build !windows

package logger

import (
	"fmt"
	"log"
	"os"
)

func debugLog(msg string, fields ...Field) {
	// No debug logs by default
}

func infoLog(msg string, fields ...Field) {
	logWithPrefix("[INFO]", msg, fields...)
}

func errorLog(msg string, fields ...Field) {
	logWithPrefix("[ERROR]", msg, fields...)
}

func logWithPrefix(prefix, msg string, fields ...Field) {
	output := prefix + " " + msg
	if len(fields) > 0 {
		output += " " + formatFields(fields)
	}
	log.Println(output)
}

func formatFields(fields []Field) string {
	result := "{"
	for i, f := range fields {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%s=%v", f.Key, f.Value)
	}
	result += "}"
	return result
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
}
