package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func SendSuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, meta ...interface{}) error {
	res := APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	if len(meta) > 0 && meta[0] != nil {
		res.Meta = meta[0]
	}

	return c.Status(statusCode).JSON(res)
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	res := APIResponse{
		Status:  statusCode,
		Message: message,
	}
	return c.Status(statusCode).JSON(res)
}