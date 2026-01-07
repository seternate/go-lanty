package game

import (
	"testing"
)

// TestGameRepository is a test to ensure the GameRepository interface is properly defined.
// Since it's an interface, we can only verify that it compiles correctly.
// Actual implementation tests would be in the infrastructure layer.
func TestGameRepository_Interface(t *testing.T) {
	// This test ensures the interface is properly defined and compiles.
	// In a real scenario, you would test implementations of this interface
	// in the infrastructure/persistence layer.
	var _ GameRepository = (*mockGameRepository)(nil)
}

// mockGameRepository is a minimal mock implementation for testing the interface
type mockGameRepository struct{}

func (m *mockGameRepository) GetGame(slug string) (*Game, error) {
	return nil, nil
}

func (m *mockGameRepository) SaveGame(game *Game) error {
	return nil
}

func (m *mockGameRepository) DeleteGame(slug string) error {
	return nil
}
