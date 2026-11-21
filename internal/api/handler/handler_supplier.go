package handler

import (
	"net/http"
	"uuid"
)

// @Summary Create a supplier
// @Description Creates a supplier and its address when the address does not already exist.
// @Tags Suppliers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param supplier body createSupplierRequest true "Supplier data"
// @Success 201 {object} getSupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [post]
func (api *API) createSupplier(writer http.ResponseWriter, request *http.Request) {
	var req createSupplierRequest
	if err := decodeJSON(writer, request, &req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if err := validateSupplier(req); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	supplier, err := api.supplier.CreateSupplier(request.Context(), toDomainSupplier(req))
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, toGetSupplierResponse(supplier))
}

// @Summary Change a supplier's address
// @Description Replaces the address associated with the specified supplier.
// @Tags Suppliers
// @Security BearerAuth
// @Accept json
// @Param id path string true "Supplier UUID"
// @Param address body updateSupplieraddressRequest true "New address"
// @Success 204 "Updated"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id}/address [put]
func (api *API) updateSupplierAddress(writer http.ResponseWriter, request *http.Request) {
	supplierID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid supplier ID")
		return
	}
	var req updateSupplieraddressRequest
	if err := decodeJSON(writer, request, &req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if err := validateAddress(req.Address); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	address := toDomainAddress(req.Address)
	if err := api.supplier.UpdateSupplierAddress(
		request.Context(), supplierID, address,
	); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete a supplier
// @Description Deletes the supplier identified by its UUID.
// @Tags Suppliers
// @Security BearerAuth
// @Param id path string true "Supplier UUID"
// @Success 204 "Deleted"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [delete]
func (api *API) deleteSupplier(writer http.ResponseWriter, request *http.Request) {
	supplierID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid supplier ID")
		return
	}
	if err := api.supplier.DeleteSupplier(request.Context(), supplierID); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary List suppliers
// @Description Returns all suppliers. An empty list is returned when there are no suppliers.
// @Tags Suppliers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} getSupplierResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [get]
func (api *API) getSuppliers(writer http.ResponseWriter, request *http.Request) {
	suppliers, err := api.supplier.GetSuppliers(request.Context())
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toGetSuppliersResponse(suppliers))
}

// @Summary Get a supplier by ID
// @Description Returns the supplier identified by its UUID.
// @Tags Suppliers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Supplier UUID"
// @Success 200 {object} getSupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [get]
func (api *API) getSupplierByID(writer http.ResponseWriter, request *http.Request) {
	supplierID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid supplier ID")
		return
	}
	supplier, err := api.supplier.GetSupplierByID(request.Context(), supplierID)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toGetSupplierResponse(supplier))
}
