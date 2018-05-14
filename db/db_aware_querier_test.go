package db_test

import (
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("DbAwareQuerier", func() {
	var dbAwareQuerier db2.DbAwareQuerier

	BeforeEach(func() {
		dbAwareQuerier = db2.DbAwareQuerier{
			DBScheme: "mysql",
		}
	})

	Context("unknown db scheme", func() {
		BeforeEach(func() {
			dbAwareQuerier.DBScheme = "unknown"
		})

		It("should return an error", func() {
			_, err := dbAwareQuerier.Queryx("")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Unable to query: Unrecognized DB dialect 'unknown'"))
		})
	})
})
