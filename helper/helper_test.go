package helper_test

import (
	"github.com/kraneware/kore-go/helper"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoPanicHandling(t *testing.T) {
	var ch = make(chan bool)
	helper.Go("anonymous_function_that_panics", func() {
		ch <- true
		panic(errors.New("TEST ERROR in TestGoPanicHandling. It should not fail"))
	})
	<-ch
}

var _ = Describe("Env Tests", func() {
	Context("Env Test", func() {
		It("should reset environment property", func() {
			environ := os.Environ()

			Expect(os.Setenv("TEST_KEY", "TEST_VALUE")).Should(BeNil())

			Expect(os.Getenv("TEST_KEY")).Should(Equal("TEST_VALUE"))

			Expect(helper.ResetEnv(environ)).Should(BeNil())

			Expect(os.Getenv("TEST_KEY")).ShouldNot(Equal("TEST_VALUE"))
		})
	})
})

func TestPanicRecovery(t *testing.T) {
	a, b, err2 := myFunc()
	assert.Equal(t, 1, a)
	assert.Equal(t, 0, b)
	assert.NotEmpty(t, err2)
	assert.Contains(t, err2.Error(), "CHECK")
}

func myFunc() (a, b int, err error) {
	defer helper.RecoverToErrorVar("myFunc", &err)
	a = 1
	if a == 1 {
		err2 := errors.New("TEST ERROR in myFunc. Code = CHECK")
		helper.ErrorHandlerf(err2, "myFunc", "%s", "arg")
	}
	b = 1
	return
}

func TestRecoverAsLogErrorf(t *testing.T) {
	defer helper.RecoverAsLogErrorf("%s", "arg")
	err2 := errors.New("TEST ERROR in TestRecoverAsLogErrorf: should not fail")
	helper.ErrorHandlerf(err2, "TestRecoverAsLogErrorf", "%s", "arg")
}
