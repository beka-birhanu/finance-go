package userrepo

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/beka-birhanu/finance-go/config"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// MockHashService is a mock implementation of the IHashService interface.
type MockHashService struct{}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return false, nil
}

var _ hash.IService = &MockHashService{}

var user, _ = usermodel.New(usermodel.Config{
	Username:       "validUser",
	PlainPassword:  "#%strongPassword#%",
	CreationTime:   time.Now().UTC(),
	PasswordHasher: &MockHashService{},
})

var userConflict, _ = usermodel.New(usermodel.Config{
	Username:       "validUser",
	PlainPassword:  "#%strongPassword#%",
	CreationTime:   time.Now().UTC(),
	PasswordHasher: &MockHashService{},
})

// CreateTestDB initializes a test database connection.
func CreateTestDB(t *testing.T) *sql.DB {
	connStr := "user=" + config.Envs.TestDBUser +
		" password=" + config.Envs.TestDBPassword +
		" dbname=" + config.Envs.TestDBName +
		" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	runMigrations(t, connStr)

	return db
}

// Run migrations on the database.
func runMigrations(t *testing.T, connStr string) {
	// Initialize the migrate instance
	m, err := migrate.New(
		"file://infrastructure/db/migrations", // Path to migration files
		connStr,                               // Connection string
	)
	if err != nil {
		t.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Hey there, migration has failed: %v", err)
	}
}

// Rollback migrations on the database.
func rollbackMigrations(t *testing.T, connStr string) {
	m, err := migrate.New(
		"file://infrastructure/db/migrations", // Path to migration files
		connStr,
	)
	if err != nil {
		t.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Failed to rollback migrations: %v", err)
	}
}

// TestUserRepository runs a suite of tests for the UserRepository.
func TestUserRepository(t *testing.T) {
	connStr := "user=test_user password=test_password dbname=test_finance sslmode=disable"
	db := CreateTestDB(t)
	defer func() {
		rollbackMigrations(t, connStr)
		db.Close()
	}()

	repo := New(db) // Pass the test database connection to the repository

	t.Run("SaveUser", func(t *testing.T) {
		err := repo.Save(user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("SaveUserWithClashingUsername", func(t *testing.T) {
		err := repo.Save(userConflict)
		if err == nil {
			t.Errorf("expected conflict error %v, got %v", erruser.UsernameConflict, err)
		}
		if !errors.Is(err, erruser.UsernameConflict) {
			t.Errorf("unexpected error: %v, %v", err, erruser.UsernameConflict)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Modify the user object
		testExpense, _ := expensemodel.New(expensemodel.Config{
			Description:  "asdfasdf",
			Amount:       43,
			UserId:       user.ID(),
			Date:         time.Now(),
			CreationTime: time.Now(),
		})

		err := user.AddExpense(testExpense, time.Now())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = repo.Save(user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		createdUser, err := repo.ById(user.ID())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if createdUser == nil {
			t.Error("expected user to be found")
		}
		if createdUser.ID() != user.ID() {
			t.Errorf("expected user ID to be %v, got %v", user.ID(), createdUser.ID())
		}
	})

	t.Run("GetUserByIdWithInvalidId", func(t *testing.T) {
		_, err := repo.ById(uuid.New()) // random invalid id
		if err == nil {
			t.Errorf("expected error %v", erruser.NotFound)
		}
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		createdUser, err := repo.ByUsername(user.Username())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if createdUser == nil {
			t.Error("expected user to be found")
		}
		if createdUser.Username() != user.Username() {
			t.Errorf("expected username to be %v, got %v", user.Username(), createdUser.Username())
		}
	})

	t.Run("GetUserByInvalidUsername", func(t *testing.T) {
		_, err := repo.ByUsername("invalidUsername")
		if err == nil {
			t.Errorf("expected error %v", erruser.NotFound)
		}
	})
}

