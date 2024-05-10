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

func (m *UserModel) FetchUserRoleAndID(id int) (int, int, error) {
	// Fetch user credentials from the database
	var role int
	var memberID int

	s := `
		SELECT role, memberID
		FROM LOGIN
		WHERE id = $1
		AND activated = TRUE
	`
	err := m.DB.QueryRow(s, id).Scan(&role, &memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, models.ErrInvalidCredentials
		} else {
			return 0, 0, err
		}
	}
	return role, memberID, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
