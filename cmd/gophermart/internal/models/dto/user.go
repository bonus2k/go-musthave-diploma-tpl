package dto

import (
	"encoding/json"
	"time"
)

type User struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type Order struct {
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual int64     `json:"accrual,omitempty"`
	Upload  time.Time `json:"uploaded_at"`
}

func (t *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		*Alias
		Upload string `json:"uploaded_at"`
	}{
		Alias:  (*Alias)(t),
		Upload: t.Upload.Format(time.RFC3339),
	})
}
