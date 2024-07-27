package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/beka-birhanu/finance-go/domain/common/hash"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
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

// TestUserRepository runs a suite of tests for the UserRepository.
func TestUserRepository(t *testing.T) {
	repo := NewUserRepository(nil) // Passing nil as we're using an in-memory implementation

	t.Run("CreateUser", func(t *testing.T) {
		err := repo.Add(user)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("CreateUserWithClashingUsername", func(t *testing.T) {
		err := repo.Add(user)
		if err == nil {
			t.Errorf("expected conflict error %v, got %v", erruser.UsernameConflict, err)
		}
		if !errors.Is(err, erruser.UsernameConflict) {
			t.Errorf("unexpected error: %v, %v", err, erruser.UsernameConflict)
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
