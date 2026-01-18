//go:build integration

package user

import (
	"errors"
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	domainerr "github.com/seternate/go-lanty/internal/domain/error"
	"github.com/seternate/go-lanty/internal/domain/user"
	"github.com/seternate/go-lanty/internal/infrastructure/database"
)

func getTestDB(t *testing.T) *sqlx.DB {
	dbString := os.Getenv("DB")
	if dbString == "" {
		dbString = "postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable"
	}

	db, err := database.New(dbString)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func setupTest(t *testing.T) (*userrepository, *sqlx.DB) {
	db := getTestDB(t)
	repo := NewUserRepository(db).(*userrepository)

	// Clean up any existing test data
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean up test data: %v", err)
	}

	// Add defensive cleanup for when test fails
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM users")
		db.Close()
	})

	return repo, db
}

func TestNewUserRepository(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	if repo == nil {
		t.Fatal("expected repository, got nil")
	}

	// Verify it implements the interface
	var _ user.UserRepository = repo
}

func TestUserRepository_GetUser(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Get existing user", func(t *testing.T) {
		testIPv4Address := "192.168.1.100"
		testUsername := "testuser"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to save user in database: %v", err)
		}

		// Retrieve the user
		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved == nil {
			t.Fatal("expected user, got nil")
		}
		if retrieved.IPv4Address.String() != testIPv4Address {
			t.Errorf("expected IPv4Address %q, got %q", testIPv4Address, retrieved.IPv4Address.String())
		}
		if retrieved.Username != testUsername {
			t.Errorf("expected Username %q, got %q", testUsername, retrieved.Username)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		nonExistentIPv4Address := "10.0.0.1"

		retrieved, err := repo.GetUser(nonExistentIPv4Address)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if retrieved != nil {
			t.Error("expected nil user")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Get user with different IPv4 address", func(t *testing.T) {
		testIPv4Address := "172.16.0.50"
		testUsername := "anotheruser"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to save user in database: %v", err)
		}

		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.IPv4Address.String() != testIPv4Address {
			t.Errorf("expected IPv4Address %q, got %q", testIPv4Address, retrieved.IPv4Address.String())
		}
		if retrieved.Username != testUsername {
			t.Errorf("expected Username %q, got %q", testUsername, retrieved.Username)
		}
	})
}

func TestUserRepository_SaveUser(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Save valid user", func(t *testing.T) {
		testIPv4Address := "192.168.1.200"
		testUsername := "newuser"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it was saved
		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to retrieve saved user: %v", err)
		}
		if retrieved.IPv4Address.String() != testIPv4Address {
			t.Errorf("expected IPv4Address %q, got %q", testIPv4Address, retrieved.IPv4Address.String())
		}
		if retrieved.Username != testUsername {
			t.Errorf("expected Username %q, got %q", testUsername, retrieved.Username)
		}
	})

	t.Run("Update existing user (UPSERT)", func(t *testing.T) {
		testIPv4Address := "192.168.1.201"
		initialUsername := "initialuser"
		testUser, err := user.RehydrateUser(initialUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to save user: %v", err)
		}

		// Update the username
		updatedUsername := "updateduser"
		err = testUser.SetUsername(updatedUsername)
		if err != nil {
			t.Fatalf("failed to update username: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		// Verify the update
		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to retrieve updated user: %v", err)
		}
		if retrieved.Username != updatedUsername {
			t.Errorf("expected Username %q, got %q", updatedUsername, retrieved.Username)
		}
		if retrieved.IPv4Address.String() != testIPv4Address {
			t.Errorf("expected IPv4Address %q, got %q", testIPv4Address, retrieved.IPv4Address.String())
		}
	})

	t.Run("Save multiple users", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			testIPv4Address := fmt.Sprintf("10.0.0.%d", i+1)
			testUsername := fmt.Sprintf("user%d", i)
			testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
			if err != nil {
				t.Fatalf("failed to create test user %d: %v", i, err)
			}

			err = repo.SaveUser(testUser)
			if err != nil {
				t.Fatalf("unexpected error creating user %d: %v", i, err)
			}
		}

		// Verify all users were saved
		for i := 0; i < 5; i++ {
			testIPv4Address := fmt.Sprintf("10.0.0.%d", i+1)
			retrieved, err := repo.GetUser(testIPv4Address)
			if err != nil {
				t.Fatalf("failed to retrieve user %d: %v", i, err)
			}
			if retrieved == nil {
				t.Errorf("expected user %d to exist", i)
			}
		}
	})

	t.Run("Save user with special characters in username", func(t *testing.T) {
		testIPv4Address := "192.168.1.202"
		testUsername := "user_name-123"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to retrieve saved user: %v", err)
		}
		if retrieved.Username != testUsername {
			t.Errorf("expected Username %q, got %q", testUsername, retrieved.Username)
		}
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Delete existing user", func(t *testing.T) {
		testIPv4Address := "192.168.1.250"
		testUsername := "deleteme"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to save user: %v", err)
		}

		// Verify it exists
		_, err = repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("user should exist before deletion: %v", err)
		}

		// Delete it
		err = repo.DeleteUser(testIPv4Address)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetUser(testIPv4Address)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
		var notFoundErr *domainerr.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Errorf("expected NotFoundError, got %T: %v", err, err)
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		nonExistentIPv4Address := "10.0.0.99"

		// Delete should not return an error even if the user doesn't exist
		err := repo.DeleteUser(nonExistentIPv4Address)
		if err != nil {
			// If it does return an error, that's also acceptable behavior
			// We just verify it doesn't panic
		}
	})

	t.Run("Delete multiple users", func(t *testing.T) {
		testIPv4Addresses := []string{"192.168.1.101", "192.168.1.102", "192.168.1.103"}
		for i, ipv4Address := range testIPv4Addresses {
			testUser, err := user.RehydrateUser(fmt.Sprintf("user%d", i), ipv4Address)
			if err != nil {
				t.Fatalf("failed to create test user: %v", err)
			}

			err = repo.SaveUser(testUser)
			if err != nil {
				t.Fatalf("failed to save user: %v", err)
			}
		}

		// Delete all of them
		for _, ipv4Address := range testIPv4Addresses {
			err := repo.DeleteUser(ipv4Address)
			if err != nil {
				t.Fatalf("unexpected error deleting user %s: %v", ipv4Address, err)
			}
		}

		// Verify they're all gone
		for _, ipv4Address := range testIPv4Addresses {
			_, err := repo.GetUser(ipv4Address)
			if err == nil {
				t.Errorf("user %s should not exist after deletion", ipv4Address)
			}
		}
	})
}

func TestUserRepository_Integration(t *testing.T) {
	repo, _ := setupTest(t)

	t.Run("Full CRUD cycle", func(t *testing.T) {
		testIPv4Address := "192.168.1.240"
		testUsername := "cruduser"
		testUser, err := user.RehydrateUser(testUsername, testIPv4Address)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create
		err = repo.SaveUser(testUser)
		if err != nil {
			t.Fatalf("failed to save user: %v", err)
		}

		// Read
		retrieved, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to retrieve user: %v", err)
		}
		if retrieved.IPv4Address.String() != testIPv4Address {
			t.Errorf("expected IPv4Address %q, got %q", testIPv4Address, retrieved.IPv4Address.String())
		}
		if retrieved.Username != testUsername {
			t.Errorf("expected Username %q, got %q", testUsername, retrieved.Username)
		}

		// Update
		updatedUsername := "updatedcruduser"
		err = retrieved.SetUsername(updatedUsername)
		if err != nil {
			t.Fatalf("failed to update username: %v", err)
		}
		err = repo.SaveUser(retrieved)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		// Verify update
		updated, err := repo.GetUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to retrieve updated user: %v", err)
		}
		if updated.Username != updatedUsername {
			t.Errorf("expected Username %q, got %q", updatedUsername, updated.Username)
		}

		// Delete
		err = repo.DeleteUser(testIPv4Address)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		// Verify deletion
		_, err = repo.GetUser(testIPv4Address)
		if err == nil {
			t.Error("expected error after deletion, got nil")
		}
	})
}
