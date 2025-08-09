package models

type Balance struct {
	Login     *string
	Current   *float64 `json:"current"`
	Withdrawn *float64 `json:"withdrawn"`
}
