package response

import "github.com/gofiber/fiber/v3"

// Envelope is the standard shape for every JSON response returned by the API.
type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Error represents the standardized error body. Details is populated for validation failures.
type Error struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details map[string][]string `json:"details,omitempty"`
}

// Success writes {"success": true, "message": ..., "data": ...} with the given HTTP status code.
func Success(c fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Envelope{Success: true, Message: message, Data: data})
}

// SuccessWithMeta writes {"success": true, "message": ..., "data": ..., "meta": ...} for paginated endpoints.
func SuccessWithMeta(c fiber.Ctx, statusCode int, message string, data any, meta *Meta) error {
	return c.Status(statusCode).JSON(Envelope{Success: true, Message: message, Data: data, Meta: meta})
}

// Fail writes {"success": false, "error": {"code": ..., "message": ...}}.
func Fail(c fiber.Ctx, statusCode int, code, message string) error {
	return c.Status(statusCode).JSON(Envelope{Success: false, Error: &Error{Code: code, Message: message}})
}

// FailWithDetails writes {"success": false, "error": {"code": ..., "message": ..., "details": ...}}.
func FailWithDetails(c fiber.Ctx, statusCode int, code, message string, details map[string][]string) error {
	return c.Status(statusCode).JSON(Envelope{Success: false, Error: &Error{Code: code, Message: message, Details: details}})
}
