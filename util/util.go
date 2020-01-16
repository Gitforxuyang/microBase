package util

import (
	"fmt"
	"strings"
)

// Must make error panic.
func Must(err error, ctxinfo ...interface{}) {
	if err == nil {
		return
	}
	if len(ctxinfo) > 0 {
		info := []string{}
		for _, a := range ctxinfo { // XXX: fmt.Sprint is not good enough..
			info = append(info, fmt.Sprintf("%v", a))
		}
		panic(fmt.Errorf("%v: %+v", strings.Join(info, " "), err))
	} else {
		panic(err)
	}
}

// Must make error panic.
func MustNotNil(t interface{}, name string) {
	if t == nil {
		panic(fmt.Errorf("%s must not nil", name))
	}
}
