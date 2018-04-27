package main_test

import (
	"testing"

	"database/sql"
	"fmt"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func TestUaaKeyRotator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UaaKeyRotator Suite")
}

var uaaRotatorBuildPath string
var db *sqlx.DB
var activeKey config.EncryptionKey
var oldKey config.EncryptionKey

var _ = BeforeSuite(func() {
	var err error
	uaaRotatorBuildPath, err = gexec.Build("github.com/cloudfoundry/uaa-key-rotator")
	Expect(err).NotTo(HaveOccurred())

	By("validating and testing the db connection", func() {
		db = testutils.TestDBConnection()
	})

	By("migrating UAA database", testutils.MigrateUaaDatabase)

	activeKey = config.EncryptionKey{Label: "active-key", Passphrase: "my-passphrase"}

	By("clearing database of any records", func() {
		_, err := db.Exec(`truncate user_google_mfa_credentials;`)
		Expect(err).NotTo(HaveOccurred())
	})

	By("adding test fixtures", testFixtures)
})

func testFixtures() {
	oldKey = config.EncryptionKey{
		Label:      "old-key-label",
		Passphrase: "secret",
	}

	secretKeyCipherValue := encryptPlainText("secret-key", oldKey.Passphrase)
	scratchCodesCipherValue := encryptPlainText("scratchCodes", oldKey.Passphrase)
	encryptedValidationCodesCipherValue := encryptPlainText("encryptedValidationCodes", oldKey.Passphrase)

	mfaCredential := entity.MfaCredential{
		UserId:                  "user-id-1",
		SecretKey:               secretKeyCipherValue,
		ScratchCodes:            scratchCodesCipherValue,
		MfaProviderId:           "mfa_provider_id",
		ZoneId:                  "zone_id",
		EncryptionKeyLabel:      oldKey.Label,
		EncryptedValidationCode: encryptedValidationCodesCipherValue,
		ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
	}

	insertSQL := testutils.RebindForSQLDialect(`insert into user_google_mfa_credentials(
		user_id, 
		secret_key, 
		validation_code, 
		scratch_codes, 
		mfa_provider_id, 
		zone_id, 
		encryption_key_label, 
		encrypted_validation_code) values(
		?, ?, ?, ?, ?, ?, ?, ?)`,
		testutils.Scheme)

	insertResult, err := db.Exec(insertSQL,
		mfaCredential.UserId,
		mfaCredential.SecretKey,
		mfaCredential.ValidationCode,
		mfaCredential.ScratchCodes,
		string(mfaCredential.MfaProviderId),
		string(mfaCredential.ZoneId),
		mfaCredential.EncryptionKeyLabel,
		mfaCredential.EncryptedValidationCode)
	Expect(err).NotTo(HaveOccurred())
	numOfRowsInserted, err := insertResult.RowsAffected()
	Expect(err).NotTo(HaveOccurred())
	Expect(numOfRowsInserted).To(Equal(int64(1)))
}

func encryptPlainText(plainText string, passphrase string) string {
	cipherValueFile, err := ioutil.TempFile(os.TempDir(), "encrypt")
	Expect(err).NotTo(HaveOccurred())
	defer os.Remove(cipherValueFile.Name())

	gradleCmd := exec.Command("./gradlew", "--quiet", "--system-prop", fmt.Sprintf(`encryptArgs=%s,%s,%s`, cipherValueFile.Name(), passphrase, plainText), "encrypt")
	gradleCmd.Dir = os.Getenv("UAA_LOCATION")
	gradleSession, err := gexec.Start(gradleCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(gradleSession, 5*time.Minute).Should(gexec.Exit())

	cipherValueContents, err := ioutil.ReadAll(cipherValueFile)
	Expect(err).NotTo(HaveOccurred())

	return string(cipherValueContents)
}

func decryptCipherValue(cipherValue string, passphrase string) string {
	plainTextFile, err := ioutil.TempFile(os.TempDir(), "decrypt")
	Expect(err).NotTo(HaveOccurred())
	defer os.Remove(plainTextFile.Name())

	gradleCmd := exec.Command("./gradlew", "--quiet", "--system-prop", fmt.Sprintf(`decryptArgs=%s,%s,%s`, plainTextFile.Name(), passphrase, cipherValue), "decrypt")
	gradleCmd.Dir = os.Getenv("UAA_LOCATION")

	gradleSession, err := gexec.Start(gradleCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(gradleSession, 5*time.Minute).Should(gexec.Exit())

	plainTextContents, err := ioutil.ReadAll(plainTextFile)
	Expect(err).NotTo(HaveOccurred())

	return string(plainTextContents)
}

var _ = AfterSuite(func() {
	Expect(os.Remove(uaaRotatorBuildPath)).To(Succeed())
	By("clearing database of any records", func() {
		if db != nil {
			_, err := db.Exec(`truncate user_google_mfa_credentials;`)
			Expect(err).NotTo(HaveOccurred())
		}
	})
	db.Close()
})
