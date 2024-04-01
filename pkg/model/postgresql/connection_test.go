package postgresql

import (
	"encoding/base64"
	"net/http"
	"reflect"
	"testing"

	models "amencia.net/ubb-campus-safety-main/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestConnectModel_ReadReport(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of your ConnectModel
	model := &ConnectModel{DB: db}

	// Define expected rows and columns
	expectedRows := sqlmock.NewRows([]string{"person_name", "type_of_incident", "location", "description", "imagename", "imagedata"}).
		AddRow("John Doe", "Theft", "Location A", "Description A", "image1.jpg", []byte{0x00, 0x01, 0x02})

	// Set up expectations
	mock.ExpectQuery("SELECT person_name, type_of_incident, location, description, imagename, imagedata FROM report").
		WillReturnRows(expectedRows)

	// Call the function
	reports, err := model.ReadReport()
	if err != nil {
		t.Fatalf("error calling ReadReport: %v", err)
	}

	// Check if the returned data matches the expected result
	expectedReport := &models.Report{
		PersonName:       "John Doe",
		TypeOfIncident:   "Theft",
		Location:         "Location A",
		Description:      "Description A",
		ImageName:        "image1.jpg",
		ImageData:        []byte{0x00, 0x01, 0x02},
		EncodedImageData: base64.StdEncoding.EncodeToString([]byte{0x00, 0x01, 0x02}),
		MimeType:         http.DetectContentType([]byte{0x00, 0x01, 0x02}),
	}

	if !reflect.DeepEqual(reports[0], expectedReport) {
		t.Errorf("returned report %+v does not match expected %+v", reports[0], expectedReport)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
