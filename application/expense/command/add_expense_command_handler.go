package commands

import (
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/models"
)

type AddExpenseCommandHandler struct {
	userRepository persistance.IUserRepository
}

func NewAddExpenseCommandHandler(userRepository persistance.IUserRepository) *AddExpenseCommandHandler {
	return &AddExpenseCommandHandler{userRepository: userRepository}
}

func (h *AddExpenseCommandHandler) Handle(command *AddExpenseCommand) (*models.Expense, error) {
	newExpense, err := fromAddExpenseCommand(command)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.GetUserById(command.UserId.String())
	if err != nil {
		return nil, err
	}

	err = user.AddExpense(newExpense)
	if err != nil {
		return nil, err
	}

	// TODO: save user and expense

	return newExpense, nil
}

func fromAddExpenseCommand(command *AddExpenseCommand) (*models.Expense, error) {
	return models.NewExpense(
		command.Description,
		command.Amount,
		command.UserId,
		command.Date,
	)
}
