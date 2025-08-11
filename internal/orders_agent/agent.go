package agent

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/config"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"time"
)

const orderChunkSizeDefault = 5

type OrdersAgent struct {
	orderService   *service.OrderService
	httpClient     *StatusUpdaterClient
	config         *config.Config
	orderChunkSize uint
}

func NewOrdersAgent(orderService *service.OrderService, config *config.Config) *OrdersAgent {
	return &OrdersAgent{
		orderService:   orderService,
		config:         config,
		httpClient:     NewClientAgent(*config.AccrualSystemAddress),
		orderChunkSize: orderChunkSizeDefault,
	}
}

func (o *OrdersAgent) Run(ctx context.Context) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	pollTicker := time.NewTicker(time.Duration(*o.config.PollInterval) * time.Second)
	defer pollTicker.Stop()

	resultCh := make(chan models.OrderResult)
	defer close(resultCh)

	orderNumberCh := o.orderGenerator(ctx, doneCh, pollTicker)
	orderNewStatusCh := o.OrdersStatusGenerator(ctx, doneCh, orderNumberCh)
	go o.updateStatusOrder(ctx, doneCh, orderNewStatusCh, resultCh)
	for res := range resultCh {
		if res.Err != nil {
			logger.Log.Warn(res.Err.Error())
		}
	}
}

func (o *OrdersAgent) updateStatusOrder(ctx context.Context, doneCh chan struct{}, orderNewStatusCh <-chan *models.Order, resultCh chan<- models.OrderResult) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Worker: Context done, exiting.")
			return
		case <-doneCh:
			return
		case order := <-orderNewStatusCh:
			err := o.orderService.UpdateStatus(ctx, order)
			if err != nil {
				resultCh <- models.OrderResult{Err: err}
			}
		}
	}
}

func (o *OrdersAgent) orderGenerator(ctx context.Context, doneCh chan struct{}, pollTicker *time.Ticker) <-chan models.OrderResult {
	orderCh := make(chan models.OrderResult, o.orderChunkSize)
	go func() {
		defer close(orderCh)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			case <-pollTicker.C:
				orders, err := o.orderService.GetNotAcceptedOrderNumbers(ctx, o.orderChunkSize)
				for _, order := range orders {
					orderCh <- models.OrderResult{Order: order, Err: err}
				}
			}
		}
	}()
	return orderCh
}

func (o *OrdersAgent) OrdersStatusGenerator(ctx context.Context, doneCh chan struct{}, orderCh <-chan models.OrderResult) <-chan *models.Order {
	orderNewStatusCh := make(chan *models.Order, o.orderChunkSize)
	go func() {
		defer close(orderNewStatusCh)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			default:
				for w := 1; w <= int(*o.config.RateLimit); w++ {
					go o.getOrdersStatus(ctx, doneCh, orderCh, orderNewStatusCh)
				}
			}
		}
	}()
	return orderNewStatusCh
}

func (o *OrdersAgent) getOrdersStatus(ctx context.Context, doneCh chan struct{}, orderCh <-chan models.OrderResult, orderNewStatusCh chan<- *models.Order) {
	select {
	case <-ctx.Done():
		logger.Log.Info("Worker: Context done, exiting.")
		return
	case <-doneCh:
		return
	case orderRes := <-orderCh:
		res := o.httpClient.getOrderStatus(orderRes.Order.Number)
		if res.Err != nil {
			logger.Log.Warn(res.Err.Error())
		}
		orderNewStatusCh <- res.Order
	}
}
