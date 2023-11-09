package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/services"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type HandlerUser struct {
	us     *services.UserService
	secret []byte
}

func NewHandlerUser(service *services.UserService, secretKey []byte) *HandlerUser {

	return &HandlerUser{us: service, secret: secretKey}
}

func (hu *HandlerUser) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	internal.Log.Debug("decoding message")
	var user internal.UserDto
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		internal.Logf.Errorf("cannot decode request JSON body %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser, err := hu.us.CreateNewUser(r.Context(), &user)
	if err != nil {
		internal.Log.Error("user hasn't been created", zap.Error(err))
		if errors.Is(err, services.ErrIllegalUserArgument) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, repositories.ErrUserIsExist) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	signed := writeSigned(newUser.ID.String(), hu.secret)
	http.SetCookie(w, &signed)
	w.WriteHeader(http.StatusOK)
}

func (hu *HandlerUser) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	internal.Log.Debug("decoding message")
	var user internal.UserDto
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		internal.Logf.Errorf("cannot decode request JSON body %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := hu.us.LoginUser(r.Context(), user)
	if err != nil {
		internal.Log.Error("authorization fault", zap.Error(err))
		if errors.Is(err, services.ErrWrongAuth) || errors.Is(err, repositories.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	signed := writeSigned(userID.String(), hu.secret)
	http.SetCookie(w, &signed)
	w.WriteHeader(http.StatusOK)
}

func (hu *HandlerUser) AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := readSigned(r, hu.secret)
	if err != nil {
		internal.Log.Error("cookie is wrong", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internal.Log.Error("can't get body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = hu.us.AddOrder(r.Context(), userID, string(body)); err != nil {
		internal.Log.Error("add order", zap.Error(err))
		if errors.Is(err, repositories.ErrOrderIsExistThisUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, repositories.ErrOrderIsExistAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, services.ErrIllegalOrder) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (hu *HandlerUser) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := readSigned(r, hu.secret)
	if err != nil {
		internal.Log.Error("cookie is wrong", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := hu.us.GetOrders(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(orders); err != nil {
		internal.Log.Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hu *HandlerUser) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := readSigned(r, hu.secret)
	if err != nil {
		internal.Log.Error("cookie is wrong", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := hu.us.GetWithdrawals(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(withdrawals); err != nil {
		internal.Log.Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeSigned(value string, secret []byte) http.Cookie {
	cookie := http.Cookie{
		Name:     "gophermart",
		Value:    value,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(cookie.Name))
	hash.Write([]byte(cookie.Value))
	signature := hash.Sum(nil)
	value = string(signature) + cookie.Value
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(value))
	return cookie
}

func readSigned(r *http.Request, secretKey []byte) (string, error) {
	cookie, err := r.Cookie("gophermart")
	if err != nil {
		return "", err
	}
	signedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write(value)
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

var ErrInvalidValue = errors.New("invalid cookie value")
