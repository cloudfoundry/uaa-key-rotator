package testutils

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	Scheme   string
	Hostname string
	Port     string
	Username string
	Password string
	DBName   string
)

func MigrateUaaDatabase() {
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

func TestDBConnection() *sqlx.DB {
	var found bool
	Scheme, found = os.LookupEnv("DB_SCHEME")
	Expect(found).To(BeTrue(), "DB_SCHEME env variable is required")
	Hostname, found = os.LookupEnv("DB_HOSTNAME")
	Expect(found).To(BeTrue(), "DB_HOSTNAME env variable is required")
	Port, found = os.LookupEnv("DB_PORT")
	Expect(found).To(BeTrue(), "DB_PORT env variable is required")
	Username, found = os.LookupEnv("DB_USERNAME")
	Expect(found).To(BeTrue(), "DB_USERNAME env variable is required")
	DBName, found = os.LookupEnv("DB_NAME")
	Expect(found).To(BeTrue(), "DB_NAME env variable is required")
	Password = os.Getenv("DB_PASSWORD")
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Scheme, Username, Password, Hostname, Port, DBName)

	db, err := sqlx.Open(Scheme, connStr)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).Should(BeNil())
	return db
}

func RebindForSQLDialect(query, dialect string) string {
	if dialect == "mysql" {
		return query
	}
	if dialect != "postgres" {
		panic(fmt.Sprintf("Unrecognized DB dialect '%s'", dialect))
	}

	strParts := strings.Split(query, "?")
	for i := 1; i < len(strParts); i++ {
		strParts[i-1] = fmt.Sprintf("%s$%d", strParts[i-1], i)
	}
	return strings.Join(strParts, "")
}
