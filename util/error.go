package util

import (
	"fmt"
	"github.com/micro/go-micro/errors"
	"net/http"
)

var (
	UnknowError error
	ParamError  error
	port        int
)

//端口号为5位  31001  错误码为4位0001-0999 为保留错误码，供框架层面使用  1001-9999为业务错误码。供每个应用使用
//0500-0599 为内部错误  参照http status code
func ErrInit(port int) {
	port = port
	UnknowError = errors.New(fmt.Sprintf("%d%s", port, "0500"), "未知错误", http.StatusInternalServerError)
	ParamError = errors.New(fmt.Sprintf("%d%s", port, "0400"), "参数错误", http.StatusBadRequest)
}

func NewUnknowError(detail string) error {
	return errors.New(fmt.Sprintf("%d%s", port, "0500"), detail, http.StatusInternalServerError)
}

func NewCustomError(id string, detail string) error {
	return errors.New(fmt.Sprintf("%d%s", port, id), detail, http.StatusInternalServerError)

}
