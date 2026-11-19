package httpapi

import (
	"electrohub/internal/api/service"
	"log"
	"net/http"
)

type API struct {
	client   *service.ClientService
	product  *service.ProductService
	supplier *service.SupplierService
	image    *service.ImageService
	auth     AuthClient

	mux          *http.ServeMux
	protectedMux *http.ServeMux
}

func NewAPI(client *service.ClientService, product *service.ProductService,
	supplier *service.SupplierService, image *service.ImageService, auth AuthClient,
) *API {
	api := &API{
		client:   client,
		product:  product,
		supplier: supplier,
		image:    image,
		auth:     auth,

		mux:          http.NewServeMux(),
		protectedMux: http.NewServeMux(),
	}

	api.routes()
	return api
}

// обозвать по-другому
func (api *API) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf(
				"panic while handling %s %s: %v",
				request.Method,
				request.URL.Path,
				recovered,
			)
			writeError(
				writer,
				http.StatusInternalServerError,
				"internal server error",
			)
		}
	}()

	api.mux.ServeHTTP(writer, request)
}

// Поправить
func (api *API) routes() {
	api.mux.HandleFunc("GET /swagger/", serveSwaggerUI)
	api.mux.HandleFunc("GET /swagger/doc.json", serveSwaggerSpec)

	api.mux.HandleFunc("POST /register", api.register)
	api.mux.HandleFunc("POST /auth", api.authenticate)
	api.mux.HandleFunc("POST /reset", api.resetPassword)
	api.mux.Handle("PUT /auth/change-password", api.authMiddleware(http.HandlerFunc(api.changePassword)))
	api.mux.Handle("/api/v1/", api.authMiddleware(api.protectedMux))

	api.protectedMux.HandleFunc("POST /api/v1/clients", api.createClient)
	api.protectedMux.HandleFunc("DELETE /api/v1/clients/{id}", api.deleteClient)
	api.protectedMux.HandleFunc("GET /api/v1/clients", api.getClients)
	api.protectedMux.HandleFunc("PUT /api/v1/clients/{id}/address", api.updateClientAddress)

	api.protectedMux.HandleFunc("POST /api/v1/products", api.createProduct)
	api.protectedMux.HandleFunc("GET /api/v1/products/{id}", api.getProductByID)
	api.protectedMux.HandleFunc("GET /api/v1/products", api.getAvailableProducts)
	api.protectedMux.HandleFunc("PATCH /api/v1/products/{id}/stock", api.decreaseProductStock)
	api.protectedMux.HandleFunc("DELETE /api/v1/products/{id}", api.deleteProduct)

	api.protectedMux.HandleFunc("POST /api/v1/suppliers", api.createSupplier)
	api.protectedMux.HandleFunc("PUT /api/v1/suppliers/{id}/address", api.updateSupplierAddress)
	api.protectedMux.HandleFunc("DELETE /api/v1/suppliers/{id}", api.deleteSupplier)
	api.protectedMux.HandleFunc("GET /api/v1/suppliers", api.getSuppliers)
	api.protectedMux.HandleFunc("GET /api/v1/suppliers/{id}", api.getSupplierByID)

	api.protectedMux.HandleFunc("POST /api/v1/images", api.addImage)
	api.protectedMux.HandleFunc("PUT /api/v1/images/{id}", api.updateImage)
	api.protectedMux.HandleFunc("DELETE /api/v1/images/{id}", api.deleteImage)
	api.protectedMux.HandleFunc("GET /api/v1/images/{id}", api.getImageByID)
	api.protectedMux.HandleFunc("GET /api/v1/products/{id}/image", api.getImageByProductID)
}
