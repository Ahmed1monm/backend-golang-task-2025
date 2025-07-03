package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ahmed1monm/backend-golang-task-2025/internal/api/dto"
	"github.com/Ahmed1monm/backend-golang-task-2025/internal/service"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/validator"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User registration details"
// @Success 201 {object} dto.UserProfileResponse
// @Failure 400 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request body
	req := new(dto.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errs := validator.Validate(req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": errs})
	}

	// Create user
	resp, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, resp)
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get a user's profile information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [get]
func (h *UserHandler) GetUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID from path parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Get user profile
	resp, err := h.userService.GetUserProfile(ctx, uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user profile")
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserProfileRequest true "Update profile request"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID from path parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Get authenticated user ID from context
	authUserID := c.Get("user_id").(uint)

	// Only allow users to update their own profile
	if authUserID != uint(id) {
		return echo.NewHTTPError(http.StatusForbidden, "You can only update your own profile")
	}

	// Parse request body
	req := new(dto.UpdateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errs := validator.Validate(req); len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": errs})
	}

	// Update user profile
	resp, err := h.userService.UpdateUserProfile(ctx, uint(id), req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		case "email is already taken":
			return echo.NewHTTPError(http.StatusConflict, "Email is already taken")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user profile")
		}
	}

	return c.JSON(http.StatusOK, resp)
}
