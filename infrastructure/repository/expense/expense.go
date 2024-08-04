package expenserepo

import (
	"database/sql"
	"fmt"
	"time"

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

// List retrieves all expenses for a given user ID.
func (e *Repository) List(userId uuid.UUID) ([]*expensemodel.Expense, error) {
	params := []interface{}{userId}
	id, _ := uuid.Parse("99193432-4063-4d43-a7d4-9b105221609b")
	additionalWhere := BuildExpenseListWhereClause(false, id, time.Now().UTC(), &params)
	orderBy := BuildExpenseListOrderByClause(false)
	limit := BuildLimitClause(5, &params)
	query := fmt.Sprintf("%s %s %s %s", listBaseQuery, additionalWhere, orderBy, limit)
	rows, err := e.db.Query(query, params...)
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
