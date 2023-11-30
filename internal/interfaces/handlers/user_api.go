package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	errors2 "github.com/bonus2k/go-musthave-diploma-tpl/internal/errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/services"
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
		if errors.Is(err, errors2.ErrIllegalUserArgument) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, errors2.ErrUserIsExist) {
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
		if errors.Is(err, errors2.ErrWrongAuth) || errors.Is(err, errors2.ErrUserNotFound) {
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
	userID := r.Header.Get("user")
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		internal.Log.Error("can't get body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = hu.us.AddOrder(r.Context(), userID, string(body)); err != nil {
		internal.Log.Error("add order", zap.Error(err))
		if errors.Is(err, errors2.ErrOrderIsExistThisUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, errors2.ErrOrderIsExistAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, errors2.ErrIllegalOrder) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (hu *HandlerUser) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user")
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
	userID := r.Header.Get("user")
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

func (hu *HandlerUser) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user")
	balance, err := hu.us.GetBalance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(balance); err != nil {
		internal.Log.Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hu *HandlerUser) AddWithdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user")
	internal.Log.Debug("decoding message")
	var dto internal.WithdrawDto
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dto); err != nil {
		internal.Logf.Errorf("cannot decode request JSON body %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := hu.us.AddWithdraw(r.Context(), dto, userID); err != nil {
		internal.Log.Error("add withdraw", zap.Error(err))
		if errors.Is(err, errors2.ErrNotEnoughAmount) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, errors2.ErrIllegalOrder) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
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
