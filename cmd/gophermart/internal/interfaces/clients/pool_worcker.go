package clients

import (
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/services"
	"go.uber.org/zap"
)

type PoolWorker struct {
	client      *ClientAccrual
	serviceUser *services.UserService
	ordersIn    chan string
	Done        chan struct{}
	err         chan error
}

func NewPoolWorker(client *ClientAccrual, serviceUser *services.UserService) *PoolWorker {
	ordersIn := make(chan string, 10)
	done := make(chan struct{})
	err := make(chan error)
	return &PoolWorker{client: client, serviceUser: serviceUser, ordersIn: ordersIn, Done: done, err: err}
}

func (p *PoolWorker) StarIntegration(countWorker int) {
	numbers, err := p.serviceUser.GetOrdersNotProcessed()
	if err != nil {
		internal.Log.Error("error start workers of integration")
		return
	}

	go func() {
		for i := 0; i < countWorker; i++ {
			name := i
			go p.worker(name)
		}
	}()

	go func() {
		for _, n := range numbers {
			number := n
			p.ordersIn <- number
		}
	}()
	for err := range p.err {
		internal.Log.Error("error integration", zap.Error(err))
		if errors.Is(err, ErrTooManyRequests) {
			close(p.Done)
			close(p.err)
		}
	}
	close(p.ordersIn)
}

func (p *PoolWorker) worker(name int) {
close:
	for {
		select {
		case order := <-p.ordersIn:
			internal.Logf.Debugf("worker %d, order %s send request to accrual services", name, order)
			accrual, err := p.client.CheckAccrual(order)
			if err != nil {
				p.err <- fmt.Errorf("error worker %d %w", name, err)
				break
			}
			internal.Logf.Debugf("worker %d, save accrual in order", name)
			p.serviceUser.UpdateOrder(accrual)
		case <-p.Done:
			internal.Logf.Debugf("worker %d stoped")
			break close
		}
	}
}
