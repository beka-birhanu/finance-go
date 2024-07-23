package repository

import (
	"errors"
	"testing"

	domainError "github.com/beka-birhanu/finance-go/domain/error"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type MockHashService struct {
}

func (m *MockHashService) Hash(word string) (string, error) {
	return word, nil
}

func (m *MockHashService) Match(hashedWord, plainWord string) (bool, error) {
	return false, nil
}

var user, _ = model.NewUser("validUser", "#%strongPassword#%", &MockHashService{})

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
			t.Errorf("expected conflict error %v, got %v", domainError.ErrUsernameConflict, err)
		}
		if !errors.Is(err, domainError.ErrUsernameConflict) {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		createdUser, err := repo.GetUserById(user.ID().String())
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
		_, err := repo.GetUserById("invalidId")
		if err == nil {
			t.Errorf("expected error %v", NotFound)
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
			t.Errorf("expected error %v", NotFound)
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
