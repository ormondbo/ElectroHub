package domain

import (
	"time"
	"uuid"
)

type Product struct {
	ID               uuid.UUID
	Name             string
	Category         string
	Price            float64
	Available_stock  uint
	Last_update_date *time.Time
	Supplier_id      uuid.UUID
	Image_id         uuid.UUID
}
