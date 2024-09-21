package graph

import (
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type Resolver struct {
	getExpenseHandler iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
}

type ResolverConfig struct {
	GetExpenseHandler iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
}

func NewResolver(c ResolverConfig) *Resolver {
	return &Resolver{
		getExpenseHandler: c.GetExpenseHandler,
	}

}
