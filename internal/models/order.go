package models

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
)

type Order struct {
	Login      *string
	Number     *string     `json:"number"`
	Status     *string     `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt *CustomTime `json:"uploaded_at"`
}

type OrderResult struct {
	Order *Order
	Err   error
}
