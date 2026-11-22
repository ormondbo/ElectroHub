package domain

import "uuid"

type Supplier struct {
	ID           uuid.UUID
	Name         string
	Address      Address
	Phone_number string
}
