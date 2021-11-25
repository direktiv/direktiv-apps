package reusable

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

const devMode = "development"

type DirektivLogger struct {
	logger zerolog.Logger
}

type RequestInfo struct {
	aid    string
	logger *DirektivLogger
}

type DirektivLoggerWriter struct {
	aid string
}

func timestamp(in interface{}) string {
	// return time.Now().Format("15:04:05.000")
	return ""
}

func consoleWriter(w io.Writer) zerolog.ConsoleWriter {
	cw := zerolog.ConsoleWriter{Out: w}
	cw.NoColor = true
	cw.FormatTimestamp = timestamp
	cw.FormatLevel = func(i interface{}) string {
		return ""
	}
	return cw
}

func newRequestInfo(aid string) *RequestInfo {

	dl := &DirektivLoggerWriter{
		aid: aid,
	}
	cw := consoleWriter(dl)

	return &RequestInfo{
		aid: aid,
		logger: &DirektivLogger{
			logger: GetZeroLogger(cw),
		},
	}

}

func GetZeroLogger(w io.Writer) zerolog.Logger {

	// setup logger
	cw := consoleWriter(os.Stderr)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	var wr io.Writer = cw
	if w != nil {
		wr = io.MultiWriter(cw, w)
	}

	l := zerolog.New(wr).With().Timestamp().Logger()
	return l

}

func (ri *RequestInfo) ActionID() string {
	return ri.aid
}

func (ri *RequestInfo) Logger() *DirektivLogger {
	return ri.logger
}

func (dl *DirektivLogger) Errorf(format string, args ...interface{}) {
	txt := fmt.Sprintf(format, args...)
	dl.logger.Error().Msg(txt)
}

func (dl *DirektivLogger) Infof(format string, args ...interface{}) {
	txt := fmt.Sprintf(format, args...)
	dl.logger.Info().Msg(txt)
}

func (dl *DirektivLogger) Debugf(format string, args ...interface{}) {
	txt := fmt.Sprintf(format, args...)
	dl.logger.Debug().Msg(txt)
}

// Write writes log output
func (dl *DirektivLoggerWriter) Write(p []byte) (n int, err error) {

	if dl.aid != devMode {
		_, err = http.Post(fmt.Sprintf("http://localhost:8889/log?aid=%s", dl.aid), "plain/text", bytes.NewBuffer(p))
		return len(p), err
	}

	return len(p), nil
}