// Package userrepo provides the implementation for handling user persistence operations in the repository.
package userrepo

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/lib/pq"

	// errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
)

// Repository handles the persistence of user models.
type Repository struct {
	db *sql.DB
}

// Ensure UserRepository implements repository.IUserRepository.
var _ irepository.IUserRepository = &Repository{}

// New creates a new UserRepository with the given database connection.
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save inserts or updates a user in the repository.
// If the user already exists, it updates the existing record.
// If the user does not exist, it adds a new record.
//
// Returns:
//   - error: An error if a conflict occurs, otherwise nil.
func (u *Repository) Save(user *usermodel.User) error {
	ctx, err := u.db.Begin()
	if err != nil {
		return errdmn.NewUnexpected(fmt.Sprintf("error starting transaction: %v", err))
	}

	defer func() {
		if err != nil {
			err = ctx.Rollback()
			if err != nil {
				log.Printf("error closing transaction: %v", err)
			}
		} else {
			err = ctx.Commit()
		}
	}()

	// Channels for handling errors from goroutines
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	// Upsert user in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := upsertUser(ctx, user); err != nil {
			errChan <- err
		}
	}()

	// Upsert expenses in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := upsertExpenses(ctx, user.Expenses()); err != nil {
			errChan <- err
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Collect and return errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// ById retrieves a user by their ID.
//
// Returns:
//   - *usermodel.User: A pointer to the retrieved user model.
//   - error: An error if the user is not found, otherwise nil.
func (u *Repository) ById(id uuid.UUID) (*usermodel.User, error) {
	row := u.db.QueryRow("SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = $1", id)

	user, err := scanRowToUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errdmn.NewNotFound("user not found.")
		}
		return nil, errdmn.NewUnexpected(fmt.Sprintf("error fetching user by their id: %v", err))
	}

	return user, nil
}

// ByUsername retrieves a user by their username.
//
// Returns:
//   - *usermodel.User: A pointer to the retrieved user model.
//   - error: An error if the user is not found, otherwise nil.
func (u *Repository) ByUsername(username string) (*usermodel.User, error) {
	row := u.db.QueryRow("SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1", username)

	user, err := scanRowToUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errdmn.NewNotFound("user not found.")
		}
		return nil, errdmn.NewUnexpected(err.Error())
	}

	return user, nil
}

// upsertUser inserts a new user or updates an existing user in the database.
// Returns a conflict error if the username already exists with a different ID.
func upsertUser(ctx *sql.Tx, user *usermodel.User) error {
	// Sorry for the lengthy query in advance. here is what the query do.
	// 1. Use a Common Table Expression (CTE) to check if a user with the same username
	//    but a different ID already exists.
	// 2. If such a user exists, the CTE returns the conflicting user's ID.
	// 3. Attempt to insert a new user with the provided details (id, username, password_hash,
	//    created_at, updated_at).
	// 4. On conflict with an existing user ID, update the user details with the new values.
	// 5. The update only occurs if no conflicting username is found (i.e., the CTE is empty).
	// 6. Return the ID of the inserted or updated user for further processing.
	//
	var userID string
	err := ctx.QueryRow(`
        WITH existing_user AS (
            SELECT id FROM users WHERE username = $1 AND id != $2
        )
        INSERT INTO users (id, username, password_hash, created_at, updated_at)
        VALUES ($2, $1, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            username = EXCLUDED.username,
            password_hash = EXCLUDED.password_hash,
            created_at = EXCLUDED.created_at,
            updated_at = EXCLUDED.updated_at
        WHERE NOT EXISTS (SELECT 1 FROM existing_user)
        RETURNING id`,
		user.Username(), user.ID(), user.PasswordHash(), user.CreatedAt(), user.UpdatedAt()).Scan(&userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				// At this point, we are sure that the conflict error is due to username because we checked for
				// user_id conflict in the query.
				return erruser.UsernameConflict
			}
		}
		// Handle other unexpected errors
		return errdmn.NewUnexpected(fmt.Sprintf("error saving user: %v", err))
	}

	return nil
}

// upsertExpenses inserts expenses into the database.
// If an expense with the same ID and user_id already exists, it returns a conflict error.
// Any other errors during insertion are also returned.
func upsertExpenses(ctx *sql.Tx, expenses []expensemodel.Expense) error {
	for _, expense := range expenses {
		_, err := ctx.Exec(`
            INSERT INTO expenses (id, description, amount, date, user_id, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			expense.ID(), expense.Description(), expense.Amount(), expense.Date(), expense.UserID(), expense.CreatedAt(), expense.UpdatedAt())

		if err != nil {
			// Check if the error is a unique constraint violation (conflict)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				return errdmn.NewUnexpected(fmt.Sprintf("conflict error: expense with ID %s and user_id %s already exists", expense.ID(), expense.UserID()))
			}
			// Return any other error that occurred during insertion
			return fmt.Errorf("error saving expense: %v", err)
		}
	}
	return nil
}

func scanRowToUser(row *sql.Row) (*usermodel.User, error) {
	var (
		id           uuid.UUID
		username     string
		passwordHash string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(&id, &username, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user, err := usermodel.NewWithExistingHash(usermodel.ConfigForExistingHash{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreationTime: createdAt,
		UpdatedAt:    updatedAt,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
