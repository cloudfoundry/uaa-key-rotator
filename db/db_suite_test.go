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
	"strings"
	"database/sql"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func() {
	By("validating and testing the db connection", testDBConnection)
	By("migrating UAA database", migrateUaaDatabase)
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

func insertGoogleMfaCredential(userId string) entity.MfaCredential {
	mfaCredential := entity.MfaCredential{
		UserId: userId,
		SecretKey: "secret-key",
		ScratchCodes: "scratch_codes",
		MfaProviderId: "mfa_provider_id",
		ZoneId: "zone_id",
		EncryptionKeyLabel: "activeKeyLabel",
		EncryptedValidationCode: "encrypted_validation_code",
		ValidationCode: sql.NullInt64{Int64: 1234, Valid: true},
	}

	insertSQL := RebindForSQLDialect(`insert into user_google_mfa_credentials(
		user_id, 
		secret_key, 
		validation_code, 
		scratch_codes, 
		mfa_provider_id, 
		zone_id, 
		encryption_key_label, 
		encrypted_validation_code) values(
		?, ?, ?, ?, ?, ?, ?, ?
		)`,
		"postgres")

	insertResult, err := db.Exec(insertSQL, mfaCredential.UserId,
		mfaCredential.SecretKey,
		mfaCredential.ValidationCode,
		mfaCredential.ScratchCodes,
		mfaCredential.MfaProviderId,
		mfaCredential.ZoneId,
		mfaCredential.EncryptionKeyLabel,
		mfaCredential.EncryptedValidationCode)

	Expect(err).NotTo(HaveOccurred())
	numOfRowsInserted, err := insertResult.RowsAffected()
	Expect(err).NotTo(HaveOccurred())
	Expect(numOfRowsInserted).To(Equal(int64(1)))
	return mfaCredential
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
