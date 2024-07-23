package error

import "errors"

var (
	ErrNegativeExpenseAmount     = errors.New("Expense.Amount cannot be negative")
	ErrExpenseDescriptionTooLong = errors.New("Expense.Description is too long")
	ErrEmptyExpenseDescription   = errors.New("Expense.Description cannot be empty")
)
