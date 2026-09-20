package delivery

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	"github.com/bagusyanuar/pharmacy-be/internal/modules/user/usecase"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	uc *usecase.UserUsecase
}

// NewUserHandler registers user routes and returns the handler.
func NewUserHandler(router fiber.Router, uc *usecase.UserUsecase, authMiddleware fiber.Handler) *UserHandler {
	h := &UserHandler{uc: uc}

	users := router.Group("/users")
	if authMiddleware != nil {
		users.Use(authMiddleware)
	}

	users.Get("/", h.List)
	users.Get("/:id", h.GetByID)

	return h
}

// List returns a paginated list of users.
func (h *UserHandler) List(c fiber.Ctx) error {
	params := response.ParsePagination(c)
	users, total, err := h.uc.ListUsers(c.Context(), params)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch users")
	}

	meta := response.NewMeta(params.Page, params.PerPage, total)
	return response.SuccessWithMeta(c, fiber.StatusOK, "", toUserResponses(users), meta)
}

// GetByID returns user details by ID.
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return response.Fail(c, fiber.StatusNotFound, "USER_NOT_FOUND", "user not found")
		}
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch user")
	}

	return response.Success(c, fiber.StatusOK, "", toUserResponse(user))
}
