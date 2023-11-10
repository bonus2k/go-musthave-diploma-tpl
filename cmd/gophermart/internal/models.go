package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
	Bill     float32   `db:"bill"`
}

type Order struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Number   int64     `db:"number"`
	Accrual  float32   `db:"accrual"`
	Status   Value     `db:"status"`
	UserID   uuid.UUID `db:"user_id"`
}

type Withdraw struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Order    int64     `db:"order_num"`
	Sum      float32   `db:"sum"`
	UserID   uuid.UUID `db:"user_id"`
}

type Value string

const (
	NEW        Value = "NEW"
	PROCESSING Value = "PROCESSING"
	INVALID    Value = "INVALID"
	PROCESSED  Value = "PROCESSED"
	REGISTERED Value = "REGISTERED"
)

type UserDto struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type OrderDto struct {
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual float32   `json:"accrual"`
	Upload  time.Time `json:"uploaded_at"`
}

func (t *OrderDto) MarshalJSON() ([]byte, error) {
	type Alias OrderDto
	return json.Marshal(&struct {
		*Alias
		Upload string `json:"uploaded_at"`
	}{
		Alias:  (*Alias)(t),
		Upload: t.Upload.Format(time.RFC3339),
	})
}

type AccrualDto struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

type WithdrawDto struct {
	Order    string    `json:"order"`
	Sum      float32   `json:"sum"`
	CreateAt time.Time `json:"processed_at"`
}

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

func (t *WithdrawDto) MarshalJSON() ([]byte, error) {
	type Alias WithdrawDto
	return json.Marshal(&struct {
		*Alias
		CreateAt string `json:"processed_at"`
	}{
		Alias:    (*Alias)(t),
		CreateAt: t.CreateAt.Format(time.RFC3339),
	})
}
