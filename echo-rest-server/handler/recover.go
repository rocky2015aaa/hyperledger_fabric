package handler

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo"
)

type (
	RecoverConfig struct {
		StackSize int `yaml:"stack_size"`

		DisableStackAll bool `yaml:"disable_stack_all"`

		DisablePrintStack bool `yaml:"disable_print_stack"`
	}
)

var (
	DefaultRecoverConfig = RecoverConfig{
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
)

func SetRecoverConfig() echo.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {

	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msgKey := "req-key=" + c.Request().Header.Get("req-key") + "-"
						c.Logger().Printf(msgKey+"PANIC RECOVER = %v %s\n", err, stack[:length])
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
