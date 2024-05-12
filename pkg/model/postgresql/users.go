// Names: Alex Peraza & Abner Mencia
// Assignment: Final Project
package postgresql

import (
	"database/sql"
	"errors"

	models "amencia.net/ubb-campus-safety-main/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Authenticate(username, password string) (int, error) {
	var id int
	var hashedPassword []byte

	s := `
		SELECT id, password
		FROM LOGIN
		WHERE username = $1
		AND activated = TRUE
	`
	err := m.DB.QueryRow(s, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// check the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) FetchUserRoleAndIDAndUsername(id int) (string, int, int, int, error) {
	// Fetch user credentials from the database
	var username string
	var role int
	var memberID int
	var LoginID int

	s := `
		SELECT id , username, role, memberID
		FROM LOGIN
		WHERE id = $1
		AND activated = TRUE
	`
	err := m.DB.QueryRow(s, id).Scan(&LoginID, &username, &role, &memberID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "username not found", 0, 0, 0, models.ErrInvalidCredentials
		} else {
			return "username not found", 0, 0, 0, err
		}
	}
	return username, LoginID, role, memberID, nil
}
func (m *UserModel) FetchUserRoleAndID(id int) (int, int, int, error) {
	// Fetch user credentials from the database
	var role int
	var memberID int
	var LoginID int

	s := `
		SELECT id, role, memberID
		FROM LOGIN
		WHERE id = $1
		AND activated = TRUE
	`
	err := m.DB.QueryRow(s, id).Scan(&LoginID, &role, &memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, 0, models.ErrInvalidCredentials
		} else {
			return 0, 0, 0, err
		}
	}
	return LoginID, role, memberID, nil
}

func (m *UserModel) FetchUserPersonName(memberID int) (string, error) {
	// Fetch first name and last name from PersonnelInfoTable based on memberID

	var personName string

	s := `
	SELECT CONCAT(FName, ' ', LName) FROM PersonnelInfoTable WHERE id = $1 LIMIT 1
	`
	err := m.DB.QueryRow(s, memberID).Scan(&personName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "User information not found", models.ErrInvalidCredentials
		} else {
			return "User information not found", err
		}
	}
	return personName, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
