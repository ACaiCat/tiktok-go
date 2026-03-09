package totp

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/pquerna/otp/totp"
)

func CreateSecret(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.TotpIssuer,
		AccountName: accountName,
		Period:      constants.TotpPeriod,
		Digits:      constants.TotpDigitLength,
	})

	if err != nil {
		log.Println("failed to generate totp secret:", err)
		return "", errno.ServiceErr
	}

	return key.String(), err
}

func ValidateCode(secret string, code string) bool {
	return totp.Validate(code, secret)
}
