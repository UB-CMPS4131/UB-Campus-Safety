package main

import (
	"database/sql"
	"fmt"
	"strings"

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
	var data struct {
		ErrorMessage string
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		data.ErrorMessage = http.StatusText(http.StatusBadRequest)
		renderTemplate(w, "./ui/html/login.tmpl", data)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	// Fetch user credentials from the database
	var storedPassword string
	var role int
	var memberID int // Added memberID variable
	err = app.ubcs.DB.QueryRow("SELECT password, role, memberID FROM LOGIN WHERE username = $1", username).Scan(&storedPassword, &role, &memberID)
	if err == sql.ErrNoRows || storedPassword != password {
		data.ErrorMessage = "Invalid username or password"
		renderTemplate(w, "./ui/html/login.tmpl", data)
		return
	}

	// Fetch first name and last name from PersonnelInfoTable based on memberID
	var personName string
	row := app.ubcs.DB.QueryRow("SELECT CONCAT(FName, ' ', LName) FROM PersonnelInfoTable WHERE id = $1 LIMIT 1", memberID)
	err = row.Scan(&personName)
	if err == sql.ErrNoRows {
		data.ErrorMessage = "User not found in PersonnelInfoTable"
		renderTemplate(w, "./ui/html/login.tmpl", data)
		return
	}

	// Set PersonName in the application struct
	app.PersonName = personName

	// Render specific template based on role
	switch role {
	case 1: // Assuming role 1 is for normal users
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	case 2: // Assuming role 2 is for admin users
		http.Redirect(w, r, "/student", http.StatusSeeOther)
	case 3: // Assuming role 3 is for guard users
		http.Redirect(w, r, "/guard", http.StatusSeeOther)
	default:
		http.Error(w, "Invalid role", http.StatusInternalServerError)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	ts, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
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

func (app *application) guard_profile(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/guard/guard-profile.tmpl")
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
func (app *application) admin_profile(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/admin-profile.tmpl")
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
func (app *application) guardreports(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/guard/guard-reports.tmpl")
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

func (app *application) addNewuser(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/adduser.tmpl")
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

func (app *application) viewreport(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/viewreport.tmpl")
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

func (app *application) view_report(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/guard/view-report.tmpl")
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

func (app *application) viewlog(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/viewlog.tmpl")
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

func (app *application) checkinout(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/guard/checkinout.tmpl")
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	incidentType := r.PostForm.Get("type_of_incident")
	location := r.PostForm.Get("location")
	description := r.PostForm.Get("description")
	anonymous := r.PostForm.Get("is_anonymous")
	deviceLocation := r.PostForm.Get("device_location")

	// Default personName to Anonymous if anonymous is checked
	var personName string
	if anonymous == "on" {
		personName = "Anonymous"
	} else {
		personName = app.PersonName
	}

	// Check if file_path is provided
	var filePath string
	if _, ok := r.PostForm["file_path"]; ok {
		filePath = r.PostForm.Get("file_path")
	}

	// check the web form fields for validity
	errors := make(map[string]string)
	// check each field
	if incidentType = strings.TrimSpace(incidentType); incidentType == "" {
		errors["type_of_incident"] = "This field cannot be left blank"
	} else if len(incidentType) > 50 {
		errors["type_of_incident"] = "This field is too long (maximum is 50 characters)"
	}
	if location = strings.TrimSpace(location); location == "" {
		errors["location"] = "This field cannot be left blank"
	} else if len(location) > 100 {
		errors["location"] = "This field is too long (maximum is 100 characters)"
	}
	if description = strings.TrimSpace(description); description == "" {
		errors["description"] = "This field cannot be left blank"
	} else if len(description) > 255 {
		errors["description"] = "This field is too long (maximum is 255 characters)"
	}

	// Validate filePath if provided
	if filePath != "" && len(filePath) > 255 {
		errors["file_path"] = "This field is too long (maximum is 255 characters)"
	}

	if deviceLocation != "" && len(deviceLocation) > 255 {
		errors["device_location"] = "This field is too long (maximum is 255 characters)"
	}

	// check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for field, errorMsg := range errors {
			fmt.Fprintf(w, "- %s: %s\n", field, errorMsg)
		}
		return
	}
	// Insert the report
	_, err = app.ubcs.Insert(incidentType, personName, location, description, filePath, deviceLocation)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
