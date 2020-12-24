package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	//"github.com/pkg/errors"
	"github.com/m/common/msg"
	//"github.com/asaskevich/govalidator"
	"encoding/json"
	"strings"
)

var RETRY_TARGET = map[string]string{
	"MVCC_READ_CONFLICT":    "RETRY",
	"PHANTOM_READ_CONFLICT": "RETRY",
}

func SetErrorConfig(e *echo.Echo) {
	e.HTTPErrorHandler = customHTTPErrorHandler

}

func checkRetryTarget(err error) (int, string) {
	cnt := 0
	errStr := err.Error()

	for code, _ := range RETRY_TARGET {

		cnt = strings.Count(errStr, code)
		if cnt > 0 {
			return cnt, code
		}

	}

	return cnt, ""
}

func makeErrData(err error) interface{} {

	var message interface{}
	var errCode string

	errStr := err.Error()
	errAsBytes := []byte(errStr)
	commonSt := msg.CommonMsg{}
	json.Unmarshal(errAsBytes, &commonSt)

	if cnt, code := checkRetryTarget(err); cnt > 0 {

		errCode = code

	} else {
		errCode = commonSt.Code
	}

	message = msg.CommonMsg{errCode, commonSt.Desc, commonSt.Message}

	return message
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		code    = http.StatusInternalServerError
		desc    string
		message interface{}
	)

	// http web error or RESTFul error
	if he, ok := err.(*echo.HTTPError); ok {

		code = he.Code

		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
		desc = "HTTP PROTOCOL ERROR"
		message = msg.CommonMsg{strconv.Itoa(he.Code), desc, he.Message.(string)}

	} else { // RESTFul ERROR OR CC ERROR

		message = makeErrData(err)

	}

	fmt.Println("----------customHTTPErrorHandler message--------------------")
	fmt.Println(code)
	fmt.Println(message)

	err = c.JSON(code, message)
}
