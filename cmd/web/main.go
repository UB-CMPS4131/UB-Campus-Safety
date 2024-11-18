package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"amencia.net/ubb-campus-safety-main/pkg/model/postgresql"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq" // Third party package
)

func setUpDB(dsn string) (*sql.DB, error) {
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
	errorLog   *log.Logger
	infoLog    *log.Logger
	Username   string
	MemberID   int
	LoginID    int
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	//using 5433 since im also using 5432 for local postgers, therefore 5433 in public maps to 5432 inside docker
	dsn := flag.String("dsn", "postgres://ub:ub@localhost:5433/ub?sslmode=disable", "PostgreSQL DSN (Data Source Name)")
	flag.Parse()
	secret := flag.String("secret", "p7Mhd+qQamgHsS*+8Tg7mNXtcjvu@egz", "Secret key")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	var db, err = setUpDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	session := sessions.New([]byte(*secret))
	session.Lifetime = 4380 * time.Hour // 6 months in hours (assuming 30 days per month)
	session.Secure = true

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		ubcs: &postgresql.ConnectModel{

			DB: db,
		},
		session: session,
		users: &postgresql.UserModel{
			DB: db,
		},
	}
	// Create a custom web server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start our server
	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
