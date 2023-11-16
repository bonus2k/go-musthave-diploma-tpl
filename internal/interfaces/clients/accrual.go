package clients

import (
	"errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type ClientAccrual struct {
	client    *resty.Client
	serverURL string
}

func NewClientAccrual(client *resty.Client, serverURL string) *ClientAccrual {
	return &ClientAccrual{client: client, serverURL: serverURL}
}

func (ca *ClientAccrual) CheckAccrual(number string) (*internal.AccrualDto, error) {
	accrual := internal.AccrualDto{}
	response, err := ca.client.R().
		SetResult(&accrual).
		SetRawPathParam("number", number).
		Get(ca.serverURL + "/api/orders/{number}")
	if response.StatusCode() == http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}
	if strings.Contains(response.Status(), http.StatusText(http.StatusNoContent)) {
		return nil, ErrNoContent
	}
	if err != nil {
		return nil, err
	}
	return &accrual, nil
}

var ErrTooManyRequests = errors.New("too many requests")
var ErrNoContent = errors.New("no content")
