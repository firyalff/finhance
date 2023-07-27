package auth

import (
	"context"
	"finhancesvc/shared"
	"log"

	"github.com/jackc/pgx/v5"
)

func getUserByCredentials(ctx context.Context, credentials loginPayload) (user userDB, err error) {
	userRecord, err := getUserByEmail(ctx, credentials.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return userDB{}, shared.ErrNotFound
		}
		return
	}

	err = validateUserPassword(credentials.Password, userRecord.Password)
	if err != nil {
		return userDB{}, shared.ErrUnauthorized
	}

	return userRecord, nil
}

func generateAthenticationToken(userID string, secretToken string, authExpirationDays int) (token string, err error) {
	return generateJWT(userID, secretToken, authExpirationDays)
}

func validateUniqueEmailRegistration(ctx context.Context, email string) (err error) {
	totalUser, err := countUserByEmail(ctx, email)
	if err != nil {
		return shared.ErrInternal
	}

	if totalUser > 0 {
		return shared.ErrExist
	}

	return
}

func registerNewUser(ctx context.Context, registrationRecord registerPayload, baseURL string) (err error) {
	const REGISTRATION_MAIL_SUBJECT = "Finhance New User"

	tx, err := AuthModuleInstance.dbPool.Begin(ctx)
	if err != nil {
		log.Print(err)
		return shared.ErrInternal
	}

	userID, err := createUser(ctx, tx, registrationRecord)
	if err != nil {
		return shared.ErrInternal
	}

	activationToken, err := createUserActivation(ctx, tx, userID)
	if err != nil {
		return shared.ErrInternal
	}

	accountActivationLink := baseURL + "/account-activation?registration_token=" + activationToken
	registrationMailContent := "Please activate your account by visiting <a href=" + accountActivationLink + ">this link</a>"

	mailSession := shared.SMTPAuth(AuthModuleInstance.serverConfig.SMTPUsername, AuthModuleInstance.serverConfig.SMTPPassword, AuthModuleInstance.serverConfig.SMTPHostURL)

	mailContent := shared.SMTPMailConfig{
		Sender:    AuthModuleInstance.serverConfig.EmailDefaultSender,
		Recipient: registrationRecord.Email,
		Subject:   REGISTRATION_MAIL_SUBJECT,
		Body:      registrationMailContent,
	}

	err = shared.SMTPSendMail(AuthModuleInstance.serverConfig.SMTPHostURL, AuthModuleInstance.serverConfig.SMTPHostPORT, mailSession, mailContent)
	if err != nil {

		log.Print(err)
		return shared.ErrInternal
	}

	tx.Commit(ctx)

	return
}

func activateAccount(ctx context.Context, activationToken string) (err error) {
	tx, err := AuthModuleInstance.dbPool.Begin(ctx)
	if err != nil {
		log.Print(err)
		return shared.ErrInternal
	}

	totalFound, err := countUserActivation(ctx, tx, activationToken)
	if err != nil {
		err = shared.ErrInternal
	}

	if totalFound < 1 {
		return shared.ErrNotFound
	}

	err = activateAccountWithActivationToken(ctx, tx, activationToken)
	if err != nil {
		err = shared.ErrInternal
	}

	err = deleteActivationToken(ctx, tx, activationToken)
	if err != nil {
		err = shared.ErrInternal
	}

	tx.Commit(ctx)

	return
}
