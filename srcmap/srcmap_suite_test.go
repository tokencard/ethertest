package srcmap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tokencard/ethertest/srcmap"
)

func TestSrcmap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Srcmap Suite")
}

var _ = Describe("Uncompress", func() {
	// -1:-1:-1:;11386:87:2;;;11452:10;;11434:15;:28;11386:87;11081:398;:::o;20850:::-;20960:22;:20;:22::i
	It("Should not change uncompressed entry", func() {
		m, err := srcmap.Uncompress("1:2:3:abc")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc"))
	})

	It("Should repeat previous entry if the next is empty", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;1:2:3:abc"))
	})

	It("Should propagate first value if unspecified", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;:4:5:def")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;1:4:5:def"))
	})

	It("Should propagate second value if unspecified", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;4::5:def")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;4:2:5:def"))
	})

	It("Should propagate third value if unspecified", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;4:5::def")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;4:5:3:def"))
	})

	It("Should propagate fourth value if unspecified", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;4:5:6:")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;4:5:6:abc"))
	})

	It("Should propagate fourth value if shorter than 4", func() {
		m, err := srcmap.Uncompress("1:2:3:abc;4:5:6")
		Expect(err).ToNot(HaveOccurred())
		Expect(m.String()).To(Equal("1:2:3:abc;4:5:6:abc"))
	})

})
