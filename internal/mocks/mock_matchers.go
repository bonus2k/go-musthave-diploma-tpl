package mock

import (
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
)

type MatchOrder struct {
	*internal.Order
}

func (m *MatchOrder) Matches(o interface{}) bool {
	order := o.(*internal.Order)
	return m.Number == order.Number
}

func (m *MatchOrder) String() string {
	return fmt.Sprintf("Order request matcher %v", m.Order)
}

type MatchWithdraw struct {
	*internal.Withdraw
}

func (w *MatchWithdraw) Matches(o interface{}) bool {
	withdraw := o.(*internal.Withdraw)
	return w.Order == withdraw.Order
}

func (w *MatchWithdraw) String() string {
	return fmt.Sprintf("Withdraw request matcher %v", w.Withdraw)
}

type MatchUser struct {
	*internal.User
}

func (u *MatchUser) Matches(o interface{}) bool {
	user := o.(*internal.User)
	return u.Login == user.Login
}

func (u *MatchUser) String() string {
	return fmt.Sprintf("User request matcher %v", u.User)
}
