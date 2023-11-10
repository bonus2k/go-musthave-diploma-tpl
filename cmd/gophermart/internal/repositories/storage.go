package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
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
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return store.db.PingContext(ctx)
}

func (store *Store) AddUser(ctx context.Context, user *internal.User) error {
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

func (store *Store) FindUserByLogin(ctx context.Context, login string) (*internal.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var user internal.User
	err := store.db.GetContext(timeout, &user, "SELECT * FROM users WHERE login=$1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get user from db, %w", err)
	}
	return &user, nil
}

func (store *Store) AddOrder(ctx context.Context, order *internal.Order) (*internal.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var existOrder internal.Order
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

func (store *Store) GetOrders(ctx context.Context, userID uuid.UUID) (*[]internal.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var orders []internal.Order
	err := store.db.SelectContext(timeout, &orders, `SELECT * FROM orders WHERE user_id=$1 
                     ORDER BY create_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *Store) GetOrder(ctx context.Context, login string, number int64) (*internal.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var user internal.User
	err := store.db.SelectContext(timeout, &user, `SELECT * FROM users WHERE login=$1`, login)
	if err != nil {
		return nil, fmt.Errorf("can't get user from db %w", err)
	}
	var orders internal.Order
	err = store.db.SelectContext(timeout, &orders, `SELECT * FROM orders WHERE user_id=$1 AND number=$2`,
		user.ID, number)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *Store) GetOrdersNotProcessed(ctx context.Context) (*[]internal.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	var orders []internal.Order
	err := store.db.SelectContext(timeout, &orders, `SELECT * FROM orders 
         WHERE status != 'INVALID' AND status != 'PROCESSED'`)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *Store) UpdateOrder(ctx context.Context, order *internal.Order) error {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	tx := store.db.MustBeginTx(timeout, nil)
	defer tx.Commit()
	if order.Status == "PROCESSED" {
		var user internal.User
		err := tx.GetContext(timeout, &user, `SELECT u.id, u.bill FROM users AS u INNER JOIN orders AS s ON u.id = s.user_id 
                    WHERE s.number=$1`, order.Number)
		if err != nil {
			tx.Rollback()
			return err
		}
		sumBill := user.Bill + order.Accrual
		_, err = tx.ExecContext(timeout, `UPDATE users SET bill = $1 WHERE id = $2`, sumBill, user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err := tx.ExecContext(timeout, `UPDATE orders SET status = $1, accrual=$2 WHERE number = $3`, order.Status, order.Accrual, order.Number)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (store *Store) SaveWithdrawal(ctx context.Context, withdrawal *internal.Withdraw) error {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5000)
	defer cancelFunc()
	tx := store.db.MustBeginTx(timeout, nil)
	defer tx.Commit()
	var user internal.User
	err := tx.GetContext(timeout, &user, "SELECT * FROM users WHERE id=$1", withdrawal.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("can't get user from db %w", err)
	}
	if user.Bill < withdrawal.Sum {
		return ErrNotEnoughAmount
	}
	_, err = tx.NamedExecContext(timeout, `INSERT INTO withdrawals (id, create_at, order_num, sum, user_id) VALUES (:id, :create_at, :order_num, :sum, :user_id)`, withdrawal)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("can't save withdrawal to db %w", err)
	}
	sumBill := user.Bill - withdrawal.Sum
	_, err = tx.ExecContext(timeout, `UPDATE users SET bill = $1 WHERE id = $2`, sumBill, withdrawal.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("can't update user bill at db %w", err)
	}
	return nil
}

func (store *Store) GetWithdrawals(ctx context.Context, userID uuid.UUID) (*[]internal.Withdraw, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	var withdrawals []internal.Withdraw
	err := store.db.SelectContext(timeout, &withdrawals, `SELECT * FROM withdrawals WHERE user_id=$1 
                     ORDER BY create_at DESC `, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &withdrawals, nil
}

func (store *Store) GetUser(ctx context.Context, id uuid.UUID) (*internal.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	var user internal.User
	err := store.db.GetContext(timeout, &user, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user from db %w", err)
	}
	return &user, nil
}

var ErrUserIsExist = errors.New("user is exist")
var ErrUserNotFound = errors.New("user not found")
var ErrOrderIsExistThisUser = errors.New("this order is exist the user")
var ErrOrderIsExistAnotherUser = errors.New("this order is exist another user")
var ErrNotEnoughAmount = errors.New("not enough amount")
