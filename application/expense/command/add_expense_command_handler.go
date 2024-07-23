package command

import (
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type AddExpenseCommandHandler struct {
	userRepository repository.IUserRepository
}

func NewAddExpenseCommandHandler(userRepository repository.IUserRepository) *AddExpenseCommandHandler {
	return &AddExpenseCommandHandler{userRepository: userRepository}
}

func (h *AddExpenseCommandHandler) Handle(command *AddExpenseCommand) (*model.Expense, error) {
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

func fromAddExpenseCommand(command *AddExpenseCommand) (*model.Expense, error) {
	return model.NewExpense(
		command.Description,
		command.Amount,
		command.UserId,
		command.Date,
	)
}
