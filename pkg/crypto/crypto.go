package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
)

type EncryptedString struct {
	KeyID     string `json:"key_id"`
	Algorithm string `json:"alg"`
	Data      []byte `json:"data"`
	Nonce     []byte `json:"nonce,omitempty"`
}

func (es *EncryptedString) IsValid() bool {
	return es.KeyID != "" && len(es.Data) > 0 && len(es.Nonce) > 0 && es.Algorithm == "aes-gcm-hkdf"
}

// SecureToken creates a new random token
func SecureToken(options ...int) string {
	length := 16
	if len(options) > 0 {
		length = options[0]
	}
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GenerateOtp generates a random n digit otp
func GenerateOtp(digits int) (string, error) {
	upper := math.Pow10(digits)
	val, err := rand.Int(rand.Reader, big.NewInt(int64(upper)))
	if err != nil {
		return "", errors.WithMessage(err, "Error generating otp")
	}
	// adds a variable zero-padding to the left to ensure otp is uniformly random
	expr := "%0" + strconv.Itoa(digits) + "v"
	otp := fmt.Sprintf(expr, val.String())
	return otp, nil
}

func GenerateTokenHash(emailOrPhone, otp string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(emailOrPhone+otp)))
}
