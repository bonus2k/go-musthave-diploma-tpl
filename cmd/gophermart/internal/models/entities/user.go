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
	Bill     int       `db:"bill"`
}

type Order struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Number   int64     `db:"number"`
	Accrual  int64     `db:"accrual"`
	Status   value     `db:"status"`
	UserID   uuid.UUID `db:"user_id"`
}

type Withdraw struct {
	ID       uuid.UUID `db:"id"`
	CreateAt time.Time `db:"create_at"`
	Order    int64     `db:"order_num"`
	Sum      int64     `db:"sum"`
	UserID   uuid.UUID `db:"user_id"`
}

type value string

const (
	NEW        value = "NEW"
	PROCESSING value = "PROCESSING"
	INVALID    value = "INVALID"
	PROCESSED  value = "PROCESSED"
)
