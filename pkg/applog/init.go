package applog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var LogLevel = os.Getenv("LOG_LEVEL")

type LogBuilder interface {
	WithAttribute(key string, value interface{}) LogBuilder
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type _builder struct {
	attributes map[string]interface{}
}

func (logBuilder *_builder) Info(format string, args ...interface{}) {
	logBuilder.log("INFO", format, args[:]...)
}

func (logBuilder *_builder) Error(format string, args ...interface{}) {
	logBuilder.log("ERROR", format, args[:]...)
}

func (logBuilder *_builder) log(level string, format string, args ...interface{}) {
	logBuilder.attributes["level"] = level
	logBuilder.attributes["time"] = time.Now().Format(time.RFC3339)
	logBuilder.attributes["message"] = fmt.Sprintf(format, args[:]...)
	str := bytes.Buffer{}
	json.NewEncoder(&str).Encode(logBuilder.attributes)
	fmt.Print(str.String())
}

func (logBuilder *_builder) WithAttribute(key string, value interface{}) LogBuilder {
	logBuilder.attributes[key] = value
	return logBuilder
}

type Logger interface {
	Tag(tag string) LogBuilder
}

type _logger struct {
}

func (log *_logger) tag(tag string) LogBuilder {
	return &_builder{
		attributes: map[string]interface{}{
			"tag": tag,
		},
	}
}

func IsDebugEnabled() bool {
	return LogLevel == "DEBUG"
}

var _log = &_logger{}

func Tag(tag string) LogBuilder {
	return _log.tag(tag)
}
