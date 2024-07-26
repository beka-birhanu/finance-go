package command

import (
	"time"

	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type AddHandler struct {
	userRepository irepository.IUserRepository
	timeService    itimeservice.IService
}

var _ icmd.IHandler[*AddExpenseCommand, *expensemodel.Expense] = &AddHandler{}

func New(userRepository irepository.IUserRepository, timeService itimeservice.IService) *AddHandler {
	return &AddHandler{userRepository: userRepository, timeService: timeService}
}

func (h *AddHandler) Handle(command *AddExpenseCommand) (*expensemodel.Expense, error) {
	newExpense, err := newExpense(command, h.timeService.NowUTC())
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.ById(command.UserId)
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
