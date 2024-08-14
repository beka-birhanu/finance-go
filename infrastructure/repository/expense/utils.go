// Package expenserepo provides utilities for working with expense models and building SQL queries.
package expenserepo

import (
	"fmt"
	"time"

	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// ScanExpense converts a database row into an Expense model.
func ScanExpense(scanner interface {
	Scan(dest ...interface{}) error
}) (*expensemodel.Expense, error) {
	var id, userId uuid.UUID
	var description string
	var amount float32
	var date, createdAt, updatedAt time.Time

	err := scanner.Scan(&id, &description, &amount, &date, &userId, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	config := expensemodel.Config{
		Description:  description,
		Amount:       amount,
		UserId:       userId,
		Date:         date,
		CreationTime: createdAt,
	}

	expense, err := expensemodel.NewWithID(id, config)
	if err != nil {
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error creating expense model: %v", err))
	}

	expense.UpdateDate(date)
	_ = expense.UpdateDescription(description)
	_ = expense.UpdateAmount(amount)

	return expense, nil
}

// BuildExpenseListWhereClause creates the WHERE clause for expense pagination queries.
func BuildExpenseListWhereClause(ascending bool, id uuid.UUID, value interface{}, field string, params *[]interface{}) string {
	if id == uuid.Nil {
		return ""
	}

	inequalitySign := "<"
	if ascending {
		inequalitySign = ">"
	}

	clause := fmt.Sprintf(`
		AND id < $%v 
		AND %s %s $%v
		`, len(*params)+1, field, inequalitySign, len(*params)+2)

	*params = append(*params, id, value)
	return clause
}

// BuildExpenseListOrderByClause creates the ORDER BY clause for expense pagination queries.
func BuildExpenseListOrderByClause(ascending bool, field string) string {
	order := "DESC"
	if ascending {
		order = "ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s, id DESC", field, order)
}

// BuildLimitClause creates the LIMIT clause for expense pagination queries.
func BuildLimitClause(limit int, params *[]interface{}) string {
	clause := fmt.Sprintf("LIMIT $%v", len(*params)+1)
	*params = append(*params, limit)

	return clause
}
