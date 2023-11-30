package services

import (
	"context"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/auth"
	errors2 "github.com/bonus2k/go-musthave-diploma-tpl/internal/errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	db repositories.Store
}

func NewUserService(storage repositories.Store) *UserService {
	return &UserService{db: storage}
}

func (us *UserService) CreateNewUser(ctx context.Context, user *internal.UserDto) (*internal.User, error) {
	if isIllegalUserArgument(user) {
		return nil, errors2.ErrIllegalUserArgument
	}
	internal.Log.Debug("create user")
	password, err := auth.SignPassword(user.Pass)
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
		return errors2.ErrIllegalOrder
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	order := &internal.Order{ID: uuid.New(), CreateAt: time.Now(), Number: int64(luna), Status: internal.OrderStatusNew, UserID: userID}
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
	ok, err := auth.CheckPassword(user.Pass, login.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors2.ErrWrongAuth
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

func (us *UserService) UpdateOrder(accrual *internal.AccrualDto) error {
	number, err := strconv.Atoi(accrual.Order)
	if err != nil {
		return fmt.Errorf("parse accrual number %s, %w", accrual.Order, err)
	}
	order := &internal.Order{Number: int64(number), Accrual: accrual.Accrual, Status: internal.OrderStatus(accrual.Status)}
	err = us.db.UpdateOrder(context.Background(), order)
	if err != nil {
		return err
	}
	return nil
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

func (us *UserService) GetBalance(ctx context.Context, id string) (*internal.Balance, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	withdrawals, err := us.db.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	user, err := us.db.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var withdrawn float32
	for _, w := range *withdrawals {
		withdrawn = withdrawn + w.Sum
	}

	return &internal.Balance{Current: user.Bill, Withdrawn: withdrawn}, nil
}

func (us *UserService) AddWithdraw(ctx context.Context, dto internal.WithdrawDto, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	luna, ok := isLuna(dto.Order)
	if !ok {
		return errors2.ErrIllegalOrder
	}

	withdraw := &internal.Withdraw{
		ID:       uuid.New(),
		CreateAt: time.Now(),
		Order:    int64(luna),
		Sum:      dto.Sum,
		UserID:   userID,
	}

	return us.db.SaveWithdrawal(ctx, withdraw)
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
