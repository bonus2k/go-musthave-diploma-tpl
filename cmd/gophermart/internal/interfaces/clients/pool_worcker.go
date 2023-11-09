package clients

import (
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/services"
	"go.uber.org/zap"
	"time"
)

type PoolWorker struct {
	client      *ClientAccrual
	serviceUser *services.UserService
	orderIn     chan string
	Err         chan error
}

func NewPoolWorker(client *ClientAccrual, serviceUser *services.UserService) *PoolWorker {
	ordersIn := make(chan string, 10)
	err := make(chan error)
	return &PoolWorker{client: client, serviceUser: serviceUser, orderIn: ordersIn, Err: err}
}

func (p *PoolWorker) StarIntegration(countWorker int, requestTime *time.Ticker) {

	pauses := make([]chan struct{}, 0)
	for i := 0; i < countWorker; i++ {
		name := i
		p := p.worker(name)
		pauses = append(pauses, p)
	}

	go func() {
		for range requestTime.C {
			numbers, err := p.serviceUser.GetOrdersNotProcessed()
			if err != nil {
				internal.Log.Error("error start workers of integration")
				break
			}
			for _, n := range numbers {
				number := n
				p.orderIn <- number
			}
		}
	}()

	for err := range p.Err {
		internal.Log.Error("error integration", zap.Error(err))
		if errors.Is(err, ErrTooManyRequests) {
			go func() {
				for _, pause := range pauses {
					ch := pause
					ch <- struct{}{}
				}
			}()
		}
	}
}

func (p *PoolWorker) worker(name int) chan struct{} {
	pause := make(chan struct{})
	go func() {
		defer close(pause)
		for {
			select {
			case order := <-p.orderIn:
				internal.Logf.Debugf("worker %d, order %s send request to accrual services", name, order)
				accrual, err := p.client.CheckAccrual(order)
				if err != nil {
					p.Err <- fmt.Errorf("error worker %d %w", name, err)
					break
				}
				internal.Logf.Debugf("worker %d, save accrual in order", name)
				p.serviceUser.UpdateOrder(accrual)
			case <-pause:
				internal.Logf.Debugf("worker %d do pause", name)
				time.Sleep(2 * time.Minute)
			}
		}
	}()
	return pause
}
