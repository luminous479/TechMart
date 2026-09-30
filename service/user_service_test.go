package service

import (
	"database/sql"
	"testing"

	"github.com/luminous479/TechMart/model"
)

type MockUserRepository struct {
	users map[int]*model.User
}

func (m *MockUserRepository) GetByID(id int) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

func TestUserService_GetByID_UserFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		users: map[int]*model.User{
			1: {
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
	}

	userService := NewUserService(mockRepo)

	user, err := userService.GetByID(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", user.Email)
	}
}

func TestUserService_GetByID_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		users: map[int]*model.User{},
	}

	userService := NewUserService(mockRepo)

	user, err := userService.GetByID(99)

	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
