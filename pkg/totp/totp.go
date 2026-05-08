package totp

import (
	"time"

	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "CreateKey failed, accountName: %s", accountName)
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

		return false, errors.Wrapf(err, "ValidateCode failed, code: %s", code)
	}
	return ok, nil
}
