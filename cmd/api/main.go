package main

import (
	"electrohub/internal/api/di"

	"go.uber.org/fx"
)

//ПЕРЕПИСАТЬ! Старая версия

// @title Store API
// @version 1.0
// @description REST API for managing clients, products, suppliers, and product images.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Send the JWT access token with the Bearer authorization scheme.
// @tag.name Authentication
// @tag.description Registration, login, password changes, and recovery.
// @tag.name Clients
// @tag.description Operations for managing store clients.
// @tag.name Products
// @tag.description Operations for managing products and stock.
// @tag.name Suppliers
// @tag.description Operations for managing suppliers.
// @tag.name Images
// @tag.description Operations for uploading and retrieving product images.
func main() {
	fx.New(di.Module).Run()
}
