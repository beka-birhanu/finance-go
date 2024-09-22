/*
Package errexpense defines expense-related errors for the application.

It provides a set of predefined errors related to expense not-found, and validation
issues. These errors are used throughout the application to handle variours error conditions
specfic to expense operations.
*/
package errexpense

import "github.com/beka-birhanu/finance-go/domain/error/common"

// Validation errors
var (
	// Amount is negative.
	NegativeAmount = errdmn.NewValidation("Expense.Amount cannot be negative or zero.")

	// Description is longer than allowed.
	DescriptionTooLong = errdmn.NewValidation("Expense.Description is too long.")

	// Description is empty.
	EmptyDescription = errdmn.NewValidation("Expense.Description cannot be empty.")
)

// NotFound errors
var (
	// Expense is does not exist.
	NotFound = errdmn.NewNotFound("Expense not found")
)
