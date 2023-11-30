package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	errors2 "github.com/bonus2k/go-musthave-diploma-tpl/internal/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type StoreImpl struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &StoreImpl{db: db}
}

func (store *StoreImpl) CheckConnection() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	return store.db.PingContext(ctx)
}

func (store *StoreImpl) AddUser(ctx context.Context, user *internal.User) error {
	var count int
	err := store.db.GetContext(ctx, &count, `SELECT count(*) FROM users WHERE login=$1`, user.Login)
	if err != nil {
		return fmt.Errorf("can't check user is exist %w", err)
	}
	if count > 0 {
		return errors2.ErrUserIsExist
	}
	_, err = store.db.NamedExecContext(ctx,
		`INSERT INTO users (id, create_at, login, password, bill) VALUES (:id, :create_at, :login, :password, :bill)`,
		user)
	if err != nil {
		return fmt.Errorf("can't save user to db %w", err)
	}
	return nil
}

func (store *StoreImpl) FindUserByLogin(ctx context.Context, login string) (*internal.User, error) {
	var user internal.User
	err := store.db.GetContext(ctx, &user, `SELECT * FROM users WHERE login=$1`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't get user from db, %w", err)
	}
	return &user, nil
}

func (store *StoreImpl) AddOrder(ctx context.Context, order *internal.Order) (*internal.Order, error) {
	_, err := store.db.NamedExecContext(ctx,
		`INSERT INTO orders (id, create_at, number, accrual, status, user_id) 
			VALUES (:id, :create_at, :number, :accrual, :status, :user_id)`,
		order)
	if err == nil {
		return order, nil
	}

	var pgError *pgconn.PgError
	ok := errors.As(err, &pgError)
	if ok && strings.EqualFold(pgError.SQLState(), "23505") {
		var existOrder internal.Order
		err := store.db.GetContext(ctx,
			&existOrder,
			`SELECT * FROM orders WHERE number = $1`,
			order.Number)

		if err != nil {
			return nil, fmt.Errorf("can't get order from db %w", err)
		}

		if existOrder.UserID == order.UserID {
			return nil, errors2.ErrOrderIsExistThisUser
		} else {
			return nil, errors2.ErrOrderIsExistAnotherUser
		}
	}
	return nil, fmt.Errorf("can't add order from db %w", err)
}

func (store *StoreImpl) GetOrders(ctx context.Context, userID uuid.UUID) (*[]internal.Order, error) {
	var orders []internal.Order
	err := store.db.SelectContext(ctx, &orders,
		`SELECT * FROM orders WHERE user_id=$1 ORDER BY create_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *StoreImpl) GetOrdersNotProcessed(ctx context.Context) (*[]internal.Order, error) {
	var orders []internal.Order
	err := store.db.SelectContext(ctx, &orders,
		`SELECT * FROM orders WHERE status !=$1 AND status !=$2`,
		internal.OrderStatusInvalid, internal.OrderStatusProcessed)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &orders, nil
}

func (store *StoreImpl) UpdateOrder(ctx context.Context, order *internal.Order) error {
	tx := store.db.MustBeginTx(ctx, nil)
	defer tx.Commit()
	if order.Status == internal.OrderStatusProcessed {
		var user internal.User
		err := tx.GetContext(ctx, &user,
			`SELECT u.id, u.bill FROM users AS u INNER JOIN orders AS s ON u.id = s.user_id WHERE s.number=$1`,
			order.Number)
		if err != nil {
			err = tx.Rollback()
			return fmt.Errorf("can't get bill from db %w", err)
		}
		sumBill := user.Bill + order.Accrual
		_, err = tx.ExecContext(ctx, `UPDATE users SET bill = $1 WHERE id = $2`, sumBill, user.ID)
		if err != nil {
			err = tx.Rollback()
			return fmt.Errorf("can't update bill from db %w", err)
		}
	}
	_, err := tx.ExecContext(ctx, `UPDATE orders SET status = $1, accrual=$2 WHERE number = $3`,
		order.Status, order.Accrual, order.Number)
	if err != nil {
		err = tx.Rollback()
		return fmt.Errorf("can't update order from db %w", err)
	}
	return nil
}

func (store *StoreImpl) SaveWithdrawal(ctx context.Context, withdrawal *internal.Withdraw) error {
	tx := store.db.MustBeginTx(ctx, nil)
	defer tx.Commit()
	var user internal.User
	err := tx.GetContext(ctx, &user, `SELECT * FROM users WHERE id=$1`, withdrawal.UserID)
	if err != nil {
		err = tx.Rollback()
		return fmt.Errorf("can't get user from db %w", err)
	}
	if user.Bill < withdrawal.Sum {
		return errors2.ErrNotEnoughAmount
	}
	_, err = tx.NamedExecContext(ctx, `INSERT INTO withdrawals (id, create_at, order_num, sum, user_id) 
											VALUES (:id, :create_at, :order_num, :sum, :user_id)`, withdrawal)
	if err != nil {
		err = tx.Rollback()
		return fmt.Errorf("can't save withdrawal to db %w", err)
	}
	sumBill := user.Bill - withdrawal.Sum
	_, err = tx.ExecContext(ctx, `UPDATE users SET bill = $1 WHERE id = $2`, sumBill, withdrawal.UserID)
	if err != nil {
		err = tx.Rollback()
		return fmt.Errorf("can't update user bill at db %w", err)
	}
	return nil
}

func (store *StoreImpl) GetWithdrawals(ctx context.Context, userID uuid.UUID) (*[]internal.Withdraw, error) {
	var withdrawals []internal.Withdraw
	err := store.db.SelectContext(ctx, &withdrawals, `SELECT * FROM withdrawals WHERE user_id=$1 
                     ORDER BY create_at DESC `, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders from db %w", err)
	}
	return &withdrawals, nil
}

func (store *StoreImpl) GetUser(ctx context.Context, id uuid.UUID) (*internal.User, error) {
	var user internal.User
	err := store.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user from db %w", err)
	}
	return &user, nil
}

type Store interface {
	CheckConnection() error
	AddUser(ctx context.Context, user *internal.User) error
	FindUserByLogin(ctx context.Context, login string) (*internal.User, error)
	AddOrder(ctx context.Context, order *internal.Order) (*internal.Order, error)
	GetOrders(ctx context.Context, userID uuid.UUID) (*[]internal.Order, error)
	GetOrdersNotProcessed(ctx context.Context) (*[]internal.Order, error)
	UpdateOrder(ctx context.Context, order *internal.Order) error
	SaveWithdrawal(ctx context.Context, withdrawal *internal.Withdraw) error
	GetWithdrawals(ctx context.Context, userID uuid.UUID) (*[]internal.Withdraw, error)
	GetUser(ctx context.Context, id uuid.UUID) (*internal.User, error)
}
