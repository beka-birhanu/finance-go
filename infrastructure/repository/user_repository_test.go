package repository

import (
	"errors"
	"testing"
	"time"

	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/google/uuid"
)

type MockHashService struct {
}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return false, nil
}

var _ hash.IHashService = &MockHashService{}

var user, _ = model.NewUser("validUser", "#%strongPassword#%", &MockHashService{}, time.Now().UTC())

func TestUserRepository(t *testing.T) {
	repo := NewUserRepository(nil) // Passing nil as we're using an in-memory implementation

	t.Run("CreateUser", func(t *testing.T) {
		err := repo.CreateUser(user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("CreateUserWithClashingUsername", func(t *testing.T) {
		err := repo.CreateUser(user)
		if err == nil {
			t.Errorf("expected conflict error %v, got %v", appError.ErrUsernameConflict, err)
		}
		if !errors.Is(err, appError.ErrUsernameConflict) {
			t.Errorf("unexpected error: %v, %v", err, appError.ErrUsernameConflict)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		createdUser, err := repo.GetUserById(user.ID())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if createdUser == nil { //nolint
			t.Error("expected user to be found")
		}
		if createdUser.ID() != user.ID() { //nolint
			t.Errorf("expected user ID to be %v, got %v", user.ID(), createdUser.ID())
		}
	})

	t.Run("GetUserByIdWithInvalidId", func(t *testing.T) {
		_, err := repo.GetUserById(uuid.New()) // random invalid id
		if err == nil {
			t.Errorf("expected error %v", appError.ErrUserNotFound)
		}

	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		createdUser, err := repo.GetUserByUsername(user.Username())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if createdUser == nil { //nolint
			t.Error("expected user to be found")
		}
		if createdUser.Username() != user.Username() { //nolint
			t.Errorf("expected username to be %v, got %v", user.Username(), createdUser.Username())
		}
	})

	t.Run("GetUserByInvalidUsername", func(t *testing.T) {
		_, err := repo.GetUserByUsername("invalidUsername")
		if err == nil {
			t.Errorf("expected error %v", appError.ErrUserNotFound)
		}
	})

	t.Run("ListUser", func(t *testing.T) {
		users, err := repo.ListUser()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(users) == 0 {
			t.Error("expected users to be listed")
		}
	})
}
