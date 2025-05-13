package runner_test

import (
	"github.com/apkatsikas/subcordant/runner"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("runner", func() {
	var subcordantRunner *runner.SubcordantRunner
	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("will run", func() {
		subcordantRunner.HandlePlay("foobar")
		Expect(1).To(Equal(1))
	})
})
