package repository

import (
	"database/sql"
	"errors"

	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	errexpense "github.com/beka-birhanu/finance-go/domain/error/expense"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

type ExpenseRepository struct {
	db *sql.DB
}

type primaryKey struct {
	Id     uuid.UUID
	UserId uuid.UUID
}

var Expenses = map[primaryKey]expensemodel.Expense{}

var _ irepository.IExpenseRepository = &ExpenseRepository{}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

// Add implements irepository.IExpenseRepository.
func (e *ExpenseRepository) Add(expense *expensemodel.Expense) error {
	key := primaryKey{
		Id:     expense.ID(),
		UserId: expense.UserID(),
	}

	if _, exists := Expenses[key]; exists {
		return errors.New("expense already exists")
	}

	Expenses[key] = *expense
	return nil
}

// ById implements irepository.IExpenseRepository.
func (e *ExpenseRepository) ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error) {
	key := primaryKey{
		Id:     id,
		UserId: userId,
	}

	expense, exists := Expenses[key]
	if !exists {
		return nil, errors.New("expense not found")
	}

	return &expense, nil
}

// Update implements irepository.IExpenseRepository.
func (e *ExpenseRepository) Update(expense *expensemodel.Expense) error {
	key := primaryKey{
		Id:     expense.ID(),
		UserId: expense.UserID(),
	}

	if _, exists := Expenses[key]; !exists {
		return errexpense.NotFound
	}

	Expenses[key] = *expense
	return nil
}

