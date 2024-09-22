package graph

import (
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type Resolver struct {
	getExpenseHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	getMultipleExpenseHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
	addExpenseHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	patchExpenseHandler       icmd.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
}

type ResolverConfig struct {
	GetExpenseHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	GetMultipleExpenseHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
	AddExpenseHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	PatchExpenseHandler       icmd.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
}

func NewResolver(c ResolverConfig) *Resolver {
	return &Resolver{
		getExpenseHandler:         c.GetExpenseHandler,
		getMultipleExpenseHandler: c.GetMultipleExpenseHandler,
		addExpenseHandler:         c.AddExpenseHandler,
		patchExpenseHandler:       c.PatchExpenseHandler,
	}

}
