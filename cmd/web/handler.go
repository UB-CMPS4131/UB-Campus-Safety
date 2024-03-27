package main

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"

	"html/template"
	"log"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/login.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) verification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	// Fetch user credentials from the database
	var storedPassword string
	var role int
	err = app.ubcs.DB.QueryRow("SELECT password, role FROM LOGIN WHERE username = $1", username).Scan(&storedPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Validate the password
	if storedPassword != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Render specific template based on role
	switch role {
	case 1: // Assuming role 1 is for normal users
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case 2: // Assuming role 2 is for admin users
		http.Redirect(w, r, "/student", http.StatusSeeOther)
	case 3: // Assuming role 3 is for admin users
		http.Redirect(w, r, "/guard", http.StatusSeeOther)

	default:
		http.Error(w, "Invalid role", http.StatusInternalServerError)
	}
}

func (app *application) admin(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/dashboard.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) student(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/panic.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) guard(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/guard/dashboard.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) reports(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/reports.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) panic(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/panic.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/profile.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) call(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/call.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) createReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/reports", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	incident_type := r.PostForm.Get("type_of_incident")
	location := r.PostForm.Get("location")
	description := r.PostForm.Get("description")
	file_path := r.PostForm.Get("file_path")
	anonymous := r.PostForm.Get("is_anonymous")
	device_location := r.PostForm.Get("device_location")

	// check the web form fields to valadity
	errors := make(map[string]string)
	// check each field
	if strings.TrimSpace(incident_type) == "" {
		errors["type_of_incident"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(incident_type) > 50 {
		errors["type_of_incident"] = "This field is too long(maximum is 50 characters)"
	}
	if strings.TrimSpace(location) == "" {
		errors["location"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(location) > 100 {
		errors["location"] = "This field is too long(maximum is 100 characters)"
	}
	if strings.TrimSpace(description) == "" {
		errors["description"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(description) > 255 {
		errors["description"] = "This field is too long(maximum is 255 characters)"
	}
	if strings.TrimSpace(file_path) == "" {
		errors["file_path"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(file_path) > 255 {
		errors["file_path"] = "This field is too long(maximum is 255 characters)"
	}
	if strings.TrimSpace(anonymous) == "" {
		errors["is_anonymous"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(anonymous) > 255 {
		errors["is_anonymous"] = "This field is too long(maximum is 255 characters)"
	}
	if strings.TrimSpace(device_location) == "" {
		errors["device_location"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(device_location) > 255 {
		errors["device_location"] = "This field is too long(maximum is 255 characters)"
	}
	// check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for field, errorMsg := range errors {
			fmt.Fprintf(w, "- %s: %s\n", field, errorMsg)
		}
		return
	}
	_, err = app.ubcs.Insert(incident_type, location, description, file_path, anonymous, device_location)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reports", http.StatusSeeOther)
}

// // func (app *application) displayQuotation(w http.ResponseWriter, r *http.Request) {
// // 	q, err := app.quotes.Read()
// // 	if err != nil {
// // 		log.Println(err.Error())
// // 		http.Error(w,
// // 			http.StatusText(http.StatusInternalServerError),
// // 			http.StatusInternalServerError)
// // 		return
// // 	}
// // 	// Display the quotes using a template
// // 	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
// // 	if err != nil {
// // 		log.Println(err.Error())
// // 		http.Error(w,
// // 			http.StatusText(http.StatusInternalServerError),
// // 			http.StatusInternalServerError)
// 		return
// 	}
// 	err = ts.Execute(w, q)

// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w,
// 			http.StatusText(http.StatusInternalServerError),
// 			http.StatusInternalServerError)
// 		return
// 	}
// }
