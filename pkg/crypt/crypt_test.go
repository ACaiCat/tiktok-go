package crypt

import (
	"crypto/rand"
	"io"
	"os"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/config"
)

const (
	testKey        = "FnUqo4Qeoiwxv3L8f04is8OhmhTbSvUifSWN6BujF+c="
	testPlaintext  = "114514"
	testCiphertext = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
)

func TestMain(m *testing.M) {
	config.AppConfig.Security.Key = testKey
	code := m.Run()
	os.Exit(code)
}

func TestEncrypt(t *testing.T) {
	type testCase struct {
		key       string
		plaintext string
		mockRand  bool
		wantErr   bool
	}

	testCases := map[string]testCase{
		"encrypt normal text": {
			key:       testKey,
			plaintext: "hello world",
			wantErr:   false,
		},
		"encrypt empty string": {
			key:       testKey,
			plaintext: "",
			wantErr:   false,
		},
		"invalid base64 key": {
			key:       "114514",
			plaintext: "test",
			wantErr:   true,
		},
		"wrong AES key size": {
			key:       "aGVsbG8=", // "hello", not 16/24/32 bytes
			plaintext: "test",
			wantErr:   true,
		},
		"rand reader failure": {
			key:       testKey,
			plaintext: "hello",
			mockRand:  true,
			wantErr:   true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			config.AppConfig.Security.Key = tc.key
			if tc.mockRand {
				Mock(io.ReadFull).To(func(_ io.Reader, _ []byte) (int, error) {
					return 0, assert.AnError
				}).Build()
			}
			_, err := Encrypt(tc.plaintext)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestDecrypt(t *testing.T) {
	type testCase struct {
		key        string
		ciphertext string
		mockNonce  bool
		wantResult string
		wantErr    bool
	}

	testCases := map[string]testCase{
		"decrypt roundtrip normal text": {
			key:        testKey,
			mockNonce:  true,
			wantResult: testPlaintext,
			wantErr:    false,
		},
		"decrypt roundtrip empty string": {
			key:        testKey,
			mockNonce:  true,
			wantResult: "",
			wantErr:    false,
		},
		"invalid base64 ciphertext": {
			key:        testKey,
			ciphertext: "not-valid-base64!!!",
			wantErr:    true,
		},
		"invalid base64 key": {
			key:        "not-valid-base64!!!",
			ciphertext: testCiphertext,
			wantErr:    true,
		},
		"wrong AES key size": {
			key:        "aGVsbG8=", // "hello", not 16/24/32 bytes
			ciphertext: testCiphertext,
			wantErr:    true,
		},
		"ciphertext shorter than nonce": {
			key:        testKey,
			ciphertext: "aGVsbG8=", // "hello", 5 bytes < 12 byte nonce
			wantErr:    true,
		},
		"tampered ciphertext": {
			key:        testKey,
			ciphertext: testCiphertext,
			wantErr:    true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			config.AppConfig.Security.Key = tc.key

			ciphertext := tc.ciphertext
			if tc.mockNonce {
				Mock(rand.Reader.Read).To(func(p []byte) (int, error) {
					clear(p)
					return len(p), nil
				}).Build()
				ciphertext, _ = Encrypt(tc.wantResult)
			}

			result, err := Decrypt(ciphertext)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
