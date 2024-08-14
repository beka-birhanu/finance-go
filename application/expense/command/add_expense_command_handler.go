// Package expensecmd provides functionality for handling commands related to expenses.
package expensecmd

import (
	"fmt"
	"time"

	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

// AddHandler handles commands for adding new expenses.
type AddHandler struct {
	userRepo irepository.IUserRepository // Repository for user data
	timeSvc  itimeservice.IService       // Service for time-related operations
}

// Ensure AddHandler implements icmd.IHandler[*AddCommand, *expensemodel.Expense].
var _ icmd.IHandler[*AddCommand, *expensemodel.Expense] = &AddHandler{}

// Config holds dependencies required for creating an AddHandler.
type Config struct {
	UserRepository irepository.IUserRepository // Repository for user data
	TimeService    itimeservice.IService       // Service for time-related operations
}

// NewAddHandler creates a new AddHandler with the specified configuration.
func NewAddHandler(config Config) *AddHandler {
	return &AddHandler{
		userRepo: config.UserRepository,
		timeSvc:  config.TimeService,
	}
}

// Handle processes an AddCommand to create a new expense and returns the expense.
func (h *AddHandler) Handle(command *AddCommand) (*expensemodel.Expense, error) {
	newExpense, err := createExpense(command, h.timeSvc.NowUTC())
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.ById(command.UserId)
	if err != nil {
		return nil, err
	}

	if err := user.AddExpense(newExpense, h.timeSvc.NowUTC()); err != nil {
		return nil, err
	}

	if err := h.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("unable to update user: %w", err)
	}

	return newExpense, nil
}

// createExpense constructs an Expense instance using the command and current time.
func createExpense(command *AddCommand, currentTime time.Time) (*expensemodel.Expense, error) {
	config := expensemodel.Config{
		Description:  command.Description,
		Amount:       command.Amount,
		UserId:       command.UserId,
		Date:         command.Date,
		CreationTime: currentTime,
	}
	return expensemodel.New(config)
}
