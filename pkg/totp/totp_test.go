package totp

import (
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestCreateKey(t *testing.T) {
	type testCase struct {
		accountName string
		wantErr     bool
	}

	testCases := map[string]testCase{
		"valid account":      {"Cai", false},
		"empty account name": {"", true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			key, err := CreateKey(tc.accountName)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.accountName, key.AccountName())
			assert.Equal(t, constants.TotpIssuer, key.Issuer())
			assert.Equal(t, constants.TotpPeriod, int(key.Period()))
			assert.Equal(t, constants.TotpDigitLength, key.Digits().Length())
		})
	}
}

func TestValidateCode(t *testing.T) {
	type testCase struct {
		secret    string
		code      string
		timestamp int64
		wantValid bool
		wantErr   bool
	}

	testSecret := "BQRIBVKD6QIBIIZMRZVMF2FL7RTOQSJO"
	testCode := "757019"
	var testTimestamp int64 = 1778496825343

	testCases := map[string]testCase{
		"valid": {
			secret:    testSecret,
			code:      testCode,
			timestamp: testTimestamp,
			wantValid: true,
			wantErr:   false,
		},
		"expired code": {
			secret:    testSecret,
			code:      testCode,
			timestamp: testTimestamp + time.Minute.Milliseconds(),
			wantValid: false,
			wantErr:   false,
		},
		"invalid code": {
			secret:    testSecret,
			code:      "1145141919810",
			timestamp: testTimestamp,
			wantValid: false,
			wantErr:   false,
		},
		"invalid secret": {
			secret:    "原神",
			code:      testCode,
			timestamp: testTimestamp,
			wantValid: false,
			wantErr:   true,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			testTime := time.UnixMilli(tc.timestamp)
			mockey.Mock(time.Now).Return(testTime).Build()

			valid, err := ValidateCode(tc.secret, tc.code)
			if tc.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tc.wantValid, valid)
		})
	}
}
