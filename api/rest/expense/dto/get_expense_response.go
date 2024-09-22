package dto

import (
	"time"

	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

type GetExpenseResponse struct {
	Id          uuid.UUID `json:"id"`
	Amount      float32   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}

func FromExpenseModel(expense *expensemodel.Expense) *GetExpenseResponse {
	return &GetExpenseResponse{
		Id:          expense.ID(),
		Amount:      expense.Amount(),
		Description: expense.Description(),
		Date:        expense.Date(),
		CreatedAt:   expense.CreatedAt(),
	}
}
