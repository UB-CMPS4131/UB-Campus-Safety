package main_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"amencia.net/ubb-campus-safety-main/pkg/model/postgresql"
	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB is a mock implementation of the database interface for testing purposes
type MockDB struct {
	ubcs postgresql.ConnectModel
}

// Connect implements the Connect method of model.DatabaseModel interface
func (m *MockDB) Connect() (*sql.DB, error) {
	return m.Connect()
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
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	// Create a new instance of the application
	ubcs := &postgresql.ConnectModel{DB: db}

	// Expect a query to be executed during the test
	mock.ExpectQuery("SELECT password, role, memberID FROM LOGIN WHERE username = ?").
		WithArgs("mpit").
		WillReturnRows(sqlmock.NewRows([]string{"password", "role", "memberID"}).
			AddRow("mpit", 1, 1))

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
		err := ubcs.DB.QueryRow("SELECT password, role, memberID FROM LOGIN WHERE username = $1", username).Scan(&storedPassword, &role, &memberID)
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

	// Additional assertions can be made here to ensure correct behavior

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
