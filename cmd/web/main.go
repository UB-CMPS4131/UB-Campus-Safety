package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"amencia.net/ubb-campus-safety-main/pkg/model/postgresql"
	_ "github.com/lib/pq" // Third party package
)

func setUpDB() (*sql.DB, error) {
	// Provide the credentials for our database
	const (
		host     = "bubble.db.elephantsql.com"
		port     = 5432
		user     = "xqymnerr"
		password = "Xgtj9QRe3ouBnLW1WN-9C_g4_DDWefMr"
		dbname   = "xqymnerr"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Test our connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Dependencies (things/variables)
// Dependency Injection (passing)
type application struct {
	ubcs       *postgresql.ConnectModel
	PersonName string
	Username   string
	MemberID   int
}

func main() {
	var db, err = setUpDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	app := &application{
		ubcs: &postgresql.ConnectModel{
			DB: db,
		},
	}

	mux := http.NewServeMux()
	// Serve HTML files
	htmlDir := http.Dir("ui/html")
	htmlFS := http.FileServer(htmlDir)
	mux.Handle("/html/", http.StripPrefix("/html/", htmlFS))

	jsDir := http.Dir("ui/js")
	jsFS := http.FileServer(jsDir)
	mux.Handle("/js/", http.StripPrefix("/js/", jsFS))

	imgDir := http.Dir("ui/images")
	imgFS := http.FileServer(imgDir)
	mux.Handle("/images/", http.StripPrefix("/images/", imgFS))

	// Serve CSS files
	cssDir := http.Dir("ui/css")
	cssFS := http.FileServer(cssDir)
	mux.Handle("/css/", http.StripPrefix("/css/", cssFS))

	mux.HandleFunc("/", app.login)
	mux.HandleFunc("/login", app.verification)
	mux.HandleFunc("/student", app.student)
	mux.HandleFunc("/guard", app.guard)
	mux.HandleFunc("/reports", app.reports)
	mux.HandleFunc("/panic", app.panic)
	mux.HandleFunc("/profile", app.profile)
	mux.HandleFunc("/report-add", app.createReport)
	mux.HandleFunc("/guard-report-add", app.guardcreateReport)
	mux.HandleFunc("/add-user", app.addNewuser)
	mux.HandleFunc("/add-contact", app.addContact)
	mux.HandleFunc("/create-user", app.createuser)
	mux.HandleFunc("/create-contact", app.createContact)
	mux.HandleFunc("/view-reports", app.viewreport)
	mux.HandleFunc("/view-contact", app.viewContact)
	mux.HandleFunc("/guard-reports", app.guardreports)
	mux.HandleFunc("/view-log", app.viewlog)
	mux.HandleFunc("/check-in-out", app.checkinout)
	mux.HandleFunc("/guard-view-report", app.view_report)
	mux.HandleFunc("/guard-profile", app.guard_profile)
	mux.HandleFunc("/admin-profile", app.admin_profile)
	mux.HandleFunc("/create-log", app.createLog)
	mux.HandleFunc("/create-notice", app.createnotice)
	mux.HandleFunc("/add-notice", app.addNotices)

	log.Println("Starting server on port :8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
