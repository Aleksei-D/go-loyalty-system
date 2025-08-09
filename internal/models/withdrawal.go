package models

type Withdrawal struct {
	Login       *string
	OrderNumber *string     `json:"order"`
	Sum         *float64    `json:"sum"`
	ProcessedAt *CustomTime `json:"processed_at"`
}
