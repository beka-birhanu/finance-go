package expense

import (
	"github.com/beka-birhanu/finance-go/application/expense/commands"
	"github.com/beka-birhanu/finance-go/domain/models"
)

type IAddExpenseCommand interface {
	Handle(command *commands.AddExpenseCommand) (*models.Expense, error)
}
