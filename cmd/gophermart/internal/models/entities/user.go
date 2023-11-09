package entities

import (
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
	Status   value     `db:"status"`
	UserID   uuid.UUID `db:"user_id"`
}

type Withdraw struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Order    int64     `db:"order_num"`
	Sum      float32   `db:"sum"`
	UserID   uuid.UUID `db:"user_id"`
}

type value string

const (
	NEW        value = "NEW"
	PROCESSING value = "PROCESSING"
	INVALID    value = "INVALID"
	PROCESSED  value = "PROCESSED"
)
