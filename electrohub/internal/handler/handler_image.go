package handler

import (
	"electrohub/internal/domain"
	"net/http"

	"uuid"
)

// @Summary Add an image to a product
// @Description Uploads an image file and associates it with the specified product.
// @Tags Images
// @Security BearerAuth
// @Accept mpfd
// @Param product_id formData string true "UUID of the product that will receive the image"
// @Param image formData file true "Image file to upload"
// @Success 201 "Image added"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Failure 415 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/images [post]
func (api *API) addImage(writer http.ResponseWriter, request *http.Request) {
	data, err := readImage(writer, request)
	if err != nil {
		return
	}
	productID, err := uuid.Parse(request.FormValue("product_id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	image := domain.Image{Data: data}
	if err := api.image.AddImage(request.Context(), image, productID); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

// @Summary Replace an image
// @Description Replaces the image data for the specified image UUID.
// @Tags Images
// @Security BearerAuth
// @Accept mpfd
// @Param id path string true "Image UUID"
// @Param image formData file true "Replacement image file"
// @Success 204 "Image replaced"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Failure 415 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/images/{id} [put]
func (api *API) updateImage(writer http.ResponseWriter, request *http.Request) {
	imageID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid image ID")
		return
	}
	data, err := readImage(writer, request)
	if err != nil {
		return
	}
	image := domain.Image{
		ID:   imageID,
		Data: data,
	}

	if err := api.image.UpdateImage(request.Context(), image); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete an image
// @Description Deletes the image and removes its association from any product.
// @Tags Images
// @Security BearerAuth
// @Param id path string true "Image UUID"
// @Success 204 "Deleted"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/images/{id} [delete]
func (api *API) deleteImage(writer http.ResponseWriter, request *http.Request) {
	imageID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid image ID")
		return
	}
	if err := api.image.DeleteImage(request.Context(), imageID); err != nil {
		writeOperationError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// @Summary Download an image by ID
// @Description Returns the image identified by its UUID as binary data.
// @Tags Images
// @Security BearerAuth
// @Produce octet-stream
// @Param id path string true "Image UUID"
// @Success 200 {file} file
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/images/{id} [get]
func (api *API) getImageByID(writer http.ResponseWriter, request *http.Request) {
	imageID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid image ID")
		return
	}
	image, err := api.image.GetImageByID(request.Context(), imageID)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeImage(writer, image.Data)
}

// @Summary Download a product image
// @Description Returns the image associated with the specified product as binary data.
// @Tags Images
// @Security BearerAuth
// @Produce octet-stream
// @Param id path string true "Product UUID"
// @Success 200 {file} file
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id}/image [get]
func (api *API) getImageByProductID(writer http.ResponseWriter, request *http.Request) {
	productID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid product ID")
		return
	}
	image, err := api.image.GetImageByProductID(request.Context(), productID)
	if err != nil {
		writeOperationError(writer, err)
		return
	}
	writeImage(writer, image.Data)
}
