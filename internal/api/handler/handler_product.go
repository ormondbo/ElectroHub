package handler

import (
	"math"
	"net/http"
	"uuid"
)

// @Summary Create a product
// @Description Creates a product for an existing supplier.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body createProductRequest true "Product data"
// @Success 201 {object} getProductResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [post]
func (api *API) createProduct(writer http.ResponseWriter, request *http.Request) {
	var req createProductRequest
	if err := decodeJSON(writer, request, &req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if err := validateProduct(req); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	product, err := api.product.CreateProduct(request.Context(), toDomainProduct(req))
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, toGetProductResponse(product))
}

// @Summary Get a product by ID
// @Description Returns the product identified by its UUID.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} getProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [get]
func (api *API) getProductByID(writer http.ResponseWriter, request *http.Request) {
	productID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	product, err := api.product.GetProductByID(request.Context(), productID)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toGetProductResponse(product))
}

// @Summary List products with available stock
// @Description Returns products whose available stock is greater than zero.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Success 200 {object} getProductsResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [get]
func (api *API) getAvailableProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := api.product.GetAvailableProducts(request.Context())
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toGetProductsResponse(products))
}

// @Summary Decrease product stock
// @Description Decreases the stock by the requested amount when sufficient stock is available.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Param id path string true "Product UUID"
// @Param amount body decreaseProductStockRequest true "Quantity to decrease"
// @Success 204 "Stock decreased"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id}/stock [patch]
func (api *API) decreaseProductStock(writer http.ResponseWriter, request *http.Request) {
	productID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	var req decreaseProductStockRequest
	if err := decodeJSON(writer, request, &req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if req.Amount == 0 {
		writeError(writer, http.StatusBadRequest, "amount must be greater than zero")
		return
	}
	if req.Amount > math.MaxInt32 {
		writeError(writer, http.StatusBadRequest, "amount is too large")
		return
	}
	err = api.product.DecreaseProductStock(request.Context(), productID, req.Amount)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete a product
// @Description Deletes the product identified by its UUID.
// @Tags Products
// @Security BearerAuth
// @Param id path string true "Product UUID"
// @Success 204 "Deleted"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
func (api *API) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	productID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	if err := api.product.DeleteProduct(request.Context(), productID); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
