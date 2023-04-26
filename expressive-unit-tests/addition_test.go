package addition

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add(x, y)", func() {
	It("sums two numbers", func(ctx SpecContext) {
		Expect(Add(4, 6)).To(Equal(10))
	})
	It("supports negative numbers", func(ctx SpecContext) {
		Expect(Add(-4, 6)).To(Equal(2))
	})
})
