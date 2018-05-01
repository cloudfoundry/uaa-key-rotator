package db_test

import (
	"errors"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"time"
)

var _ = Describe("Writer", func() {

	var (
		defaultMfaCredential entity.MfaCredential
		credentialsDB        db2.GoogleMfaCredentialsDBFetcher
		credentialsDBUpdater db2.GoogleMfaCredentialsDBUpdater
	)

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		newUserID := getRandomTimestamp()
		defaultMfaCredential = insertGoogleMfaCredential(newUserID, "activeKeyLabel")

		credentialsDB = db2.GoogleMfaCredentialsDBFetcher{
			DB:             db2.DbAwareQuerier{DB: db, DBScheme: testutils.Scheme},
			ActiveKeyLabel: "some-active-key-label",
		}
		credentialsDBUpdater = db2.GoogleMfaCredentialsDBUpdater{
			DB: db2.DbAwareQuerier{DB: db, DBScheme: testutils.Scheme},
		}
	})

	It("should update a single mfa record", func() {
		var err error

		updatedMfaCredential := entity.MfaCredential{
			UserId:                  defaultMfaCredential.UserId,
			MfaProviderId:           defaultMfaCredential.MfaProviderId,
			ZoneId:                  defaultMfaCredential.ZoneId,
			ValidationCode:          defaultMfaCredential.ValidationCode,
			ScratchCodes:            getRandomTimestamp(),
			SecretKey:               getRandomTimestamp(),
			EncryptionKeyLabel:      getRandomTimestamp(),
			EncryptedValidationCode: getRandomTimestamp(),
		}
		mfaCredentialId3 := insertGoogleMfaCredential("userid_3", "activeKeyLabel")

		err = credentialsDBUpdater.Write(updatedMfaCredential)
		Expect(err).NotTo(HaveOccurred())

		var errChan <-chan error
		mfaCredentials, errChan := credentialsDB.RowsToRotate()
		Consistently(errChan).ShouldNot(Receive())

		var rotatedMfaCredential1 entity.MfaCredential
		var rotatedMfaCredential2 entity.MfaCredential
		Eventually(mfaCredentials).Should(Receive(&rotatedMfaCredential1))
		Eventually(mfaCredentials).Should(Receive(&rotatedMfaCredential2))

		Eventually([]entity.MfaCredential{rotatedMfaCredential1, rotatedMfaCredential2}).Should(ConsistOf(mfaCredentialId3, updatedMfaCredential))
	})

	Describe("when db error occurs", func() {
		var mockDb *dbfakes.FakeQueryer
		BeforeEach(func() {
			mockDb = &dbfakes.FakeQueryer{}
			mockDb.QueryxReturns(nil, errors.New("some db error"))
			credentialsDBUpdater = db2.GoogleMfaCredentialsDBUpdater{
				DB: mockDb,
			}
		})
		It("should return meaningful error", func() {
			err := credentialsDBUpdater.Write(entity.MfaCredential{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Unable to update mfa db record: some db error"))
		})
	})
})

func getRandomTimestamp() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
