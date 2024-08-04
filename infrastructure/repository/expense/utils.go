package expenserepo

import (
	"fmt"
	"time"

	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// ScanExpense scans a database row into an Expense model.
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

func BuildExpenseListWhereClause(ascending bool, id uuid.UUID, time time.Time, params *[]interface{}) string {
	inequalitySign := "<"
	if ascending {
		inequalitySign = ">"
	}

	clause := fmt.Sprintf(`
		AND id < $%v 
		AND created_at %s $%v
		`, len(*params)+1, inequalitySign, len(*params)+2)

	*params = append(*params, id, time)
	return clause
}

func BuildExpenseListOrderByClause(ascending bool) string {
	order := "DESC"
	if ascending {
		order = "ASC"
	}

	return fmt.Sprintf("ORDER BY id DESC, created_at %s", order)
}

func BuildLimitClause(limit int, params *[]interface{}) string {
	clause := fmt.Sprintf("LIMIT $%v", len(*params)+1)
	*params = append(*params, limit)

	return clause
}
