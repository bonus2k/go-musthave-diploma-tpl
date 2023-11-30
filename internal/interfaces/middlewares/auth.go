package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	errors2 "github.com/bonus2k/go-musthave-diploma-tpl/internal/errors"
	"go.uber.org/zap"
	"net/http"
)

func Authentication(secretKey []byte) func(http.Handler) http.Handler {
	secret := secretKey
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := readSigned(r, secret)
			if err != nil {
				if err != nil {
					internal.Log.Error("cookie is wrong", zap.Error(err))
					http.Error(w, "cookie is wrong", http.StatusUnauthorized)
					return
				}
			}
			r.Header.Add("user", userID)
			h.ServeHTTP(w, r)
		})
	}
}

func readSigned(r *http.Request, secretKey []byte) (string, error) {
	cookie, err := r.Cookie("gophermart")
	if err != nil {
		return "", err
	}
	signedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", errors2.ErrInvalidValue
	}
	if len(signedValue) < sha256.Size {
		return "", errors2.ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write(value)
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors2.ErrInvalidValue
	}

	return string(value), nil
}
