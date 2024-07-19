package repositories

import (
	"database/sql"

	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/entities"
)

type UserRepository struct {
	DB *sql.DB
}

// Ensure UserRepository implements interfaces.persistance.IUserRepository
var _ persistance.IUserRepository = &UserRepository{}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// CreateExpense implements persistance.IUserRepository.
func (u *UserRepository) CreateExpense(expense *entities.User) error {
	panic("unimplemented")
}

// GetExpense implements persistance.IUserRepository.
func (u *UserRepository) GetExpense(id string) (*entities.User, error) {
	panic("unimplemented")
}

// ListExpenses implements persistance.IUserRepository.
func (u *UserRepository) ListExpenses() ([]*entities.User, error) {
	panic("unimplemented")
}

// package repositories
//
// import (
// 	"errors"
// 	"sync"
//
// 	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
// 	"github.com/beka-birhanu/finance-go/domain/entities"
// )
//
// type InMemoryUserRepository struct {
// 	mu    sync.RWMutex
// 	users map[string]*entities.User
// }
//
// // Ensure InMemoryUserRepository implements persistance.IUserRepository
// var _ persistance.IUserRepository = &InMemoryUserRepository{}
//
// func NewInMemoryUserRepository() *InMemoryUserRepository {
// 	return &InMemoryUserRepository{
// 		users: make(map[string]*entities.User),
// 	}
// }
//
// // CreateUser implements persistance.IUserRepository
// func (r *InMemoryUserRepository) CreateUser(user *entities.User) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
//
// 	if _, exists := r.users[user.ID]; exists {
// 		return errors.New("user already exists")
// 	}
//
// 	r.users[user.ID] = user
// 	return nil
// }
//
// // GetUser implements persistance.IUserRepository
// func (r *InMemoryUserRepository) GetUser(id string) (*entities.User, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
//
// 	user, exists := r.users[id]
// 	if !exists {
// 		return nil, errors.New("user not found")
// 	}
//
// 	return user, nil
// }
//
// // ListUsers implements persistance.IUserRepository
// func (r *InMemoryUserRepository) ListUsers() ([]*entities.User, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
//
// 	var users []*entities.User
// 	for _, user := range r.users {
// 		users = append(users, user)
// 	}
//
// 	return users, nil
// }
