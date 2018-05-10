package db_test

import (
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MakeSqlQueriesCompatWithDbs", func() {
	var query string

	BeforeEach(func() {
		query = "select 1 from table where col1 = ?"
	})

	Describe("mysql", func() {
		It("should not modify mysql queries", func() {
			reboundQuery, err := db2.RebindForSQLDialect(query, "mysql")
			Expect(err).ToNot(HaveOccurred())
			Expect(reboundQuery).To(Equal(query))
		})
	})

	Describe("postgresql", func() {
		It("should not modify postgresql queries", func() {
			reboundQuery, err := db2.RebindForSQLDialect(query, "postgres")
			Expect(err).ToNot(HaveOccurred())
			Expect(reboundQuery).To(Equal("select 1 from table where col1 = $1"))
		})
	})

	Describe("sql server", func() {
		It("should modify sql server queries", func() {
			reboundQuery, err := db2.RebindForSQLDialect(query, "sqlserver")
			Expect(err).ToNot(HaveOccurred())
			Expect(reboundQuery).To(Equal("select 1 from table where col1 = @p1"))
		})
	})

	It("should fail if invalid db scheme is provided", func() {
		_, err := db2.RebindForSQLDialect(query, "some-unknown-database")
		Expect(err).To(HaveOccurred())
	})
})
