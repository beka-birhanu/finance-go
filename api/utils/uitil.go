package utils

import (
	"context"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func ConfirmUserID(ctx context.Context, userId uuid.UUID) error {
	claims, ok := ctx.Value(middleware.ContextUserClaims).(jwt.MapClaims)
	if !ok {
		return errapi.NewServerError("Error on retrieving user id from context")
	}

	userIDStr, ok := claims["user_id"].(string)
	if userIDStr == "" {
		return errapi.NewAuthentication("User claims not found!")

	} else if !ok || userId.String() != userIDStr {
		return errapi.NewForbidden("The response does not belong to the user requesting.")
	}

	return nil
}
