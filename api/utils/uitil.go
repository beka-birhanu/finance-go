package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/middleware"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
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

// constructQueryParams constructs the query parameters for retrieving multiple expenses,
// based on the user ID, cursor, limit, sort field, and sort order.
func ConstructQueryParams(userId uuid.UUID, cursor string, limit int, sortField string, sortOrder string) (*expensqry.GetMultipleQuery, error) {
	var lastSeenID uuid.UUID
	var lastSeenDate time.Time
	var lastSeenAmt float64
	var ascending bool

	if cursor != "" {
		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format1")
		}

		cursor = string(cursorByte)
		cursorParts := strings.Split(cursor, ",")
		if len(cursorParts) != 2 {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format1")
		}
		lastSeenID, err = uuid.Parse(cursorParts[0])
		if err != nil {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format2")
		}

		if sortField == "date" {
			lastSeenDate, err = time.Parse(time.RFC3339Nano, cursorParts[1])
			if err != nil {
				return &expensqry.GetMultipleQuery{}, fmt.Errorf("invalid cursor format for createdAt: %v", err)
			}
		} else if sortField == "amount" {
			lastSeenAmt, err = strconv.ParseFloat(cursorParts[1], 32)
			if err != nil {
				return &expensqry.GetMultipleQuery{}, fmt.Errorf("invalid cursor format for amount")
			}
		}
	}

	if sortOrder == "asc" {
		ascending = true
	}

	log.Println(userId, limit, sortField, &lastSeenID, &lastSeenDate, lastSeenAmt, ascending)
	return &expensqry.GetMultipleQuery{
		UserID:       userId,
		Limit:        limit,
		By:           sortField,
		LastSeenID:   &lastSeenID,
		LastSeenDate: &lastSeenDate,
		LastSeenAmt:  lastSeenAmt,
		Ascending:    ascending,
	}, nil
}

// BuildCursor constructs a cursor string for pagination, based on the last expense and the sort field.
func BuildCursor(lastExpense *expensemodel.Expense, field string) string {
	nextCursor := ""
	if lastExpense != nil {
		if field == "amount" {
			nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%f", lastExpense.ID(), lastExpense.Amount())))
		} else {
			date := lastExpense.Date().Format(time.RFC3339Nano)
			nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%v", lastExpense.ID(), date)))
		}
	}

	return nextCursor
}
