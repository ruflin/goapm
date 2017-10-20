package goapm

import (
	"fmt"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	apm := NewAPM("test", "123")

	tr := apm.StartTransaction()
	time.Sleep(200 * time.Microsecond)
	tr.Stop()

	tr2 := apm.StartTransaction()
	time.Sleep(200 * time.Microsecond)
	tr2.Stop()

	apm.send()
}

func TestAppToMapStr(t *testing.T) {
	app := app{
		Name:    "test",
		Version: "123",
	}

	fmt.Printf("%+v", app.toMapStr())
}
