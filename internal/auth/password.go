package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func SignPassword(password string) (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := hmac.New(sha256.New, salt)
	_, err = hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	sum := hash.Sum(salt)
	return hex.EncodeToString(sum), nil
}

func CheckPassword(password string, passwordIsExs string) (bool, error) {
	hashPass, err := hex.DecodeString(passwordIsExs)
	if err != nil {
		return false, err
	}
	hash := hmac.New(sha256.New, hashPass[:8])
	_, err = hash.Write([]byte(password))
	if err != nil {
		return false, err
	}
	sign := hash.Sum(nil)
	return hmac.Equal(sign, hashPass[8:]), nil
}
