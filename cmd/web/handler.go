package main

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

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
	app.Username = username
	app.MemberID = memberID

	// Render specific template based on role
	switch role {
	case 1: // Assuming role 1 is for normal users
		http.Redirect(w, r, "/add-notice", http.StatusSeeOther)
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

func (app *application) addNotices(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/admin/addNotice.tmpl")
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
	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		println("ERROR:" + err.Error())
	}

	err = ts.Execute(w, not)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) guard(w http.ResponseWriter, r *http.Request) {
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

func (app *application) reports(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/student/reports.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		println("ERROR:" + err.Error())
	}

	err = ts.Execute(w, not)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) guard_profile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the request context or session
	username := app.Username

	// Pass the username to the ReadProfile method to fetch the profile
	profiles, err := app.ubcs.ReadProfile(username, false)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if len(profiles.DATA) == 0 {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Display the profile using a template
	ts, err := template.ParseFiles("./ui/guard/guard-profile.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, profiles.DATA)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
func (app *application) admin_profile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the request context or session
	username := app.Username

	// Pass the username to the ReadProfile method to fetch the profile
	profiles, err := app.ubcs.ReadProfile(username, false)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if len(profiles.DATA) == 0 {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Display the profile using a template
	ts, err := template.ParseFiles("./ui/admin/admin-profile.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, profiles.DATA)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
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
	// Retrieve the username from the request context or session
	username := app.Username

	// Pass the username to the ReadProfile method to fetch the profile
	profiles, err := app.ubcs.ReadProfile(username, true)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if len(profiles.DATA) == 0 {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Display the profile using a template
	ts, err := template.ParseFiles("./ui/student/profile.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, profiles)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
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

func (app *application) createuser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/add-user", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(10 << 20) //10 MB limit for filesize
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	usertype := r.FormValue("usertype")
	username := r.FormValue("username")
	firstname := r.FormValue("fname")
	middlename := r.FormValue("mname")
	lastname := r.FormValue("lname")
	dob := r.FormValue("dob")
	gender := r.FormValue("gender")

	file, handler, err := r.FormFile("imagedata")
	if err != nil {
		http.Error(w, "ERROR REtRIEVING IMAGE", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imagedata, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error Reading File Data", http.StatusInternalServerError)
		return
	}

	imageName := handler.Filename

	//checking for errors map
	errors := make(map[string]string)
	//validation code
	if username = strings.TrimSpace(username); username == "" {
		errors["Username"] = "This Field connat be left"
	} else if len(firstname) > 25 {
		errors["Username"] = "This Field is too Long, max lenght is 25 characters"
	} else if len(firstname) < 3 {
		errors["Username"] = "This Field is Too short, minimum accepted lenght is 3 characters"
	}

	if firstname = strings.TrimSpace(firstname); firstname == "" {
		errors["First Name"] = "This Field connat be left"
	} else if len(firstname) > 25 {
		errors["First Name"] = "This Field is too Long, max lenght is 25 characters"
	} else if len(firstname) < 3 {
		errors["First Name"] = "This Field is Too short, minimum accepted lenght is 3 characters"
	}
	if middlename = strings.TrimSpace(middlename); firstname == "" {
		errors["Middle Name"] = "This Field connat be left"
	} else if len(firstname) > 25 {
		errors["Middle Name"] = "This Field is too Long, max lenght is 25 characters"
	} else if len(firstname) < 3 {
		errors["Middle Name"] = "This Field is Too short, minimum accepted lenght is 3 characters"
	}
	if lastname = strings.TrimSpace(lastname); firstname == "" {
		errors["Last Name"] = "This Field connat be left empty"
	} else if len(firstname) > 25 {
		errors["Last Name"] = "This Field is too Long, max lenght is 25 characters"
	} else if len(firstname) < 3 {
		errors["Last Name"] = "This Field is Too short, minimum accepted lenght is 3 characters"
	}
	roleINT := 0
	if usertype == "Student" {
		roleINT = 2
	} else if usertype == "Admin" {
		roleINT = 1
	} else if usertype == "Guard" {
		roleINT = 3
	} else {
		errors["Eser Type"] = "Invalid User Type"
	}

	// Validation code for date of birth (dob)
	if dob == "" {
		errors["dob"] = "Date of birth is required"
	} else {
		_, err := time.Parse("2006-01-02", dob)
		if err != nil {
			errors["dob"] = "Invalid date format for date of birth"
		}
	}

	// Validation code for gender
	if gender == "" {
		errors["gender"] = "Gender is required"
	} else if gender != "Male" && gender != "Female" && gender != "Other" {
		errors["gender"] = "Invalid gender specified"
	}

	if len(errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for field, errorMsg := range errors {
			fmt.Fprintf(w, "- %s: %s\n", field, errorMsg)
		}
		return
	}

	_, err = app.ubcs.NewUser(username, firstname, lastname, middlename, gender, dob, imagedata, imageName, roleINT)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/add-user", http.StatusSeeOther)

}

func (app *application) viewreport(w http.ResponseWriter, r *http.Request) {
	q, err := app.ubcs.ReadReport()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	// Display the quotes using a template
	ts, err := template.ParseFiles("./ui/admin/viewreport.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, q)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) view_report(w http.ResponseWriter, r *http.Request) {
	q, err := app.ubcs.ReadReport()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	// Display the quotes using a template
	ts, err := template.ParseFiles("./ui/guard/view-report.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, q)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) viewlog(w http.ResponseWriter, r *http.Request) {
	q, err := app.ubcs.ReadLog()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	// Display the quotes using a template
	ts, err := template.ParseFiles("./ui/admin/viewlog.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, q)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
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

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	incidentType := r.FormValue("type_of_incident")
	location := r.FormValue("location")
	description := r.FormValue("description")
	anonymous := r.FormValue("is_anonymous")
	deviceLocation := r.FormValue("device_location")

	// Default personName to Anonymous if anonymous is checked
	var personName string
	if anonymous == "on" {
		personName = "Anonymous"
	} else {
		personName = app.PersonName
	}

	var imageName string
	var imageData []byte

	// Check if file is provided
	file, handler, err := r.FormFile("file_path")
	if err == nil {
		defer file.Close()

		// Read the file data into a byte slice
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error Reading File Data", http.StatusInternalServerError)
			return
		}

		imageName = handler.Filename
	} else {
		// No image submitted by the user
		imageName = "No image submitted"
		imageData = []byte{} // Empty byte slice
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
	} else if len(description) > 700 {
		errors["description"] = "This field is too long (maximum is 700 characters)"
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
	_, err = app.ubcs.Insert(incidentType, personName, location, description, imageName, imageData, deviceLocation)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reports", http.StatusSeeOther)
}

func (app *application) guardcreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/guard-reports", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	incidentType := r.FormValue("type_of_incident")
	location := r.FormValue("location")
	description := r.FormValue("description")
	anonymous := r.FormValue("is_anonymous")
	deviceLocation := r.FormValue("device_location")

	// Default personName to Anonymous if anonymous is checked
	var personName string
	if anonymous == "on" {
		personName = "Anonymous"
	} else {
		personName = app.PersonName
	}

	var imageName string
	var imageData []byte

	// Check if file is provided
	file, handler, err := r.FormFile("file_path")
	if err == nil {
		defer file.Close()

		// Read the file data into a byte slice
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error Reading File Data", http.StatusInternalServerError)
			return
		}

		imageName = handler.Filename
	} else {
		// No image submitted by the user
		imageName = "No image submitted"
		imageData = []byte{} // Empty byte slice
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
	} else if len(description) > 700 {
		errors["description"] = "This field is too long (maximum is 700 characters)"
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
	_, err = app.ubcs.Insert(incidentType, personName, location, description, imageName, imageData, deviceLocation)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/guard-reports", http.StatusSeeOther)
}

func (app *application) createLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/check-in-out", http.StatusSeeOther)
		return
	}

	logDate := r.FormValue("date")
	logTime := r.FormValue("time")
	checkType := r.FormValue("checkType")

	// Assign personName directly without checking for anonymous
	personName := app.PersonName

	// check the web form fields for validity
	errors := make(map[string]string)
	// check each field
	if logDate = strings.TrimSpace(logDate); logDate == "" {
		errors["date"] = "This field cannot be left blank"
	} else if len(logDate) > 50 {
		errors["date"] = "This field is too long (maximum is 50 characters)"
	}
	if logTime = strings.TrimSpace(logTime); logTime == "" {
		errors["time"] = "This field cannot be left blank"
	} else if len(logTime) > 100 {
		errors["time"] = "This field is too long (maximum is 100 characters)"
	}
	if checkType = strings.TrimSpace(checkType); checkType == "" {
		errors["checkType"] = "This field cannot be left blank"
	} else if len(checkType) > 100 {
		errors["checkType"] = "This field is too long (maximum is 100 characters)"
	}
	// check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for field, errorMsg := range errors {
			fmt.Fprintf(w, "- %s: %s\n", field, errorMsg)
		}
		return
	}

	// Declare err variable
	var err error

	// Insert the report
	_, err = app.ubcs.Insertlog(personName, logDate, logTime, checkType)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/check-in-out", http.StatusSeeOther)
}

func (app *application) createnotice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/add-notice", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	message := r.FormValue("message")

	userid := app.MemberID

	// check the web form fields for validity
	errors := make(map[string]string)
	// check each field
	if title = strings.TrimSpace(title); title == "" {
		errors["title"] = "This field cannot be left blank"
	} else if len(title) > 50 {
		errors["title"] = "This field is too long (maximum is 50 characters)"
	}
	if message = strings.TrimSpace(message); message == "" {
		errors["message"] = "This field cannot be left blank"
	} else if len(message) > 1500 {
		errors["message"] = "This field is too long (maximum is 1500 characters)"
	}
	// check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for field, errorMsg := range errors {
			fmt.Fprintf(w, "- %s: %s\n", field, errorMsg)
		}
		return
	}

	// Declare err variable
	var err error

	// Insert the report
	_, err = app.ubcs.InsertNotice(userid, title, message)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/add-notice", http.StatusSeeOther)
}
