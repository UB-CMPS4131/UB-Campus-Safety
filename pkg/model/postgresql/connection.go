package postgresql

import (
	"database/sql"
)

type ConnectModel struct {
	DB *sql.DB
}

func (m *ConnectModel) Insert(incident_type, location, description, file_path, anonymous, device_location string) (int, error) {
	var id int

	s := `INSERT INTO Report (type_of_incident,location,description,file_path, is_anonymous,device_location)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING report_id`

	err := m.DB.QueryRow(s, incident_type, location, description, file_path, anonymous, device_location).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (m *Conn) Read() ([]*models., error) {
// 	// SQL statement
// 	s := `
// 		SELECT author_name, category quote
// 		FROM quotations
// 		LIMIT 10
// 	`
// 	rows, err := m.DB.Query(s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	quotes := []*models.Quote{}

// 	for rows.Next() {
// 		q := &models.Quote{}
// 		err = rows.Scan(&q.Author_name, &q.Category,
// 			&q.Body)

// 		if err != nil {
// 			return nil, err
// 		}
// 		quotes = append(quotes, q)
// 	}

// 	err = rows.Err()
// 	if err != nil {

// 		return nil, err
// 	}
// 	return quotes, nil
// }
