package expensecmd

import (
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

// PatchHandler handles patching an expense.
type PatchHandler struct {
	expenseRepository irepository.IExpenseRepository
}

// NewPatchHandler creates a new PatchHandler.
func NewPatchHandler(expenseRepository irepository.IExpenseRepository) *PatchHandler {
	return &PatchHandler{
		expenseRepository: expenseRepository,
	}
}

// Handle handles the patching of an expense.
func (h *PatchHandler) Handle(cmd *PatchCommand) (*expensemodel.Expense, error) {
	// Retrieve the expense by ID and User ID
	expense, err := h.expenseRepository.ById(cmd.Id, cmd.UserId)
	if err != nil {
		return nil, err
	}

	// Update the expense fields if provided
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

	err = h.expenseRepository.Save(expense)
	if err != nil {
		return nil, err
	}

	// TODO: publish an expense udpated domain event
	return expense, nil
}
