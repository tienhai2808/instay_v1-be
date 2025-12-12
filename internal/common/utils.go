package common

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mr-tron/base58"
)

func HandleValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is require", strings.ToLower(e.Field()))
			case "email":
				return fmt.Sprintf("%s is not a valid email", strings.ToLower(e.Field()))
			case "min":
				return fmt.Sprintf("%s must be at least %s characters long", strings.ToLower(e.Field()), e.Param())
			case "max":
				return fmt.Sprintf("%s cannot exceed %s characters", strings.ToLower(e.Field()), e.Param())
			case "len":
				return fmt.Sprintf("%s must have exactly %s characters", strings.ToLower(e.Field()), e.Param())
			case "numeric":
				return fmt.Sprintf("%s must be a number", strings.ToLower(e.Field()))
			case "uuid4":
				return fmt.Sprintf("%s must be a valid UUID v4", strings.ToLower(e.Field()))
			case "oneof":
				return fmt.Sprintf("%s must have the value: %s", strings.ToLower(e.Field()), e.Param())
			default:
				return fmt.Sprintf("%s is not valid", strings.ToLower(e.Field()))
			}
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return fmt.Sprintf("%s must be a %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return fmt.Sprintf("invalid JSON at byte %d", syntaxError.Offset)
	}

	if err != nil {
		return err.Error()
	}

	return "invalid data"
}

func ToAPIResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, types.APIResponse{
		Message: message,
		Data:    data,
	})
}

func IsUniqueViolation(err error) (bool, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return true, pgErr.ConstraintName
		}
	}

	return false, ""
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

func GenerateSlug(str string) string {
	return slug.Make(str)
}

func GenerateBase58ID(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return base58.Encode(b)
}
