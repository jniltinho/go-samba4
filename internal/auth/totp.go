package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/pquerna/otp/totp"
)

// GenerateSecret generates a new base32 TOTP secret for a user
func GenerateSecret(accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Samba4 Admin",
		AccountName: accountName,
	})

	if err != nil {
		return "", "", fmt.Errorf("could not generate TOTP secret: %w", err)
	}

	return key.Secret(), key.URL(), nil
}

// ValidatePasscode checks if the given code matches the secret
func ValidatePasscode(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}

// GenerateRecoveryCodes generates 8 backup codes of 10 digits/characters
func GenerateRecoveryCodes() []string {
	codes := make([]string, 8)
	for i := 0; i < 8; i++ {
		b := make([]byte, 5) // 40 bits = 8 base32 chars
		rand.Read(b)
		codes[i] = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	}
	return codes
}
