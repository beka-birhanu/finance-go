// Package expensecmd provides functionality for handling commands related to expenses.
package expensecmd

import (
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

// PatchHandler manages the patching of expenses.
type PatchHandler struct {
	expenseRepository irepository.IExpenseRepository // Repository for expense data
}

// NewPatchHandler creates a new PatchHandler with the provided expense repository.
func NewPatchHandler(expenseRepository irepository.IExpenseRepository) *PatchHandler {
	return &PatchHandler{
		expenseRepository: expenseRepository,
	}
}

// Handle processes a PatchCommand to update an existing expense.
// It retrieves the expense by its ID and user ID, updates the expense fields if provided,
// and saves the changes to the repository.
//
// Returns:
//   - *expensemodel.Expense: The updated expense.
//   - error: An error if the update fails. Possible errors include issues with retrieving
//     or updating the expense, or saving the changes to the repository.
func (h *PatchHandler) Handle(cmd *PatchCommand) (*expensemodel.Expense, error) {
	expense, err := h.expenseRepository.ById(cmd.Id, cmd.UserId)
	if err != nil {
		return nil, err
	}

	if cmd.Amount != nil {
		if err := expense.UpdateAmount(*cmd.Amount); err != nil {
			return nil, err
		}
	}
	if cmd.Description != nil {
		if err := expense.UpdateDescription(*cmd.Description); err != nil {
			return nil, err
		}
	}
	if cmd.Date != nil {
		expense.UpdateDate(*cmd.Date)
	}

	if err := h.expenseRepository.Save(expense); err != nil {
		return nil, err
	}

	// TODO: Publish an expense updated domain event
	return expense, nil
}
