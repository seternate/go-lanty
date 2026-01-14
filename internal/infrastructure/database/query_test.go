package database

import (
	"reflect"
	"testing"
)

func TestGenerateColumnValueMap(t *testing.T) {
	tests := []struct {
		name           string
		resource       any
		columnBlacklist []string
		expectedKeys    []string
		expectedValues  map[string]any
		expectError    bool
	}{
		{
			name: "struct with db tags",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "name", "email"},
			expectedValues: map[string]any{
				"id":    1,
				"name":  "John",
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "struct without db tags uses NameMapper",
			resource: struct {
				ID    int
				Name  string
				Email string
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "name", "email"},
			expectedValues: map[string]any{
				"id":    1,
				"name":  "John",
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "mixed db tags and no tags",
			resource: struct {
				ID    int    `db:"id"`
				Name  string // no tag
				Email string `db:"email_address"`
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "name", "email_address"},
			expectedValues: map[string]any{
				"id":            1,
				"name":          "John",
				"email_address": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "includes zero values",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    0,
				Name:  "",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "name", "email"},
			expectedValues: map[string]any{
				"id":    0,
				"name":  "",
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "with column blacklist",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: []string{"id", "email"},
			expectedKeys:    []string{"name"},
			expectedValues: map[string]any{
				"name": "John",
			},
			expectError: false,
		},
		{
			name: "all columns blacklisted returns error",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: []string{"id", "name", "email"},
			expectedKeys:    nil,
			expectedValues:  nil,
			expectError:     true,
		},
		{
			name: "different data types",
			resource: struct {
				ID     int     `db:"id"`
				Active bool    `db:"active"`
				Score  float64 `db:"score"`
				Name   string  `db:"name"`
			}{
				ID:     42,
				Active: true,
				Score:  95.5,
				Name:   "Test",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "active", "score", "name"},
			expectedValues: map[string]any{
				"id":     42,
				"active": true,
				"score":  95.5,
				"name":   "Test",
			},
			expectError: false,
		},
		{
			name: "empty struct returns error",
			resource: struct {
			}{},
			columnBlacklist: nil,
			expectedKeys:    nil,
			expectedValues:  nil,
			expectError:     true,
		},
		{
			name: "single field struct",
			resource: struct {
				ID int `db:"id"`
			}{
				ID: 1,
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id"},
			expectedValues: map[string]any{
				"id": 1,
			},
			expectError: false,
		},
		{
			name: "blacklist with non-existent column",
			resource: struct {
				ID   int    `db:"id"`
				Name string `db:"name"`
			}{
				ID:   1,
				Name: "John",
			},
			columnBlacklist: []string{"nonexistent"},
			expectedKeys:    []string{"id", "name"},
			expectedValues: map[string]any{
				"id":   1,
				"name": "John",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateColumnValueMap(tt.resource, tt.columnBlacklist...)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("expected result map, got nil")
			}

			// Check that all expected keys are present
			for _, key := range tt.expectedKeys {
				if _, ok := result[key]; !ok {
					t.Errorf("expected key %q not found in result", key)
				}
			}

			// Check that result has no extra keys
			if len(result) != len(tt.expectedKeys) {
				t.Errorf("result has %d keys, expected %d", len(result), len(tt.expectedKeys))
			}

			// Check values match
			for key, expectedValue := range tt.expectedValues {
				actualValue, ok := result[key]
				if !ok {
					t.Errorf("key %q not found in result", key)
					continue
				}
				if !reflect.DeepEqual(actualValue, expectedValue) {
					t.Errorf("value for key %q = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestGenerateNonZeroColumnValueMap(t *testing.T) {
	tests := []struct {
		name           string
		resource       any
		columnBlacklist []string
		expectedKeys    []string
		expectedValues  map[string]any
		expectError    bool
	}{
		{
			name: "excludes zero values",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    0,
				Name:  "",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"email"},
			expectedValues: map[string]any{
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "includes non-zero values",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id", "name", "email"},
			expectedValues: map[string]any{
				"id":    1,
				"name":  "John",
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "all zero values returns error",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    0,
				Name:  "",
				Email: "",
			},
			columnBlacklist: nil,
			expectedKeys:    nil,
			expectedValues:  nil,
			expectError:     true,
		},
		{
			name: "zero values with blacklist",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "",
				Email: "john@example.com",
			},
			columnBlacklist: []string{"id"},
			expectedKeys:    []string{"email"},
			expectedValues: map[string]any{
				"email": "john@example.com",
			},
			expectError: false,
		},
		{
			name: "zero values excluded even if not blacklisted",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    0,
				Name:  "John",
				Email: "",
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"name"},
			expectedValues: map[string]any{
				"name": "John",
			},
			expectError: false,
		},
		{
			name: "boolean false is zero value",
			resource: struct {
				ID     int  `db:"id"`
				Active bool `db:"active"`
				Count  int  `db:"count"`
			}{
				ID:     1,
				Active: false,
				Count:  0,
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id"},
			expectedValues: map[string]any{
				"id": 1,
			},
			expectError: false,
		},
		{
			name: "boolean true is not zero value",
			resource: struct {
				ID     int  `db:"id"`
				Active bool `db:"active"`
			}{
				ID:     0,
				Active: true,
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"active"},
			expectedValues: map[string]any{
				"active": true,
			},
			expectError: false,
		},
		{
			name: "float zero value excluded",
			resource: struct {
				ID    int     `db:"id"`
				Score float64 `db:"score"`
			}{
				ID:    1,
				Score: 0.0,
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"id"},
			expectedValues: map[string]any{
				"id": 1,
			},
			expectError: false,
		},
		{
			name: "non-zero float included",
			resource: struct {
				ID    int     `db:"id"`
				Score float64 `db:"score"`
			}{
				ID:    0,
				Score: 95.5,
			},
			columnBlacklist: nil,
			expectedKeys:    []string{"score"},
			expectedValues: map[string]any{
				"score": 95.5,
			},
			expectError: false,
		},
		{
			name: "all fields zero or blacklisted returns error",
			resource: struct {
				ID    int    `db:"id"`
				Name  string `db:"name"`
				Email string `db:"email"`
			}{
				ID:    1,
				Name:  "",
				Email: "",
			},
			columnBlacklist: []string{"id"},
			expectedKeys:    nil,
			expectedValues:  nil,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateNonZeroColumnValueMap(tt.resource, tt.columnBlacklist...)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("expected result map, got nil")
			}

			// Check that all expected keys are present
			for _, key := range tt.expectedKeys {
				if _, ok := result[key]; !ok {
					t.Errorf("expected key %q not found in result", key)
				}
			}

			// Check that result has no extra keys
			if len(result) != len(tt.expectedKeys) {
				t.Errorf("result has %d keys, expected %d. Result: %v", len(result), len(tt.expectedKeys), result)
			}

			// Check values match
			for key, expectedValue := range tt.expectedValues {
				actualValue, ok := result[key]
				if !ok {
					t.Errorf("key %q not found in result", key)
					continue
				}
				if !reflect.DeepEqual(actualValue, expectedValue) {
					t.Errorf("value for key %q = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestGenerateColumnValueMap_EdgeCases(t *testing.T) {
	t.Run("nil resource returns error", func(t *testing.T) {
		result, err := GenerateColumnValueMap(nil)
		if err == nil {
			t.Error("expected error for nil resource, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result on error, got %v", result)
		}
	})

	t.Run("non-struct type returns error", func(t *testing.T) {
		result, err := GenerateColumnValueMap("not a struct")
		if err == nil {
			t.Error("expected error for non-struct type, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result on error, got %v", result)
		}
	})

	t.Run("nil pointer returns error", func(t *testing.T) {
		type TestStruct struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}
		var ptr *TestStruct = nil

		result, err := GenerateColumnValueMap(ptr)
		if err == nil {
			t.Error("expected error for nil pointer, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result on error, got %v", result)
		}
	})

	t.Run("pointer to struct works correctly", func(t *testing.T) {
		type TestStruct struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}
		ptr := &TestStruct{
			ID:   1,
			Name: "John",
		}

		result, err := GenerateColumnValueMap(ptr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if result == nil {
			t.Fatal("expected result map, got nil")
		}
		expected := map[string]any{
			"id":   1,
			"name": "John",
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("result = %v, want %v", result, expected)
		}
	})

	t.Run("empty blacklist vs nil blacklist", func(t *testing.T) {
		resource := struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}{
			ID:   1,
			Name: "John",
		}

		result1, err1 := GenerateColumnValueMap(resource)
		if err1 != nil {
			t.Fatalf("unexpected error: %v", err1)
		}

		result2, err2 := GenerateColumnValueMap(resource, []string{}...)
		if err2 != nil {
			t.Fatalf("unexpected error: %v", err2)
		}

		if !reflect.DeepEqual(result1, result2) {
			t.Errorf("nil blacklist and empty blacklist should produce same result")
		}
	})
}
