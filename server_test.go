package microBase

import "testing"

func TestMicroInit(t *testing.T) {
	service := MicroInit()
	service.Run()
}
