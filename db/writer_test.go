package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	"time"
	"strconv"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"errors"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
)

var _ = Describe("Writer", func() {

	var (
		defaultMfaCredential entity.MfaCredential
	)

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		newUserID := getRandomTimestamp()
		defaultMfaCredential = insertGoogleMfaCredential(newUserID)
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
		mfaCredentialId3 := insertGoogleMfaCredential("userid_3")

		err = db2.Write(db, updatedMfaCredential)
		Expect(err).NotTo(HaveOccurred())

		mfaCredentials, err := db2.ReadAll(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(mfaCredentials).To(ConsistOf(updatedMfaCredential, mfaCredentialId3))
	})

	Describe("when db error occurs", func() {
		var mockDb *dbfakes.FakeUpdater
		BeforeEach(func() {
			mockDb = &dbfakes.FakeUpdater{}
			mockDb.QueryReturns(nil, errors.New("some db error"))
		})
		It("should return meaningful error", func() {
			err := db2.Write(mockDb, entity.MfaCredential{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Unable to update mfa db record: some db error"))
		})
	})
})

func getRandomTimestamp() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
