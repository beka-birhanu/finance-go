package expensecmd

import (
	"fmt"
	"time"

	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

// AddHandler handles commands to add new expenses.
type AddHandler struct {
	userRepo irepository.IUserRepository
	timeSvc  itimeservice.IService
}

// Ensure AddHandler implements icmd.IHandler[*AddCommand, *expensemodel.Expense].
var _ icmd.IHandler[*AddCommand, *expensemodel.Expense] = &AddHandler{}

// Config holds the configuration for creating a new AddHandler.
type Config struct {
	UserRepository irepository.IUserRepository
	TimeService    itimeservice.IService
}

// NewAddHandler creates a new instance of AddHandler with the provided configuration.
func NewAddHandler(config Config) *AddHandler {
	return &AddHandler{
		userRepo: config.UserRepository,
		timeSvc:  config.TimeService,
	}
}

// Handle processes the AddCommand to add a new expense and returns the created expense.
func (h *AddHandler) Handle(command *AddCommand) (*expensemodel.Expense, error) {
	newExpense, err := createExpense(command, h.timeSvc.NowUTC())
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.ById(command.UserId)
	if err != nil {
		return nil, err
	}

	err = user.AddExpense(newExpense, h.timeSvc.NowUTC())
	if err != nil {
		return nil, err
	}

	// TODO: make this operation transactional
	if err := h.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("unable to update user: %w", err)
	}

	return newExpense, nil
}

// createExpense constructs a new Expense instance based on the provided command and current time.
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

