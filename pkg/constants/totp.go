package constants

import "time"

const (
	TotpIssuer             = "Cai"
	TotpPeriod             = 30
	TotpDigitLength        = 6
	TotpTempSecretCacheKey = "totp_secret:%d" // totp_secret:user_id
	TotpTempSecretTTL      = 5 * time.Minute  // 5分钟
)
