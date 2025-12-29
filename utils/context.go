package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/labstack/echo"
)

// CustomApplicationContextStub is a stub to replace util.CustomApplicationContext
type CustomApplicationContextStub struct {
	echo.Context
}

// CustomResponse sends a JSON response
func (c *CustomApplicationContextStub) CustomResponse(status string, data interface{}, message string, httpStatus int, errorCode string, meta interface{}) error {
	response := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	} else {
		response["data"] = []interface{}{}
	}

	if errorCode != "" {
		response["error_code"] = errorCode
	}

	if meta != nil {
		response["meta"] = meta
	}

	return c.JSON(httpStatus, response)
}

// Bind binds the request body to a struct
func (c *CustomApplicationContextStub) Bind(i interface{}) error {
	return c.Context.Bind(i)
}

// Validate validates the struct using reflection
func (c *CustomApplicationContextStub) Validate(i interface{}) error {
	if i == nil {
		return fmt.Errorf("cannot validate nil interface")
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil // Not a struct, nothing to validate
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Check for required tag
		if required := field.Tag.Get("validate"); required == "required" {
			if fieldValue.Kind() == reflect.String {
				if fieldValue.String() == "" {
					return fmt.Errorf("field %s is required", field.Name)
				}
			} else if fieldValue.Kind() == reflect.Int || fieldValue.Kind() == reflect.Int64 {
				if fieldValue.Int() == 0 {
					return fmt.Errorf("field %s is required", field.Name)
				}
			} else if fieldValue.IsNil() {
				return fmt.Errorf("field %s is required", field.Name)
			}
		}
	}

	return nil
}

// GetCustomApplicationContextStub safely converts echo.Context to CustomApplicationContextStub
func GetCustomApplicationContextStub(c echo.Context) *CustomApplicationContextStub {
	return &CustomApplicationContextStub{
		Context: c,
	}
}

// GetUserNameFromRequest extracts user name from echo.Context request
func GetUserNameFromRequest(c echo.Context) string {
	// Try to get from X-User-Name header
	if userName := c.Request().Header.Get("X-User-Name"); userName != "" {
		return userName
	}

	// Try to get from X-User-Email header
	if userEmail := c.Request().Header.Get("X-User-Email"); userEmail != "" {
		return userEmail
	}

	// Try to parse from cookie "u" (user data JSON)
	cookie, err := c.Cookie("u")
	if err == nil && cookie != nil && cookie.Value != "" {
		if cookie.Value != "undefined" && cookie.Value != "null" && cookie.Value != "" {
			var userData struct {
				ID    interface{} `json:"id"`
				Name  string      `json:"name"`
				Email string      `json:"email"`
			}
			if err := json.Unmarshal([]byte(cookie.Value), &userData); err == nil {
				if userData.Name != "" {
					return userData.Name
				}
				if userData.Email != "" {
					return userData.Email
				}
			}
		}
	}

	// Default to Anonymous
	return "Anonymous"
}
