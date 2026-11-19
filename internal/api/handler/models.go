package httpapi

import (
	"time"
	"uuid"
)

// ======================================================================================================
// Auth
// ======================================================================================================

type registerRequest struct {
	Email     string `json:"email" binding:"required" format:"email" example:"user@example.com"`
	FirstName string `json:"first_name" binding:"required" example:"Иван"`
	LastName  string `json:"last_name" binding:"required" example:"Иванов"`
	Phone     string `json:"phone" binding:"required" example:"+79990000000"`
	Password  string `json:"password" binding:"required" example:"strong-password"`
}

type authRequest struct {
	Email    string `json:"email" binding:"required" format:"email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"strong-password"`
}

type resetRequest struct {
	Email string `json:"email" binding:"required" format:"email" example:"user@example.com"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"current-password"`
	NewPassword string `json:"new_password" binding:"required" example:"new-password"`
}

type changePasswordResponse struct {
	Message string `json:"message" example:"password changed successfully"`
}

// ======================================================================================================
// Addresses
// ======================================================================================================

type AddressRequest struct {
	Country string `json:"country" binding:"required" example:"Россия"`
	City    string `json:"city" binding:"required" example:"Москва"`
	Street  string `json:"street" binding:"required" example:"ул. Тверская, д. 1"`
}

type AddressResponse struct {
	ID      uuid.UUID `json:"id" example:"6b8a1d97-859b-4c5e-913a-01bfa34b201c"`
	Country string    `json:"country" example:"Россия"`
	City    string    `json:"city" example:"Москва"`
	Street  string    `json:"street" example:"ул. Тверская, д. 1"`
}

// ======================================================================================================
// Clients
// ======================================================================================================

type CreateClientRequest struct {
	Name     string         `json:"name" binding:"required" example:"Александр"`
	Surname  string         `json:"surname" binding:"required" example:"Иванов"`
	Birthday time.Time      `json:"birthday" binding:"required" example:"1990-05-20T00:00:00Z"`
	Gender   string         `json:"gender" binding:"required" enums:"male,female,helicopter,other" example:"other"`
	Address  AddressRequest `json:"address" binding:"required"`
}

type GetClientResponse struct {
	ID               uuid.UUID       `json:"id" example:"f8c93d19-130e-46d3-86dd-7f84880ed83d"`
	Name             string          `json:"name" example:"Александр"`
	Surname          string          `json:"surname" example:"Иванов"`
	Birthday         time.Time       `json:"birthday" example:"1990-05-20T00:00:00Z"`
	Gender           string          `json:"gender" example:"other"`
	RegistrationDate time.Time       `json:"registration_date" example:"2025-01-01T12:00:00Z"`
	Address          AddressResponse `json:"address"`
}

type GetClientsResponse struct {
	Clients []GetClientResponse `json:"clients"`
}

type UpdateClientAddressRequest struct {
	Address AddressRequest `json:"address" binding:"required"`
}

// ======================================================================================================
// Products
// ======================================================================================================

type CreateProductRequest struct {
	Name           string    `json:"name" binding:"required" example:"Электрический чайник"`
	Category       string    `json:"category" binding:"required" example:"Бытовая техника"`
	Price          float64   `json:"price" binding:"required" example:"49.99"`
	AvailableStock uint      `json:"available_stock" example:"12"`
	SupplierID     uuid.UUID `json:"supplier_id" binding:"required" format:"uuid" example:"f8c93d19-130e-46d3-86dd-7f84880ed83d"`
	ImageID        uuid.UUID `json:"image_id" format:"uuid" example:"00000000-0000-0000-0000-000000000000"`
}

type GetProductResponse struct {
	ID             uuid.UUID  `json:"id" example:"f8c93d19-130e-46d3-86dd-7f84880ed83d"`
	Name           string     `json:"name" example:"Электрический чайник"`
	Category       string     `json:"category" example:"Бытовая техника"`
	Price          float64    `json:"price" example:"49.99"`
	AvailableStock uint       `json:"available_stock" example:"12"`
	LastUpdateDate *time.Time `json:"last_update_date"`
	SupplierID     uuid.UUID  `json:"supplier_id" format:"uuid" example:"f8c93d19-130e-46d3-86dd-7f84880ed83d"`
	ImageID        uuid.UUID  `json:"image_id" format:"uuid" example:"00000000-0000-0000-0000-000000000000"`
}

type GetProductsResponse struct {
	Products []GetProductResponse `json:"products"`
}

type DecreaseProductStockRequest struct {
	Amount uint `json:"amount" binding:"required" example:"2"`
}

// ======================================================================================================
// Suppliers
// ======================================================================================================

type CreateSupplierRequest struct {
	Name        string         `json:"name" binding:"required" example:"ООО «ТехноСнаб»"`
	PhoneNumber string         `json:"phone_number" binding:"required" example:"+7-495-123-45-67"`
	Address     AddressRequest `json:"address" binding:"required"`
}

type GetSupplierResponse struct {
	ID          uuid.UUID       `json:"id" example:"f8c93d19-130e-46d3-86dd-7f84880ed83d"`
	Name        string          `json:"name" example:"ООО «ТехноСнаб»"`
	PhoneNumber string          `json:"phone_number" example:"+7-495-123-45-67"`
	Address     AddressResponse `json:"address"`
}

type GetSuppliersResponse struct {
	Suppliers []GetSupplierResponse `json:"suppliers"`
}

type UpdateSupplierAddressRequest struct {
	Address AddressRequest `json:"address" binding:"required"`
}
