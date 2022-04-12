package helper

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoPanicHandling(t *testing.T) {
	var ch = make(chan bool)
	Go("anonymous_function_that_panics", func() {
		ch <- true
		panic(errors.New("TEST ERROR in TestGoPanicHandling. It should not fail"))
	})
	<-ch
}

func TestPanicRecovery(t *testing.T) {
	a, b, err2 := myFunc()
	assert.Equal(t, 1, a)
	assert.Equal(t, 0, b)
	assert.NotEmpty(t, err2)
	assert.Contains(t, err2.Error(), "CHECK")
}

func myFunc() (a, b int, err error) {
	defer RecoverToErrorVar("myFunc", &err)
	a = 1
	if a == 1 {
		err2 := errors.New("TEST ERROR in myFunc. Code = CHECK")
		ErrorHandlerf(err2, "myFunc", "%s", "arg")
	}
	b = 1
	return
}

//func TestRecoverAsLogErrorf(t *testing.T) {
//	deferRecoverAsLogErrorf("%s", "arg")
//	err2 := errors.New("TEST ERROR in TestRecoverAsLogErrorf: should not fail")
//	//(err2, "TestRecoverAsLogErrorf", "%s", "arg")
//}
