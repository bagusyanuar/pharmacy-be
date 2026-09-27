package delivery

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"

	authdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/domain"
	"github.com/bagusyanuar/pharmacy-be/internal/modules/auth/usecase"
	userdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	"github.com/bagusyanuar/pharmacy-be/pkg/jwt"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	uc                     *usecase.AuthUsecase
	refreshTokenCookieName string
	cookieSecure           bool
	refreshExpiration      time.Duration
}

// NewAuthHandler mounts auth routes and returns the handler.
func NewAuthHandler(
	router fiber.Router,
	uc *usecase.AuthUsecase,
	authMiddleware fiber.Handler,
	refreshTokenCookieName string,
	cookieSecure bool,
	refreshExpiration time.Duration,
) *AuthHandler {
	h := &AuthHandler{
		uc:                     uc,
		refreshTokenCookieName: refreshTokenCookieName,
		cookieSecure:           cookieSecure,
		refreshExpiration:      refreshExpiration,
	}

	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)

	if authMiddleware != nil {
		auth.Get("/me", authMiddleware, h.Me)
	}

	return h
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user account with role binding
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Register Request Payload"
// @Success      201      {object}  response.Response{data=ProfileResponse}
// @Failure      400      {object}  response.Response  "Invalid request body format or validation failure"
// @Failure      409      {object}  response.Response  "Email already registered"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "INVALID_BODY", "cannot parse request body")
	}

	if ok, err := response.ValidateOrFail(c, req); !ok {
		return err
	}

	user, err := h.uc.Register(c.Context(), req.Email, req.RoleID, req.Password)
	if err != nil {
		if errors.Is(err, userdomain.ErrEmailAlreadyExists) {
			return response.Fail(c, fiber.StatusConflict, "EMAIL_ALREADY_EXISTS", err.Error())
		}
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user")
	}

	return response.Success(c, fiber.StatusCreated, "user registered successfully", toProfileResponse(user))
}

// Login godoc
// @Summary      Web ERP Login
// @Description  Authenticate using email and password to receive an Access Token and HttpOnly Refresh Cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login Request Payload"
// @Success      200      {object}  response.Response{data=TokenResponse}
// @Failure      400      {object}  response.Response  "Invalid request body or validation failure"
// @Failure      401      {object}  response.Response  "Invalid email or password"
// @Failure      403      {object}  response.Response  "Account is inactive"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "INVALID_BODY", "cannot parse request body")
	}

	if ok, err := response.ValidateOrFail(c, req); !ok {
		return err
	}

	tokens, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, authdomain.ErrInvalidCredentials) {
			return response.Fail(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		}
		if errors.Is(err, authdomain.ErrAccountInactive) {
			return response.Fail(c, fiber.StatusForbidden, "ACCOUNT_INACTIVE", err.Error())
		}
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to process login")
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return response.Success(c, fiber.StatusOK, "login successful", TokenResponse{
		AccessToken: tokens.AccessToken,
		TokenType:   "Bearer",
	})
}

// Refresh godoc
// @Summary      Refresh Access Token
// @Description  Exchange the HttpOnly refresh token cookie for a new access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200      {object}  response.Response{data=TokenResponse}
// @Failure      401      {object}  response.Response  "Invalid or expired refresh token"
// @Failure      403      {object}  response.Response  "Account is inactive"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies(h.refreshTokenCookieName)
	if refreshToken == "" {
		return response.Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "missing refresh token cookie")
	}

	tokens, err := h.uc.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, authdomain.ErrInvalidRefreshToken) {
			return response.Fail(c, fiber.StatusUnauthorized, "INVALID_REFRESH_TOKEN", err.Error())
		}
		if errors.Is(err, authdomain.ErrAccountInactive) {
			return response.Fail(c, fiber.StatusForbidden, "ACCOUNT_INACTIVE", err.Error())
		}
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to refresh token")
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return response.Success(c, fiber.StatusOK, "token refreshed successfully", TokenResponse{
		AccessToken: tokens.AccessToken,
		TokenType:   "Bearer",
	})
}

// Logout godoc
// @Summary      Logout Session
// @Description  Clear the active user session and purge the refresh token cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200      {object}  response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})
	return response.Success(c, fiber.StatusOK, "logged out successfully", nil)
}

// Me godoc
// @Summary      Get Authenticated User Profile
// @Description  Retrieve profile details for the currently authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.Response{data=ProfileResponse}
// @Failure      401      {object}  response.Response  "Unauthorized"
// @Failure      404      {object}  response.Response  "User not found"
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c fiber.Ctx) error {
	userID, ok := c.Locals(string(jwt.Claims{}.UserID)).(string)
	if !ok || userID == "" {
		// Fallback to "user_id" string
		if uid, isStr := c.Locals("user_id").(string); isStr {
			userID = uid
		}
	}

	if userID == "" {
		return response.Fail(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}

	user, err := h.uc.Me(c.Context(), userID)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			return response.Fail(c, fiber.StatusNotFound, "USER_NOT_FOUND", "user not found")
		}
		return response.Fail(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "failed to retrieve profile")
	}

	return response.Success(c, fiber.StatusOK, "", toProfileResponse(user))
}

func (h *AuthHandler) setRefreshTokenCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(h.refreshExpiration),
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})
}
