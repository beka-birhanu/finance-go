package persistance

import "github.com/beka-birhanu/finance-go/domain/entities"

type IUserRepository interface {
	CreateExpense(expense *entities.User) error
	GetExpense(id string) (*entities.User, error)
	ListExpenses() ([]*entities.User, error)
}
