package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/graph/model"
	"github.com/beka-birhanu/finance-go/api/graph/utils"
	generalUtil "github.com/beka-birhanu/finance-go/api/utils"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	ierr "github.com/beka-birhanu/finance-go/domain/common/error"
	"github.com/google/uuid"
)

// CreateExpense is the resolver for the createExpense field.
func (r *mutationResolver) CreateExpense(ctx context.Context, data model.CreateExpenseInput) (*model.Expense, error) {
	if err := generalUtil.ConfirmUserID(ctx, data.UserID); err != nil {
		return nil, utils.NewGQLError(err.(errapi.Error))
	}

	expense, err := r.addExpenseHandler.Handle(&expensecmd.AddCommand{
		UserId:      data.UserID,
		Date:        data.Date,
		Description: data.Description,
		Amount:      data.Amount,
	})
	if err != nil {
		return nil, utils.NewGQLError(errapi.Map(err.(ierr.IErr)))
	}

	return utils.NewExpense(expense), nil
}

// UpdateExpense is the resolver for the updateExpense field.
func (r *mutationResolver) UpdateExpense(ctx context.Context, data model.UpdateExpenseInput) (*model.Expense, error) {
	if err := generalUtil.ConfirmUserID(ctx, data.UserID); err != nil {
		return nil, utils.NewGQLError(err.(errapi.Error))
	}

	expense, err := r.patchExpenseHandler.Handle(&expensecmd.PatchCommand{
		Id:          data.ID,
		UserId:      data.UserID,
		Date:        data.Date,
		Description: data.Description,
		Amount:      data.Amount,
	})
	if err != nil {
		return nil, utils.NewGQLError(errapi.Map(err.(ierr.IErr)))
	}

	return utils.NewExpense(expense), nil
}

// Expense is the resolver for the expense field.
func (r *queryResolver) Expense(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Expense, error) {
	if err := generalUtil.ConfirmUserID(ctx, userID); err != nil {
		return nil, utils.NewGQLError(err.(errapi.Error))
	}
	expense, err := r.getExpenseHandler.Handle(&expensqry.GetQuery{UserId: userID, ExpenseId: id})
	if err != nil {
		return nil, utils.NewGQLError(errapi.Map(err.(ierr.IErr)))
	}
	return utils.NewExpense(expense), nil
}

// Expenses is the resolver for the expenses field.
func (r *queryResolver) Expenses(ctx context.Context, params model.GetMultipleInput) (*model.PaginatedExpenseResponse, error) {
	if err := generalUtil.ConfirmUserID(ctx, params.UserID); err != nil {
		return nil, utils.NewGQLError(err.(errapi.Error))
	}

	cursor := ""
	limit := 0
	sortField := "date"
	sortOrder := "desc"
	if params.SortField != nil {
		sortField = string(*params.SortField)
	}
	if params.SortOrder != nil {
		sortOrder = string(*params.SortOrder)
	}
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	if params.Limit != nil {

		limit = int(*params.Limit)
	}

	query, err := generalUtil.ConstructQueryParams(params.UserID, cursor, limit, sortField, sortOrder)
	if err != nil {
		return nil, utils.NewGQLError(errapi.NewBadRequest(err.Error()))
	}

	expenses, err := r.getMultipleExpenseHandler.Handle(query)
	if err != nil {
		return nil, utils.NewGQLError(err.(errapi.Error))
	}

	return utils.NewPaginatedExpenseResponse(expenses, sortField), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
