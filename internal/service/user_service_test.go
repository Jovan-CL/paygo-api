package service

import (
	"context"
	"testing"
	"time"

	"paygo-api/internal/domain"
)

// 1. CREATE A MOCK REPOSITORY
// We build a dummy struct that tricks the service into thinking it is talking to Postgres
type mockUserRepo struct {
	domain.UserRepository // Embed the interface contract
	mockUser              *domain.User
	mockErr               error
}

// Implement the specific function your login logic calls
func (m *mockUserRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.mockUser, m.mockErr
}

// 2. WRITE THE ACTUAL TEST FUNCTION
// Test functions must always start with the word "Test" and take *testing.T
func TestLogin_DeactivatedUser(t *testing.T) {
	// Arrange: Create a deactivated user scenario
	deactivationTime := time.Date(2026, 6, 8, 17, 0, 0, 0, time.UTC)

	deactivatedUser := &domain.User{
		Email:        "jovan@gmail.com",
		PasswordHash: "any_hash",
		Status:       "deactivated",
		DeletedAt:    deactivationTime, // Or a regular time if using Option A
	}

	// Inject the deactivated user into our mock database
	mockRepo := &mockUserRepo{mockUser: deactivatedUser, mockErr: nil}

	// Instantiate the actual service engine with our fake database linked
	userService := NewUserService(mockRepo)

	// Act: Trigger the Login function under test
	_, err := userService.Login(context.Background(), "jovan@gmail.com", "password123")

	// Assert: Verify the results match your exact expectations
	if err == nil {
		t.Fatal("Expected an error for a deactivated user, but got nil")
	}

	expectedErrorMsg := "this user account was deactivated at Jun 08, 2026 - 05:00 PM"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Wrong error message returned.\nGot:  %s\nWant: %s", err.Error(), expectedErrorMsg)
	}
}
