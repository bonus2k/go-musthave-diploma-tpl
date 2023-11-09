package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/models/entities"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	db *sqlx.DB
}

var s *Store

func NewStore(databaseURI string) (*Store, error) {
	if s != nil {
		return s, nil
	}
	open, err := sqlx.Open("pgx", databaseURI)
	if err != nil {
		return nil, fmt.Errorf("can't create connection to DB %w", err)
	}
	s = &Store{db: open}
	return s, nil
}

func (store *Store) CheckConnection() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	return store.db.PingContext(ctx)
}

func (store *Store) AddUser(ctx context.Context, user *entities.User) error {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var count int
	err := store.db.GetContext(timeout, &count, "SELECT count(*) FROM users WHERE login=$1", user.Login)
	if err != nil {
		return fmt.Errorf("can't check user is exist %w", err)
	}
	if count > 0 {
		return ErrUserIsExist
	}
	_, err = store.db.NamedExecContext(timeout,
		`INSERT INTO users (id, create_at, login, password, bill) VALUES (:id, :create_at, :login, :password, :bill)`,
		user)
	if err != nil {
		return fmt.Errorf("can't save user to db %w", err)
	}
	return nil
}

func (store *Store) FindUserByLogin(ctx context.Context, login string) (*entities.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var user entities.User
	err := store.db.GetContext(timeout, &user, "SELECT * FROM users WHERE login=$1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get user from db, %w", err)
	}
	return &user, nil
}

func (store *Store) AddOrder(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*500)
	defer cancelFunc()
	var existOrder entities.Order
	err := store.db.GetContext(timeout,
		&existOrder,
		`SELECT * FROM orders WHERE number = $1`,
		order.Number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("can't get order from db %w", err)
	}
	if existOrder.ID != uuid.Nil {
		if existOrder.UserID == order.UserID {
			return nil, ErrOrderIsExistThisUser
		} else {
			return nil, ErrOrderIsExistAnotherUser
		}
	}
	_, err = store.db.NamedExecContext(timeout, "INSERT INTO orders (id, create_at, number, accrual, status, user_id) "+
		"VALUES (:id, :create_at, :number, :accrual, :status, :user_id)", order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (store *Store) GetOrders(ctx context.Context, userID uuid.UUID) (*[]entities.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var orders []entities.Order
	err := store.db.SelectContext(timeout, &orders, `SELECT * FROM orders WHERE user_id=$1 
                     ORDER BY create_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *Store) GetOrder(ctx context.Context, login string, number int64) (*entities.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var user entities.User
	err := store.db.SelectContext(timeout, &user, `SELECT * FROM users WHERE login=$1`, login)
	if err != nil {
		return nil, fmt.Errorf("can't get user from db %w", err)
	}
	var orders entities.Order
	err = store.db.SelectContext(timeout, &orders, `SELECT * FROM orders WHERE user_id=$1 AND number=$2`,
		user.ID, number)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *Store) getWithdrawals(ctx context.Context, login string) (*[]entities.Withdraw, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var user entities.User
	err := store.db.SelectContext(timeout, &user, `SELECT * FROM users WHERE login=$1`, login)
	if err != nil {
		return nil, fmt.Errorf("can't get user from db %w", err)
	}
	var withdrawals []entities.Withdraw
	err = store.db.SelectContext(timeout, &withdrawals, `SELECT * FROM withdrawals WHERE user_id=$1 
                     ORDER BY create_at`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &withdrawals, nil
}

var ErrUserIsExist = errors.New("user is exist")
var ErrUserNotFound = errors.New("user not found")
var ErrOrderIsExistThisUser = errors.New("this order is exist the user")
var ErrOrderIsExistAnotherUser = errors.New("this order is exist another user")
