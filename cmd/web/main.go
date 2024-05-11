package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"amencia.net/ubb-campus-safety-main/pkg/model/postgresql"
	"github.com/bmizerany/pat"
	"github.com/golangcollege/sessions"
	"github.com/justinas/alice"
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
	users      *postgresql.UserModel
	PersonName string
	session    *sessions.Session
	Username   string
	MemberID   int
	LoginID    int
}

func main() {
	var db, err = setUpDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	secret := flag.String("secret", "p7Mhd+qQamgHsS*+8Tg7mNXtcjvu@egz", "Secret key")
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	app := &application{
		ubcs: &postgresql.ConnectModel{
			DB: db,
		},
		session: session,
		users: &postgresql.UserModel{
			DB: db,
		},
	}

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	// Serve HTML files
	htmlDir := http.Dir("ui/html")
	htmlFS := http.FileServer(htmlDir)
	mux.Get("/html/", http.StripPrefix("/html/", htmlFS))

	jsDir := http.Dir("ui/js")
	jsFS := http.FileServer(jsDir)
	mux.Get("/js/", http.StripPrefix("/js/", jsFS))

	imgDir := http.Dir("ui/images")
	imgFS := http.FileServer(imgDir)
	mux.Get("/images/", http.StripPrefix("/images/", imgFS))

	// Serve CSS files
	cssDir := http.Dir("ui/css")
	cssFS := http.FileServer(cssDir)
	mux.Get("/css/", http.StripPrefix("/css/", cssFS))

	mux.Get("/", dynamicMiddleware.ThenFunc(app.login))
	mux.Post("/login", dynamicMiddleware.ThenFunc(app.verification))
	mux.Get("/student", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.student))
	mux.Get("/guard", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.guard))
	mux.Get("/reports", dynamicMiddleware.ThenFunc(app.reports))
	mux.Get("/panic", dynamicMiddleware.ThenFunc(app.panic))
	mux.Get("/profile", dynamicMiddleware.ThenFunc(app.profile))
	mux.Post("/report-add", dynamicMiddleware.ThenFunc(app.createReport))
	mux.Post("/guard-report-add", dynamicMiddleware.ThenFunc(app.guardcreateReport))
	mux.Get("/add-user", dynamicMiddleware.ThenFunc(app.addNewuser))
	mux.Get("/add-contact", dynamicMiddleware.ThenFunc(app.addContact))
	mux.Post("/create-user", dynamicMiddleware.ThenFunc(app.createuser))
	mux.Post("/create-contact", dynamicMiddleware.ThenFunc(app.createContact))
	mux.Get("/view-reports", dynamicMiddleware.ThenFunc(app.viewreport))
	mux.Get("/view-contact", dynamicMiddleware.ThenFunc(app.viewContact))
	mux.Get("/guard-reports", dynamicMiddleware.ThenFunc(app.guardreports))
	mux.Get("/view-log", dynamicMiddleware.ThenFunc(app.viewlog))
	mux.Get("/check-in-out", dynamicMiddleware.ThenFunc(app.checkinout))
	mux.Get("/my-contact", dynamicMiddleware.ThenFunc(app.viewMyContact))
	mux.Get("/add-mycontact", dynamicMiddleware.ThenFunc(app.addMyContact))
	mux.Post("/create-mycontact", dynamicMiddleware.ThenFunc(app.createMyContact))
	mux.Post("/remove-mycontact", dynamicMiddleware.ThenFunc(app.removeMyContact))
	mux.Get("/guard-view-report", dynamicMiddleware.ThenFunc(app.view_report))
	mux.Get("/guard-profile", dynamicMiddleware.ThenFunc(app.guard_profile))
	mux.Get("/admin-profile", dynamicMiddleware.ThenFunc(app.admin_profile))
	mux.Post("/create-log", dynamicMiddleware.ThenFunc(app.createLog))
	mux.Post("/create-notice", dynamicMiddleware.ThenFunc(app.createnotice))
	mux.Get("/add-notice", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.addNotices))
	mux.Get("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	log.Println("Starting server on port :8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
