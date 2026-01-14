package user

import (
	"errors"
	"net"
	"testing"

	domainerr "github.com/seternate/go-lanty/internal/domain/error"
)

func TestNewUser(t *testing.T) {
	validIPv4 := "192.168.1.1"
	validIPv4IP := net.ParseIP(validIPv4)

	t.Run("Valid user", func(t *testing.T) {
		user, err := NewUser("testuser", validIPv4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.Username != "testuser" {
			t.Errorf("expected Username %q, got %q", "testuser", user.Username)
		}
		if !user.IPv4Address.Equal(validIPv4IP) {
			t.Errorf("expected IPv4Address %v, got %v", validIPv4IP, user.IPv4Address)
		}
	})

	t.Run("Empty username", func(t *testing.T) {
		user, err := NewUser("", validIPv4)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid username format", func(t *testing.T) {
		invalidUsernames := []string{
			"test user",  // space
			"test@user",  // special char
			"test.user",  // dot
			"test#user",  // hash
			"test$user",  // dollar
		}
		for _, username := range invalidUsernames {
			t.Run(username, func(t *testing.T) {
				user, err := NewUser(username, validIPv4)
				if err == nil {
					t.Error("expected error, got nil")
				}
				if user != nil {
					t.Error("expected nil user")
				}
			})
		}
	})

	t.Run("Valid username formats", func(t *testing.T) {
		validUsernames := []string{
			"testuser",
			"test_user",
			"test-user",
			"test123",
			"123test",
			"test_user-123",
			"a",
			"User123",
		}
		for _, username := range validUsernames {
			t.Run(username, func(t *testing.T) {
				user, err := NewUser(username, validIPv4)
				if err != nil {
					t.Errorf("unexpected error for username %q: %v", username, err)
				}
				if user == nil {
					t.Errorf("expected user for username %q, got nil", username)
				}
			})
		}
	})

	t.Run("Username too long", func(t *testing.T) {
		longUsername := ""
		for i := 0; i < 65; i++ {
			longUsername += "a"
		}
		user, err := NewUser(longUsername, validIPv4)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
	})


	t.Run("Invalid IPv4 address (IPv6)", func(t *testing.T) {
		ipv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		user, err := NewUser("testuser", ipv6)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid IPv4 address (empty string)", func(t *testing.T) {
		user, err := NewUser("testuser", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid IPv4 address (invalid format)", func(t *testing.T) {
		user, err := NewUser("testuser", "not.an.ip.address")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
	})

	t.Run("Valid IPv4 addresses", func(t *testing.T) {
		validIPs := []string{
			"192.168.1.1",
			"10.0.0.1",
			"172.16.0.1",
			"127.0.0.1",
			"0.0.0.0",
			"255.255.255.255",
		}
		for _, ipStr := range validIPs {
			t.Run(ipStr, func(t *testing.T) {
				user, err := NewUser("testuser", ipStr)
				if err != nil {
					t.Errorf("unexpected error for IP %q: %v", ipStr, err)
				}
				if user == nil {
					t.Errorf("expected user for IP %q, got nil", ipStr)
				}
			})
		}
	})
}

func TestRehydrateUser(t *testing.T) {
	validIPv4 := "192.168.1.1"
	validIPv4IP := net.ParseIP(validIPv4)

	t.Run("Valid user", func(t *testing.T) {
		user, err := RehydrateUser("testuser", validIPv4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.Username != "testuser" {
			t.Errorf("expected Username %q, got %q", "testuser", user.Username)
		}
		if !user.IPv4Address.Equal(validIPv4IP) {
			t.Errorf("expected IPv4Address %v, got %v", validIPv4IP, user.IPv4Address)
		}
	})

	t.Run("Empty username returns TrustedInvariantViolationError", func(t *testing.T) {
		user, err := RehydrateUser("", validIPv4)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})

	t.Run("Invalid IPv4 returns TrustedInvariantViolationError", func(t *testing.T) {
		ipv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		user, err := RehydrateUser("testuser", ipv6)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user")
		}
		var trustedErr *domainerr.TrustedInvariantViolationError
		if !errors.As(err, &trustedErr) {
			t.Errorf("expected TrustedInvariantViolationError, got %T", err)
		}
	})
}

func TestUser_SetUsername(t *testing.T) {
	validIPv4 := "192.168.1.1"

	t.Run("Valid username", func(t *testing.T) {
		user, _ := NewUser("olduser", validIPv4)
		err := user.SetUsername("newuser")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.Username != "newuser" {
			t.Errorf("expected Username %q, got %q", "newuser", user.Username)
		}
	})

	t.Run("Empty username", func(t *testing.T) {
		user, _ := NewUser("testuser", validIPv4)
		err := user.SetUsername("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
		// Username should not be changed
		if user.Username != "testuser" {
			t.Errorf("expected Username to remain %q, got %q", "testuser", user.Username)
		}
	})

	t.Run("Invalid username format", func(t *testing.T) {
		user, _ := NewUser("testuser", validIPv4)
		err := user.SetUsername("test user")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
		// Username should not be changed
		if user.Username != "testuser" {
			t.Errorf("expected Username to remain %q, got %q", "testuser", user.Username)
		}
	})

	t.Run("Username too long", func(t *testing.T) {
		user, _ := NewUser("testuser", validIPv4)
		longUsername := ""
		for i := 0; i < 65; i++ {
			longUsername += "a"
		}
		err := user.SetUsername(longUsername)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var invariantErr *domainerr.InvariantViolationError
		if !errors.As(err, &invariantErr) {
			t.Errorf("expected InvariantViolationError, got %T", err)
		}
		// Username should not be changed
		if user.Username != "testuser" {
			t.Errorf("expected Username to remain %q, got %q", "testuser", user.Username)
		}
	})
}

func TestValidateUsername(t *testing.T) {
	t.Run("Valid usernames", func(t *testing.T) {
		validUsernames := []string{
			"testuser",
			"test_user",
			"test-user",
			"test123",
			"123test",
			"test_user-123",
			"a",
			"User123",
			"user_name",
			"user-name",
		}
		for _, username := range validUsernames {
			t.Run(username, func(t *testing.T) {
				err := validateUsername(username)
				if err != nil {
					t.Errorf("unexpected error for username %q: %v", username, err)
				}
			})
		}
	})

	t.Run("Invalid usernames", func(t *testing.T) {
		invalidUsernames := []string{
			"",
			"test user",  // space
			"test@user",  // special char
			"test.user",  // dot
			"test#user",  // hash
			"test$user",  // dollar
			"test%user",  // percent
			"test&user",  // ampersand
			"test*user",  // asterisk
			"test+user",  // plus
			"test=user",  // equals
			"test(user",  // parenthesis
			"test)user",  // parenthesis
			"test[user",  // bracket
			"test]user",  // bracket
			"test{user",  // brace
			"test}user",  // brace
			"test|user",  // pipe
			"test\\user", // backslash
			"test/user",  // slash
			"test:user",  // colon
			"test;user",  // semicolon
			"test<user",  // less than
			"test>user",  // greater than
			"test,user",  // comma
			"test?user",  // question mark
			"test!user",  // exclamation
			"test~user",  // tilde
			"test`user",  // backtick
			"test^user",  // caret
		}
		for _, username := range invalidUsernames {
			t.Run(username, func(t *testing.T) {
				err := validateUsername(username)
				if err == nil {
					t.Errorf("expected error for username %q, got nil", username)
					return
				}
				var validationErrs *domainerr.ValidationErrors
				if !errors.As(err, &validationErrs) {
					t.Errorf("expected ValidationErrors, got %T", err)
				}
			})
		}
	})

	t.Run("Username too long", func(t *testing.T) {
		longUsername := ""
		for i := 0; i < 65; i++ {
			longUsername += "a"
		}
		err := validateUsername(longUsername)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		var validationErrs *domainerr.ValidationErrors
		if !errors.As(err, &validationErrs) {
			t.Errorf("expected ValidationErrors, got %T", err)
		}
	})
}

func TestValidateIPv4(t *testing.T) {
	t.Run("Valid IPv4 addresses", func(t *testing.T) {
		validIPs := []string{
			"192.168.1.1",
			"10.0.0.1",
			"172.16.0.1",
			"127.0.0.1",
			"0.0.0.0",
			"255.255.255.255",
			"192.168.0.100",
			"10.10.10.10",
		}
		for _, ipStr := range validIPs {
			t.Run(ipStr, func(t *testing.T) {
				ip, err := validateIPv4(ipStr)
				if err != nil {
					t.Errorf("unexpected error for IP %q: %v", ipStr, err)
					return
				}
				if ip == nil {
					t.Errorf("expected non-nil IP for %q, got nil", ipStr)
					return
				}
				expectedIP := net.ParseIP(ipStr).To4()
				if !ip.Equal(expectedIP) {
					t.Errorf("expected IP %v, got %v", expectedIP, ip)
				}
			})
		}
	})

	t.Run("Invalid IP (empty string)", func(t *testing.T) {
		ip, err := validateIPv4("")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if ip != nil {
			t.Error("expected nil IP, got non-nil")
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "IPv4 address" {
			t.Errorf("expected Field %q, got %q", "IPv4 address", validationErr.Field)
		}
	})

	t.Run("Invalid IP (invalid format)", func(t *testing.T) {
		ip, err := validateIPv4("not.an.ip.address")
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if ip != nil {
			t.Error("expected nil IP, got non-nil")
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "IPv4 address" {
			t.Errorf("expected Field %q, got %q", "IPv4 address", validationErr.Field)
		}
	})

	t.Run("Invalid IPv4 (IPv6)", func(t *testing.T) {
		ipv6Str := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		ip, err := validateIPv4(ipv6Str)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if ip != nil {
			t.Error("expected nil IP, got non-nil")
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
		if validationErr != nil && validationErr.Field != "IPv4 address" {
			t.Errorf("expected Field %q, got %q", "IPv4 address", validationErr.Field)
		}
	})

	t.Run("Invalid IPv4 (IPv6 loopback)", func(t *testing.T) {
		ipv6Str := "::1"
		ip, err := validateIPv4(ipv6Str)
		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if ip != nil {
			t.Error("expected nil IP, got non-nil")
		}
		var validationErr *domainerr.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})
}
