package order

import (
	"database/sql"
	"encoding/json"
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
	ID              uint           `json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	ProductName     string         `json:"productName"`
	Supplier        string         `json:"supplier"`
	AdditionalNotes sql.NullString `json:"-"`
	Status          OrderStatus    `json:"status"`
	Quantity        int            `json:"quantity"`
	UserID          uint           `json:"userId"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	additionalNotes := ""
	if o.AdditionalNotes.Valid {
		additionalNotes = o.AdditionalNotes.String
	}

	return json.Marshal(&struct {
		*Alias
		AdditionalNotes string `json:"additionalNotes"`
	}{
		Alias:           (*Alias)(&o),
		AdditionalNotes: additionalNotes,
	})
}
