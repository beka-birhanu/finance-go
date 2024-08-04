package expenserepo

import (
	"database/sql"
	"fmt"

	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	errexpense "github.com/beka-birhanu/finance-go/domain/error/expense"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

var _ irepository.IExpenseRepository = &Repository{}

const listBaseQuery = `
		SELECT id, description, amount, date, user_id, created_at, updated_at
		FROM expenses
		WHERE user_id = $1
	`

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save inserts or updates an expense in the database.
// If the expense already exists, it updates the existing record.
func (e *Repository) Save(expense *expensemodel.Expense) error {
	_, err := e.db.Exec(`
		INSERT INTO expenses (id, description, amount, date, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id, user_id) DO UPDATE
		SET description = EXCLUDED.description,
				amount = EXCLUDED.amount,
				date = EXCLUDED.date,
				updated_at = EXCLUDED.updated_at`,
		expense.ID(), expense.Description(), expense.Amount(), expense.Date(), expense.UserID(), expense.CreatedAt(), expense.UpdatedAt())

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("conflict error: expense with ID %s and user_id %s already exists", expense.ID(), expense.UserID())
		}
		return errdmn.NewUnexpected(fmt.Sprintf("error saving expense: %v", err))
	}
	return nil
}

// ById retrieves an expense by its unique identifier and user ID.
func (e *Repository) ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error) {
	row := e.db.QueryRow(`
		SELECT id, description, amount, date, user_id, created_at, updated_at
		FROM expenses
		WHERE id = $1 AND user_id = $2`, id, userId)

	expense, err := ScanExpense(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errexpense.NotFound
		}
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error retrieving expense: %v", err))
	}

	return expense, nil
}

// ListByTime retrieves paginated expenses for a user based on creation time.
//
// Params:
// - params: A struct containing:
//   - UserID: UUID of the user whose expenses are being queried.
//   - Limit: Maximum number of expenses to retrieve.
//   - LastSeenID: UUID of the last seen expense to start pagination from.
//   - LastSeenTime: Time of the last seen expense to start pagination from.
//   - Ascending: Boolean to determine the order of sorting (true for ascending, false for descending).
//
// Returns:
// - A slice of Expense pointers.
// - An error, if any occurs during the query execution or scanning.
func (e *Repository) ListByTime(params irepository.ListByTimeParams) ([]*expensemodel.Expense, error) {
	queryParams := []interface{}{params.UserID}
	additionalWhere := BuildExpenseListWhereClause(params.Ascending, *params.LastSeenID, params.LastSeenTime, "created_at", &queryParams)
	orderBy := BuildExpenseListOrderByClause(params.Ascending, "created_at")
	limitClause := BuildLimitClause(params.Limit, &queryParams)

	query := fmt.Sprintf("%s %s %s %s", listBaseQuery, additionalWhere, orderBy, limitClause)
	rows, err := e.db.Query(query, queryParams...)
	if err != nil {
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error listing expenses: %v", err))
	}
	defer rows.Close()

	var expenses []*expensemodel.Expense
	for rows.Next() {
		expense, err := ScanExpense(rows)
		if err != nil {
			return nil, errdmn.NewUnexpected(fmt.Sprintf("error scanning expense: %v", err))
		}
		expenses = append(expenses, expense)
	}
	if err = rows.Err(); err != nil {
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error with rows: %v", err))
	}
	return expenses, nil
}

// Params:
// - params: A struct containing:
//   - UserID: UUID of the user whose expenses are being queried.
//   - Limit: Maximum number of expenses to retrieve.
//   - LastSeenID: UUID of the last seen expense to start pagination from.
//   - LastSeenTime: Time of the last seen expense to start pagination from.
//   - Ascending: Boolean to determine the order of sorting (true for ascending, false for descending).
//
// Returns:
// - A slice of Expense pointers.
// - An error, if any occurs during the query execution or scanning.
func (e *Repository) ListByAmount(params irepository.ListByAmountParams) ([]*expensemodel.Expense, error) {
	queryParams := []interface{}{params.UserID}
	additionalWhere := BuildExpenseListWhereClause(params.Ascending, *params.LastSeenID, params.LastSeenAmt, "amount", &queryParams)
	orderBy := BuildExpenseListOrderByClause(params.Ascending, "amount")
	limitClause := BuildLimitClause(params.Limit, &queryParams)

	query := fmt.Sprintf("%s %s %s %s", listBaseQuery, additionalWhere, orderBy, limitClause)
	rows, err := e.db.Query(query, queryParams...)
	if err != nil {
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error listing expenses: %v", err))
	}
	defer rows.Close()

	var expenses []*expensemodel.Expense
	for rows.Next() {
		expense, err := ScanExpense(rows)
		if err != nil {
			return nil, errdmn.NewUnexpected(fmt.Sprintf("error scanning expense: %v", err))
		}
		expenses = append(expenses, expense)
	}
	if err = rows.Err(); err != nil {
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error with rows: %v", err))
	}
	return expenses, nil
}
