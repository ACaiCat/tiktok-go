package totp

import (
	"errors"
	"log"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func CreateKey(accountName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.TotpIssuer,
		AccountName: accountName,
		Period:      constants.TotpPeriod,
		Digits:      constants.TotpDigitLength,
	})

	if err != nil {
		log.Println("failed to generate totp secret:", err)
		return nil, err
	}

	return key, err
}

func ValidateCode(secret string, code string) (bool, error) {
	ok, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period: constants.TotpPeriod,
		Digits: constants.TotpDigitLength,
	})
	if err != nil {
		if errors.Is(err, otp.ErrValidateInputInvalidLength) {
			return false, nil
		}

		log.Println("failed to validate totp code:", err)
		return false, err
	}
	return ok, nil
}
