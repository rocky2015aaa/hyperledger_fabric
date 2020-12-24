package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/m/common/const"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
	"time"

	"net/http"
	"net/url"
	"strings"
)

func SetLoggerConfig(e *echo.Echo, logMode string) {
	fileIoWriterForSvc := &lumberjack.Logger{
		Filename:   os.Getenv("GOPATH") + const_rest.LOG_LOCATION + "svc.log",
		MaxSize:    100,
		MaxBackups: 50,
		MaxAge:     30,
		Compress:   true, // disabled by default
	}

	fileIoWriterForAcc := &lumberjack.Logger{
		Filename:   os.Getenv("GOPATH") + const_rest.LOG_LOCATION + "acc.log",
		MaxSize:    100,
		MaxBackups: 100,
		MaxAge:     60,
		Compress:   true, // disabled by default
	}

	var ioWriter1 io.Writer
	var ioWriter2 io.Writer

	if logMode == "dev" {
		e.Logger.SetLevel(log.DEBUG) // 1
		ioWriter1 = io.MultiWriter(os.Stdout, fileIoWriterForSvc)

		ioWriter2 = io.MultiWriter(os.Stdout, fileIoWriterForAcc)

	} else {
		log.SetLevel(log.INFO) // 2
		e.Logger.SetLevel(log.INFO)

		ioWriter1 = fileIoWriterForSvc
		ioWriter2 = fileIoWriterForAcc

	}

	e.Logger.SetOutput(ioWriter1)

	acclogger.SetLevel(log.INFO)
	acclogger.SetOutput(ioWriter2)
	e.Use(accessLogger(acclogger, e))

}

type test_struct struct {
	Key, Val string
}

var ACC_CNT = 0

func accessLogger(l *log.Logger, e *echo.Echo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			id := fmt.Sprintf("%v", start.UnixNano())
			req := c.Request()
			refID := c.FormValue("refID")
			if refID != "" {
				id = id + ":" + refID
			}
			req.Header.Set("req-key", id)
			res := c.Response()

			path := c.ParamNames()
			if len(path) > 0 {

				if (req.Method == echo.PUT || req.Method == echo.POST) && req.ContentLength == 0 {

					f := make(url.Values)
					for _, p := range path {

						f.Set(p, c.Param(p))
					}

					req2, _ := http.NewRequest(req.Method, req.URL.String(), strings.NewReader(f.Encode()))
					req2.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
					c.SetRequest(req2)

				} else if req.Method == echo.GET {
					q := c.Request().URL.Query()

					for _, p := range path {
						q.Set(p, c.Param(p))

					}

					c.Request().URL.RawQuery = q.Encode()
				}

			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()
			ACC_CNT += 1

			accInfo := []interface{}{"req-key=" + id, "status=" + strconv.Itoa(res.Status), "uri=" + req.RequestURI, "method=" + req.Method,
				"host=" + req.Host, "remote_ip=" + c.RealIP(), "user_agent=" + req.UserAgent(), "ACC_CNT=" + strconv.Itoa(ACC_CNT), "latency_human=" + stop.Sub(start).String()}

			l.Info(accInfo)
			return nil
		}
	}
}
