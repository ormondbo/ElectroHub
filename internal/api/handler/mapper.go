package httpapi

import (
	"shop/internal/domain"
)

// ======================================================================================================
// Address
// ======================================================================================================

func toDomainAddress(req AddressRequest) domain.Address {
	return domain.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
}

// ======================================================================================================
// Client
// ======================================================================================================

func toDomainClient(req CreateClientRequest) domain.Client {
	client := domain.Client{
		Name:     req.Name,
		Surname:  req.Surname,
		Birthday: req.Birthday,
		Address: domain.Address{
			Country: req.Address.Country,
			City:    req.Address.City,
			Street:  req.Address.Street,
		},
	}
	switch req.Gender {
	case "male":
		client.Gender = domain.GenderMale
	case "female":
		client.Gender = domain.GenderFemale
	case "helicopter":
		client.Gender = domain.GenderHelicopter
	default:
		client.Gender = domain.GenderOther
	}
	return client
}

func toGetClientResponse(client domain.Client) GetClientResponse {
	response := GetClientResponse{
		ID:               client.ID,
		Name:             client.Name,
		Surname:          client.Surname,
		Birthday:         client.Birthday,
		RegistrationDate: client.Registration_date,
		Address: AddressResponse{
			ID:      client.Address.ID,
			Country: client.Address.Country,
			City:    client.Address.City,
			Street:  client.Address.Street,
		},
	}
	switch client.Gender {
	case domain.GenderMale:
		response.Gender = "male"
	case domain.GenderFemale:
		response.Gender = "female"
	case domain.GenderHelicopter:
		response.Gender = "helicopter"
	default:
		response.Gender = "other"
	}
	return response
}

func toGetClientsResponse(clients []domain.Client) GetClientsResponse {
	response := GetClientsResponse{
		Clients: make([]GetClientResponse, 0, len(clients)),
	}
	for _, client := range clients {
		response.Clients = append(
			response.Clients,
			toGetClientResponse(client),
		)
	}
	return response
}

// ======================================================================================================
// Product
// ======================================================================================================

func toDomainProduct(req CreateProductRequest) domain.Product {
	return domain.Product{
		Name:            req.Name,
		Category:        req.Category,
		Price:           req.Price,
		Available_stock: req.AvailableStock,
		Supplier_id:     req.SupplierID,
		Image_id:        req.ImageID,
	}
}

func toGetProductResponse(product domain.Product) GetProductResponse {
	return GetProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.Available_stock,
		LastUpdateDate: product.Last_update_date,
		SupplierID:     product.Supplier_id,
		ImageID:        product.Image_id,
	}
}

func toGetProductsResponse(products []domain.Product) GetProductsResponse {
	response := GetProductsResponse{
		Products: make([]GetProductResponse, 0, len(products)),
	}
	for _, product := range products {
		response.Products = append(
			response.Products,
			toGetProductResponse(product),
		)
	}
	return response
}

// ======================================================================================================
// Suppliers
// ======================================================================================================

func toDomainSupplier(req CreateSupplierRequest) domain.Supplier {
	return domain.Supplier{
		Name:         req.Name,
		Phone_number: req.PhoneNumber,
		Address: domain.Address{
			Country: req.Address.Country,
			City:    req.Address.City,
			Street:  req.Address.Street,
		},
	}
}

func toGetSupplierResponse(supplier domain.Supplier) GetSupplierResponse {
	return GetSupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		PhoneNumber: supplier.Phone_number,
		Address: AddressResponse{
			ID:      supplier.Address.ID,
			Country: supplier.Address.Country,
			City:    supplier.Address.City,
			Street:  supplier.Address.Street,
		},
	}
}

func toGetSuppliersResponse(suppliers []domain.Supplier) GetSuppliersResponse {
	response := GetSuppliersResponse{
		Suppliers: make([]GetSupplierResponse, 0, len(suppliers)),
	}

	for _, supplier := range suppliers {
		response.Suppliers = append(
			response.Suppliers,
			toGetSupplierResponse(supplier),
		)
	}

	return response
}
