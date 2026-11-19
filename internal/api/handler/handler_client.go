package httpapi

import (
	"net/http"
	"strconv"
	"uuid"
)

// @Summary Create a client
// @Description Creates a client and its address when the address does not already exist.
// @Tags Clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param client body CreateClientRequest true "Client data"
// @Success 201 {object} GetClientResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/clients [post]
func (api *API) createClient(writer http.ResponseWriter, request *http.Request) {
	var client_req CreateClientRequest
	if err := decodeJSON(writer, request, &client_req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if err := validateClient(client_req); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	client, err := api.client.CreateClient(request.Context(), toDomainClient(client_req))
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, toGetClientResponse(client))
}

// @Summary Delete a client
// @Description Deletes the client identified by its UUID.
// @Tags Clients
// @Security BearerAuth
// @Param id path string true "Client UUID"
// @Success 204 "Deleted"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/clients/{id} [delete]
func (api *API) deleteClient(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid client ID")
		return
	}
	if err := api.client.DeleteClient(request.Context(), id); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary List clients or search by name and surname
// @Description Returns all clients with optional pagination, or searches by both name and surname.
// @Tags Clients
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Maximum number of clients to return"
// @Param offset query int false "Number of clients to skip before returning results"
// @Param name query string false "Client first name; requires surname"
// @Param surname query string false "Client surname; requires name"
// @Success 200 {object} GetClientsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/clients [get]
func (api *API) getClients(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	var limit, offset uint
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 31)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = uint(parsed)
	}
	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 31)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = uint(parsed)
	}
	if name, surname := query.Get("name"), query.Get("surname"); name != "" || surname != "" {
		if name == "" || surname == "" {
			writeError(writer, http.StatusBadRequest, "both name and surname are required")
			return
		}
		clients, err := api.client.GetClientsByNameSurname(request.Context(), name, surname)
		if err != nil {
			writeOperationError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, toGetClientsResponse(clients))
		return
	}
	clients, err := api.client.GetClients(request.Context(), limit, offset)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toGetClientsResponse(clients))

}

// @Summary Change a client's address
// @Description Replaces the address associated with the specified client.
// @Tags Clients
// @Security BearerAuth
// @Accept json
// @Param id path string true "Client UUID"
// @Param address body UpdateClientAddressRequest true "New address"
// @Success 204 "Updated"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/clients/{id}/address [put]
func (api *API) updateClientAddress(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid client ID")
		return
	}
	var req UpdateClientAddressRequest
	if err := decodeJSON(writer, request, &req); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON request: "+err.Error())
		return
	}
	if err := validateAddress(req.Address); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	address := toDomainAddress(req.Address)
	if err := api.client.UpdateClientAddress(request.Context(), id, address); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)

}
