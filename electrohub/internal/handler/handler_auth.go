package handler

import (
	"net/http"

	authpb "contracts/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary Register a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body registerRequest true "User registration data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (api *API) register(writer http.ResponseWriter, request *http.Request) {
	var body registerRequest

	if err := decodeJSON(writer, request, &body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON")
		return
	}
	response, err := api.auth.Register(request.Context(), &authpb.RegisterRequest{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Phone:     body.Phone,
		Password:  body.Password,
	})
	if err != nil {
		writeAuthError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, response)
}

// @Summary Authenticate a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body authRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth [post]
func (api *API) authenticate(writer http.ResponseWriter, request *http.Request) {
	var body authRequest
	if err := decodeJSON(writer, request, &body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON")
		return
	}
	response, err := api.auth.Authenticate(request.Context(), &authpb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		writeAuthError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, response)
}

// @Summary Recover a user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param account body resetRequest true "Account email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reset [post]
func (api *API) resetPassword(writer http.ResponseWriter, request *http.Request) {
	var body resetRequest
	if err := decodeJSON(writer, request, &body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON")
		return
	}
	response, err := api.auth.ResetPassword(request.Context(), &authpb.ResetPasswordRequest{
		Email: body.Email,
	})
	if err != nil {
		writeAuthError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, response)
}

// @Summary Change the authenticated user's password
// @Description Changes the password after verifying the current password. Requires a valid access token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body changePasswordRequest true "Current and new passwords"
// @Success 200 {object} changePasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /auth/change-password [put]
func (api *API) changePassword(writer http.ResponseWriter, request *http.Request) {
	var body changePasswordRequest
	if err := decodeJSON(writer, request, &body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON")
		return
	}

	accessToken, ok := request.Context().Value(accessTokenKey{}).(string)
	if !ok || accessToken == "" {
		writeError(writer, http.StatusUnauthorized, "missing or invalid authorization header")
		return
	}

	response, err := api.auth.ChangePassword(request.Context(), &authpb.ChangePasswordRequest{
		AccessToken: accessToken,
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	})
	if err != nil {
		writeAuthError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, changePasswordResponse{Message: response.GetMessage()})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		writeError(w, http.StatusBadRequest, status.Convert(err).Message())
	case codes.Unauthenticated:
		writeError(w, http.StatusUnauthorized, "invalid credentials or token")
	case codes.PermissionDenied:
		writeError(w, http.StatusForbidden, "permission denied")
	case codes.NotFound:
		writeError(w, http.StatusNotFound, "resource not found")
	case codes.AlreadyExists:
		writeError(w, http.StatusConflict, "email already exists")
	case codes.Unavailable:
		writeError(w, http.StatusServiceUnavailable, "authentication service unavailable")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
