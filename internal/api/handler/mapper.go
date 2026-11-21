package handler

import "electrohub/internal/api/domain"

// ======================================================================================================
// Address
// ======================================================================================================

func toDomainAddress(req addressRequest) domain.Address {
	return domain.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
}

// ======================================================================================================
// Client
// ======================================================================================================

func toDomainClient(req createClientRequest) domain.Client {
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

func toGetClientResponse(client domain.Client) getClientResponse {
	response := getClientResponse{
		ID:               client.ID,
		Name:             client.Name,
		Surname:          client.Surname,
		Birthday:         client.Birthday,
		RegistrationDate: client.Registration_date,
		Address: addressResponse{
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

func toGetClientsResponse(clients []domain.Client) getClientsResponse {
	response := getClientsResponse{
		Clients: make([]getClientResponse, 0, len(clients)),
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

func toDomainProduct(req createProductRequest) domain.Product {
	return domain.Product{
		Name:            req.Name,
		Category:        req.Category,
		Price:           req.Price,
		Available_stock: req.AvailableStock,
		Supplier_id:     req.SupplierID,
		Image_id:        req.ImageID,
	}
}

func toGetProductResponse(product domain.Product) getProductResponse {
	return getProductResponse{
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

func toGetProductsResponse(products []domain.Product) getProductsResponse {
	response := getProductsResponse{
		Products: make([]getProductResponse, 0, len(products)),
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

func toDomainSupplier(req createSupplierRequest) domain.Supplier {
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

func toGetSupplierResponse(supplier domain.Supplier) getSupplierResponse {
	return getSupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		PhoneNumber: supplier.Phone_number,
		Address: addressResponse{
			ID:      supplier.Address.ID,
			Country: supplier.Address.Country,
			City:    supplier.Address.City,
			Street:  supplier.Address.Street,
		},
	}
}

func toGetSuppliersResponse(suppliers []domain.Supplier) getSuppliersResponse {
	response := getSuppliersResponse{
		Suppliers: make([]getSupplierResponse, 0, len(suppliers)),
	}

	for _, supplier := range suppliers {
		response.Suppliers = append(
			response.Suppliers,
			toGetSupplierResponse(supplier),
		)
	}

	return response
}
