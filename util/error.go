package util

import (
	"fmt"
	"github.com/micro/go-micro/errors"
	"net/http"
)

var (
	UnknowError error
	port        int
)

//端口号为5位  31001  错误码为4位0001-0999 为保留错误码，供框架层面使用
//0500-0599 为内部错误  参照http status code
func ErrInit(port int) {
	port = port
	UnknowError = errors.New(fmt.Sprintf("%d%s", port, "0500"), "未知错误", http.StatusInternalServerError)
}

func NewUnkownError(detail string) error {
	return errors.New(fmt.Sprintf("%d%s", port, "0500"), detail, http.StatusInternalServerError)
}
