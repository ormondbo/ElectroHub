package handler

import (
	"electrohub/internal/api/repository"
	"electrohub/internal/api/service"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ======================================================================================================
// Decode/Encode
// ======================================================================================================

func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) error {
	request.Body = http.MaxBytesReader(writer, request.Body, 10<<20) //Max 20 MB

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("request must contain exactly one JSON value")
	}

	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func writeImage(writer http.ResponseWriter, data []byte) {
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", `attachment; filename="image"`)

	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(data); err != nil {
		log.Printf("write image response: %v", err)
	}
}

func readImage(writer http.ResponseWriter, request *http.Request) ([]byte, error) {
	request.Body = http.MaxBytesReader(writer, request.Body, 10<<20)

	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil && request.Header.Get("Content-Type") != "" {
		writeError(writer, http.StatusBadRequest, "invalid image Content-Type")
		return nil, err
	}
	/*
		The technical specifications required checking only for the "application/octet-stream" type;
		however, to support an older version of OpenAPI, the "multipart/form-data" type was also defined.
		When accessing the API directly, the code associated with this tag is either not utilized
		or represents an obsolete method of access. Please do not use "multipart/form-data".

		OPEN_TAG: old_openapi
	*/
	if mediaType == "multipart/form-data" {
		if err := request.ParseMultipartForm(10 << 20); err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				writeError(writer, http.StatusRequestEntityTooLarge, "image exceeds the 10 MiB limit")
				return nil, err
			}

			writeError(writer, http.StatusBadRequest, "invalid multipart image request")
			return nil, err
		}

		if request.MultipartForm != nil {
			defer request.MultipartForm.RemoveAll()
		}

		file, _, err := request.FormFile("image")
		if err != nil {
			writeError(writer, http.StatusBadRequest, "image file is required")
			return nil, err
		}

		data, readErr := io.ReadAll(io.LimitReader(file, (10<<20)+1))
		closeErr := file.Close()

		if readErr != nil {
			writeError(writer, http.StatusBadRequest, "failed to read image")
			return nil, readErr
		}

		if closeErr != nil {
			writeError(writer, http.StatusBadRequest, "failed to read image")
			return nil, closeErr
		}

		if len(data) > 10<<20 {
			writeError(writer, http.StatusRequestEntityTooLarge, "image exceeds the 10 MiB limit")
			return nil, errors.New("image exceeds the size limit")
		}

		if len(data) == 0 {
			writeError(writer, http.StatusBadRequest, "image cannot be empty")
			return nil, errors.New("empty image")
		}

		return data, nil
	}
	// CLOSE_TAG: old_openapi

	if mediaType != "" && mediaType != "application/octet-stream" {
		writeError(writer, http.StatusUnsupportedMediaType, "image must use application/octet-stream")
		return nil, errors.New("unsupported image Content-Type")
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(writer, http.StatusRequestEntityTooLarge, "image exceeds the 10 MiB limit")
			return nil, err
		}

		writeError(writer, http.StatusBadRequest, "failed to read image")
		return nil, err
	}

	if len(data) == 0 {
		writeError(writer, http.StatusBadRequest, "image cannot be empty")
		return nil, errors.New("empty image")
	}

	return data, nil
}

// ======================================================================================================
// Validators
// ======================================================================================================

func validateAddress(address addressRequest) error {
	if strings.TrimSpace(address.Country) == "" {
		return errors.New("country is required")
	}
	if strings.TrimSpace(address.City) == "" {
		return errors.New("city is required")
	}
	if strings.TrimSpace(address.Street) == "" {
		return errors.New("street is required")
	}
	if utf8.RuneCountInString(address.Country) > 100 ||
		utf8.RuneCountInString(address.City) > 100 ||
		utf8.RuneCountInString(address.Street) > 255 {
		return errors.New("address field is too long")
	}
	return nil
}

func validateClient(request createClientRequest) error {
	if strings.TrimSpace(request.Name) == "" || utf8.RuneCountInString(request.Name) > 100 {
		return errors.New("name is required and must be at most 100 characters")
	}
	if strings.TrimSpace(request.Surname) == "" || utf8.RuneCountInString(request.Surname) > 100 {
		return errors.New("surname is required and must be at most 100 characters")
	}
	if request.Birthday.IsZero() || request.Birthday.After(time.Now().UTC()) {
		return errors.New("birthday must be a valid date that is not in the future")
	}
	switch request.Gender {
	case "male", "female", "other":
	default:
		return errors.New("gender must be male, female or other")
	}
	return validateAddress(request.Address)
}

func validateProduct(request createProductRequest) error {
	if strings.TrimSpace(request.Name) == "" || utf8.RuneCountInString(request.Name) > 255 {
		return errors.New("name is required and must be at most 255 characters")
	}
	if strings.TrimSpace(request.Category) == "" || utf8.RuneCountInString(request.Category) > 100 {
		return errors.New("category is required and must be at most 100 characters")
	}
	if math.IsNaN(request.Price) || math.IsInf(request.Price, 0) || request.Price < 0 || request.Price >= 10_000_000_000 {
		return errors.New("price must be a finite non-negative number below 10000000000")
	}
	if request.AvailableStock > math.MaxInt32 {
		return errors.New("available_stock is too large")
	}
	if request.SupplierID == uuid.Nil() {
		return errors.New("supplier_id is required")
	}
	return nil
}

func validateSupplier(request createSupplierRequest) error {
	if strings.TrimSpace(request.Name) == "" || utf8.RuneCountInString(request.Name) > 255 {
		return errors.New("name is required and must be at most 255 characters")
	}
	if strings.TrimSpace(request.PhoneNumber) == "" || utf8.RuneCountInString(request.PhoneNumber) > 20 {
		return errors.New("phone_number is required and must be at most 20 characters")
	}
	return validateAddress(request.Address)
}

// ======================================================================================================
// ErrorWritters
// ======================================================================================================

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}

func writeOperationError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound), errors.Is(err, pgx.ErrNoRows):
		writeError(writer, http.StatusNotFound, "resource not found")

	case errors.Is(err, service.ErrInsufficientStock):
		writeError(writer, http.StatusConflict, err.Error())

	case errors.Is(err, service.ErrProductAlreadyHasImage):
		writeError(writer, http.StatusConflict, err.Error())

	default:
		var databaseError *pgconn.PgError

		if errors.As(err, &databaseError) {
			switch databaseError.Code {
			case "23503", "23505":
				writeError(writer, http.StatusConflict, "resource conflicts with existing data")
				return

			case "23514", "22001", "22P02":
				writeError(writer, http.StatusBadRequest, "request data violates a database constraint")
				return
			}
		}

		log.Printf("request failed: %v", err)
		writeError(writer, http.StatusInternalServerError, "internal server error")
	}
}
