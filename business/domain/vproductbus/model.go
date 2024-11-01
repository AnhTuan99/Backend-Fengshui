package vproductbus

import (
	"time"

	"fengshui.com/back-fengshui/business/types/money"
	"fengshui.com/back-fengshui/business/types/name"
	"fengshui.com/back-fengshui/business/types/quantity"
	"github.com/google/uuid"
)

// Product represents an individual product with extended information.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
	UserName    name.Name
}
