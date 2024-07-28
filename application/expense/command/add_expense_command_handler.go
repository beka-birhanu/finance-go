package expensecmd

import (
	"fmt"
	"time"

	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type AddHandler struct {
	expenseRepository irepository.IExpenseRepository
	userRepository    irepository.IUserRepository
	timeService       itimeservice.IService
}

var _ icmd.IHandler[*AddCommand, *expensemodel.Expense] = &AddHandler{}

type Config struct {
	UserRepository    irepository.IUserRepository
	TimeService       itimeservice.IService
	ExpenseRepository irepository.IExpenseRepository
}

func NewAddHandler(config Config) *AddHandler {
	return &AddHandler{
		userRepository:    config.UserRepository,
		timeService:       config.TimeService,
		expenseRepository: config.ExpenseRepository,
	}
}

func (h *AddHandler) Handle(command *AddCommand) (*expensemodel.Expense, error) {
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

	// TODO: make this in one transaction
	if err := h.userRepository.Save(user); err != nil {
		return nil, fmt.Errorf("unable to update user: %w", err)
	}

	if err := h.expenseRepository.Save(newExpense); err != nil {
		return nil, fmt.Errorf("unable to save expense: %w", err)
	}

	return newExpense, nil
}

func newExpense(command *AddCommand, currentUTCTime time.Time) (*expensemodel.Expense, error) {
	config := expensemodel.Config{
		Description:  command.Description,
		Amount:       command.Amount,
		UserId:       command.UserId,
		Date:         command.Date,
		CreationTime: currentUTCTime,
	}
	return expensemodel.New(config)
}
