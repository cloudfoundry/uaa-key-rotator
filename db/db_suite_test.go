package db_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/lib/pq"
	"os"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"path/filepath"
	"io"
	"time"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func() {
	By("validating and testing the db connection", testDBConnection)
	By("migrating UAA database", migrateUaaDatabase)
	// create config file to point to db (env variables)
	// insert data
})

var db *sqlx.DB

func migrateUaaDatabase() {
	uaaLocation, found := os.LookupEnv("UAA_LOCATION")
	Expect(found).To(BeTrue(), "UAA_LOCATION env variable is required")
	gradlePath := filepath.Join(uaaLocation, "gradlew")
	gradleMigrateCommand := exec.Command(gradlePath, "flywayMigrate", "-Dspring.profiles.active=postgresql")
	gradleMigrateCommand.Dir = uaaLocation
	gradleMigrateCommand.Env = append(gradleMigrateCommand.Env, fmt.Sprintf("JAVA_HOME=%s", os.Getenv("JAVA_HOME")))
	session, err := gexec.Start(gradleMigrateCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	go func() {
		io.Copy(GinkgoWriter, session.Out)
	}()
	go func() {
		io.Copy(GinkgoWriter, session.Err)
	}()
	Eventually(session, 5*time.Minute).Should(gexec.Exit())
}

func testDBConnection() {
	scheme, found := os.LookupEnv("DB_SCHEME")
	Expect(found).To(BeTrue(), "DB_SCHEME env variable is required")
	hostname, found := os.LookupEnv("DB_HOSTNAME")
	Expect(found).To(BeTrue(), "DB_HOSTNAME env variable is required")
	username, found := os.LookupEnv("DB_USERNAME")
	Expect(found).To(BeTrue(), "DB_USERNAME env variable is required")
	dbname, found := os.LookupEnv("DB_NAME")
	Expect(found).To(BeTrue(), "DB_NAME env variable is required")
	password := os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", scheme, username, password, hostname, dbname)

	var err error
	db, err = sqlx.Open("postgres", connStr)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).Should(BeNil())
}
