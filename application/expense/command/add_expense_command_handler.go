package command

import (
	"time"

	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type AddExpenseCommandHandler struct {
	userRepository repository.IUserRepository
	timeService    timeservice.ITimeService
}

func NewAddExpenseCommandHandler(userRepository repository.IUserRepository, timeService timeservice.ITimeService) *AddExpenseCommandHandler {
	return &AddExpenseCommandHandler{userRepository: userRepository, timeService: timeService}
}

func (h *AddExpenseCommandHandler) Handle(command *AddExpenseCommand) (*model.Expense, error) {
	newExpense, err := fromAddExpenseCommand(command, h.timeService.NowUTC())
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.GetUserById(command.UserId)
	if err != nil {
		return nil, err
	}

	err = user.AddExpense(newExpense, h.timeService.NowUTC())
	if err != nil {
		return nil, err
	}

	// TODO: save user and expense

	return newExpense, nil
}

func fromAddExpenseCommand(command *AddExpenseCommand, currentUTCTime time.Time) (*model.Expense, error) {
	return model.NewExpense(
		command.Description,
		command.Amount,
		command.UserId,
		command.Date,
		currentUTCTime,
	)
}
