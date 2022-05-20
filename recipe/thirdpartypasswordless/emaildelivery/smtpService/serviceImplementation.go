package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	evsmtpService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	plesssmtpService "github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeServiceImplementation(config emaildelivery.SMTPServiceConfig) emaildelivery.SMTPServiceInterface {
	sendRawEmail := func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
		return emaildelivery.SendSMTPEmail(config, input)
	}

	evServiceImpl := evsmtpService.MakeServiceImplementation(config)
	plessServiceImpl := plesssmtpService.MakeServiceImplementation(config)

	getContent := func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
		if input.EmailVerification != nil {
			return (*evServiceImpl.GetContent)(input, userContext)

		} else if input.PasswordlessLogin != nil {
			return (*plessServiceImpl.GetContent)(input, userContext)

		} else {
			return emaildelivery.SMTPGetContentResult{}, errors.New("should never come here")
		}
	}

	return emaildelivery.SMTPServiceInterface{
		SendRawEmail: &sendRawEmail,
		GetContent:   &getContent,
	}
}
