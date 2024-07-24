package command

import (
	"time"

	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type AddExpenseCommandHandler struct {
	userRepository irepository.IUserRepository
	timeService    itimeservice.ITimeService
}

var _ icmd.ICommandHandler[*AddExpenseCommand, *expensemodel.Expense] = &AddExpenseCommandHandler{}

func New(userRepository irepository.IUserRepository, timeService itimeservice.ITimeService) *AddExpenseCommandHandler {
	return &AddExpenseCommandHandler{userRepository: userRepository, timeService: timeService}
}

func (h *AddExpenseCommandHandler) Handle(command *AddExpenseCommand) (*expensemodel.Expense, error) {
	newExpense, err := newExpense(command, h.timeService.NowUTC())
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

func newExpense(command *AddExpenseCommand, currentUTCTime time.Time) (*expensemodel.Expense, error) {
	config := expensemodel.Config{
		Description:  command.Description,
		Amount:       command.Amount,
		UserId:       command.UserId,
		Date:         command.Date,
		CreationTime: currentUTCTime,
	}
	return expensemodel.New(config)
}
