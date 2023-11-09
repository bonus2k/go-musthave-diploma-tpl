package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	db *repositories.Store
}

func NewUserService(storage *repositories.Store) *UserService {
	return &UserService{db: storage}
}

func (us *UserService) CreateNewUser(ctx context.Context, user *internal.UserDto) (*internal.User, error) {
	if isIllegalUserArgument(user) {
		return nil, ErrIllegalUserArgument
	}
	internal.Log.Debug("create user")
	password, err := utils.SignPassword(user.Pass)
	if err != nil {
		return nil, fmt.Errorf("can't create password, %w", err)
	}
	entity := &internal.User{ID: uuid.New(), CreateAt: time.Now(), Login: user.Login, Password: password}
	if err := us.db.AddUser(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (us *UserService) AddOrder(ctx context.Context, id string, orderID string) error {
	luna, ok := isLuna(orderID)
	if !ok {
		return ErrIllegalOrder
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	order := &internal.Order{ID: uuid.New(), CreateAt: time.Now(), Number: int64(luna), Status: internal.NEW, UserID: userID}
	_, err = us.db.AddOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) LoginUser(ctx context.Context, user internal.UserDto) (*uuid.UUID, error) {
	login, err := us.db.FindUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	ok, err := utils.CheckPassword(user.Pass, login.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrWrongAuth
	}
	return &login.ID, nil
}

func (us *UserService) GetOrdersNotProcessed() ([]string, error) {
	orders, err := us.db.GetOrdersNotProcessed(context.Background())
	if err != nil {
		return nil, err
	}
	numbers := make([]string, 0)
	for _, order := range *orders {
		numbers = append(numbers, strconv.FormatInt(order.Number, 10))
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	return numbers, nil
}

func (us *UserService) GetOrders(ctx context.Context, id string) (*[]internal.OrderDto, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	orders, err := us.db.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	ordersDto := make([]internal.OrderDto, 0)
	for _, order := range *orders {
		d := internal.OrderDto{
			Number:  strconv.FormatInt(order.Number, 10),
			Status:  string(order.Status),
			Accrual: order.Accrual,
			Upload:  order.CreateAt,
		}
		ordersDto = append(ordersDto, d)
	}
	return &ordersDto, nil
}

func (us *UserService) UpdateOrder(accrual *internal.AccrualDto) {
	number, err := strconv.Atoi(accrual.Order)
	if err != nil {
		internal.Logf.Errorf("parse accrual number %s", accrual.Order)
		return
	}
	order := &internal.Order{Number: int64(number), Accrual: accrual.Accrual, Status: internal.Value(accrual.Status)}
	err = us.db.UpdateOrder(context.Background(), order)
	if err != nil {
		internal.Log.Error("update order", zap.Error(err))
	}
}

func (us *UserService) GetWithdrawals(ctx context.Context, id string) (*[]internal.WithdrawDto, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	withdrawals, err := us.db.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]internal.WithdrawDto, 0)
	for _, withdraw := range *withdrawals {
		dto := internal.WithdrawDto{
			Order:    strconv.FormatInt(withdraw.Order, 10),
			Sum:      withdraw.Sum,
			CreateAt: withdraw.CreateAt,
		}
		dtos = append(dtos, dto)
	}
	return &dtos, nil
}

func isLuna(order string) (int, bool) {
	number, err := strconv.Atoi(order)
	if err != nil {
		return number, false
	}
	return number, (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int
	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func isIllegalUserArgument(user *internal.UserDto) bool {
	trimLogin := strings.TrimSpace(user.Login)
	trimPassword := strings.TrimSpace(user.Pass)
	if len(trimLogin) == 0 || len(trimPassword) == 0 {
		internal.Logf.Errorf("parameters is wrong; login=%s, password=%s",
			user.Login, user.Pass)
		return true
	}
	return false
}

var ErrIllegalUserArgument = errors.New("illegal user argument")
var ErrIllegalOrder = errors.New("illegal order")
var ErrWrongAuth = errors.New("wrong authorization")
