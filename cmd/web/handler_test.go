package main_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockDB is a mock implementation of the database interface for testing purposes
type MockDB struct {
}

// QueryRow mocks the database query operation
func (m *MockDB) QueryRow(query string, args ...interface{}) mockRow {
	// Simulating a row with username "testuser" and password "testpassword"
	return mockRow{"mpit", 1, 1} // Assuming role and memberID values
}

// mockRow is a mock implementation of sql.Row for testing purposes
type mockRow struct {
	password string
	role     int
	memberID int
}

// Scan mocks the scanning operation on a row
func (m mockRow) Scan(dest ...interface{}) error {
	switch len(dest) {
	case 3:
		dest[0] = m.password
		dest[1] = m.role
		dest[2] = m.memberID
	default:
		return nil // Return nil error for simplicity
	}
	return nil // Return nil error for simplicity
}

func TestVerificationHandler(t *testing.T) {
	// Create a new instance of the MockDB
	mockDB := &MockDB{}

	// Create a new HTTP request to simulate a POST request
	reqBody := strings.NewReader("username=mpit&password=mpit")
	req, err := http.NewRequest("POST", "/verification", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the verification handler function with the mock database
	// Since we're using a mock database, we need to pass it to the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch user credentials from the database
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		var role int
		var memberID int
		// Call the QueryRow method of the mock database
		row := mockDB.QueryRow("SELECT password, role, memberID FROM LOGIN WHERE username = ?", username)
		err := row.Scan(&storedPassword, &role, &memberID)
		if err == sql.ErrNoRows || storedPassword != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Assuming verification succeeded, redirect to success page
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	})

	handler.ServeHTTP(rr, req)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// No need to check for unfulfilled expectations as we are not using sqlmock anymore
}
