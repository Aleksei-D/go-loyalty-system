package models

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusInvalid    = "INVALID"
)

type Order struct {
	Login      string     `json:"-"`
	Number     string     `json:"number"`
	Status     string     `json:"status"`
	Accrual    *float64   `json:"accrual,omitempty"`
	UploadedAt CustomTime `json:"uploaded_at"`
}

type OrderStatusResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

func (o *OrderStatusResponse) ToOrder() *Order {
	return &Order{
		Number:  o.Order,
		Status:  o.Status,
		Accrual: o.Accrual,
	}
}
