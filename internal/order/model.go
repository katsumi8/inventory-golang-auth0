package order

import (
	"time"
)

type OrderStatus string

const (
	Ordered   OrderStatus = "Ordered"
	OnProcess OrderStatus = "OnProcess"
	Pending   OrderStatus = "Pending"
	Delivered OrderStatus = "Delivered"
	Canceled  OrderStatus = "Canceled"
)

type Order struct {
	ID              uint        `json:"id"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	ProductName     string      `json:"product_name"`
	Supplier        string      `json:"supplier"`
	AdditionalNotes string      `json:"additionalNotes,omitempty"`
	Status          OrderStatus `json:"status"`
	Quantity        int         `json:"quantity"`
	UserID          uint        `json:"userId"`
}
