package createaccount

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Username     string  `json:"username" validate:"required"`
	Password     string  `json:"password" validate:"required"`
	ReferralCode *string `json:"referral_code,omitempty"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
}

type CreateUserResponseError struct {
	Error string `json:"error"`
}

type Handler struct {
	svc CreateUserService
}

func NewHandler(svc CreateUserService) *Handler {
	return &Handler{svc: svc}
}

// CreateUser godoc
//
//	@Summary		Create User
//	@Description	Create a new user
//	@Tags			accounts
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	CreateUserResponse	"Success"
//	@Failure		400	{object}	CreateUserResponseError	"Bad Request"
//	@Failure		500	{object}	CreateUserResponseError	"Internal Server Error"
//	@Router			/api/v1/accounts [post]
//	@Param			request	body		CreateUserRequest	true	"Create user request"
func (h *Handler) CreateUser(c echo.Context) error {
	var request CreateUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, CreateUserResponseError{
			Error: err.Error(),
		})
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, CreateUserResponseError{
			Error: err.Error(),
		})
	}

	createUserResponse, err := h.svc.CreateUser(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateUserResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, CreateUserResponse{
		Message: createUserResponse,
	})
}
