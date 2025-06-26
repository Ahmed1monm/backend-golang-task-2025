package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-golang-task-2025/pkg/errors"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	return errors.NewValidationError("Invalid user data", map[string]string{
		"email": "email is required",
		"password": "password must be at least 8 characters",
	}, http.StatusBadRequest)
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	return errors.NewBusinessError("User not found", errors.ErrCodeResourceNotFound, http.StatusNotFound)
}

func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	return errors.NewBusinessError("Unauthorized access", errors.ErrCodeUnauthorized, http.StatusUnauthorized)
}
