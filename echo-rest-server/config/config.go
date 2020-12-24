package config

import (
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
	//"github.com/m/echo-rest-server/logger"
	"github.com/m/common/const"
	"github.com/m/echo-rest-server/handler"
	//"fmt"
	//"net"
	//"net/http"
	//"github.com/pkg/errors"
	//	"io"
	//"strconv"
	//"time"
	//"os"
	//"github.com/labstack/gommon/log"
)

func Set() *echo.Echo {
	// Echo instance
	e := echo.New()
	/*
		e.Use(handler.CORSWithConfig(handler.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		}))
	*/
	e.HideBanner = false

	//e.AutoTLSManager.Cache = autocert.DirCache("/home/vagrant/workspace/go/src/github.com/m/chaincode/src/rest/CRT/.cache")

	handler.SetLoggerConfig(e, const_rest.LOGMODE)

	handler.SetErrorConfig(e)

	e.Use(handler.SetRecoverConfig())

	return e
}
