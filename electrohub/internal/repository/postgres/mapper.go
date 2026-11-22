package postgres

import "electrohub/internal/domain"

// ======================================================================================================
// Address
// ======================================================================================================

func toDBAddress(da domain.Address) dbAddress {
	return dbAddress{
		id:      da.ID,
		country: da.Country,
		city:    da.City,
		street:  da.Street,
	}
}

func toDomainAddress(dba dbAddress) domain.Address {
	return domain.Address{
		ID:      dba.id,
		Country: dba.country,
		City:    dba.city,
		Street:  dba.street,
	}
}

// ======================================================================================================
// Client
// ======================================================================================================

func toDBClient(dc domain.Client) dbClient {
	dbc := dbClient{
		id:                dc.ID,
		name:              dc.Name,
		surname:           dc.Surname,
		birthday:          dc.Birthday,
		registration_date: dc.Registration_date,
		address_id:        dc.Address.ID,
	}

	switch dc.Gender {
	case domain.GenderMale:
		dbc.gender = "male"
	case domain.GenderFemale:
		dbc.gender = "female"
	default:
		dbc.gender = "other"
	}
	return dbc
}

func toDomainClient(dbc dbClient) domain.Client {
	dc := domain.Client{
		ID:                dbc.id,
		Name:              dbc.name,
		Surname:           dbc.surname,
		Birthday:          dbc.birthday,
		Registration_date: dbc.registration_date,
		Address:           domain.Address{ID: dbc.address_id},
	}
	switch dbc.gender {
	case "male":
		dc.Gender = domain.GenderMale
	case "female":
		dc.Gender = domain.GenderFemale
	case "helicopter":
		dc.Gender = domain.GenderHelicopter
	default:
		dc.Gender = domain.GenderOther
	}
	return dc
}

// ======================================================================================================
// Supplier
// ======================================================================================================

func toDBSupplier(ds domain.Supplier) dbSupplier {
	return dbSupplier{
		id:           ds.ID,
		name:         ds.Name,
		address_id:   ds.Address.ID,
		phone_number: ds.Phone_number,
	}
}

func toDomainSupplier(dbs dbSupplier) domain.Supplier {
	return domain.Supplier{
		ID:           dbs.id,
		Name:         dbs.name,
		Address:      domain.Address{ID: dbs.address_id},
		Phone_number: dbs.phone_number,
	}
}

// ======================================================================================================
// Image
// ======================================================================================================

func toDBImage(di domain.Image) dbImage {
	return dbImage{
		id:   di.ID,
		data: di.Data,
	}
}

func toDomainImage(dbi dbImage) domain.Image {
	return domain.Image{
		ID:   dbi.id,
		Data: dbi.data,
	}
}

// ======================================================================================================
// Product
// ======================================================================================================

func toDBProduct(dp domain.Product) dbProduct {
	return dbProduct{
		id:               dp.ID,
		name:             dp.Name,
		category:         dp.Category,
		price:            dp.Price,
		available_stock:  dp.Available_stock,
		last_update_date: dp.Last_update_date,
		supplier_id:      dp.Supplier_id,
		image_id:         dp.Image_id,
	}
}

func toDomainProduct(dbp dbProduct) domain.Product {
	return domain.Product{
		ID:               dbp.id,
		Name:             dbp.name,
		Category:         dbp.category,
		Price:            dbp.price,
		Available_stock:  dbp.available_stock,
		Last_update_date: dbp.last_update_date,
		Supplier_id:      dbp.supplier_id,
		Image_id:         dbp.image_id,
	}
}
