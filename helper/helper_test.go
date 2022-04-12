package helper_test

import (
	"github.com/kraneware/kore-go/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

type TestStruct struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
	Key3 string `json:"key3"`
}

func TestHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "helper test suite")
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

	//Context("WriteToCSV() Test", func() {
	//	It("should write to csv", func() {
	//		testStruct := TestStruct{
	//			Key1: "v1",
	//			Key2: "v2",
	//			Key3: "v3",
	//		}
	//		headers := []string{"header1", "header2", "header3"}
	//		data := make([]map[string]interface{}, 1)
	//		data[0] = map(
	//			_: TestStruct{
	//				Key1: "",
	//				Key2: "",
	//				Key3: "",
	//			},
	//	)
	//
	//		rows, err := helper.WriteToCSVFile("test.csv", headers, data)
	//		Expect(err).ToNot(BeNil())
	//		Expect(rows).To(BeEquivalentTo(1))
	//	})
	//})
})
