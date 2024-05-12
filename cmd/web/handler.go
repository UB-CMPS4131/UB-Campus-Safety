package main

import (
	"errors"
	"io"
	"strings"
	"time"

	"html/template"
	"log"
	"net/http"

	models "amencia.net/ubb-campus-safety-main/pkg/model"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already authenticated
	if app.isAuthenticated(r) {
		// Get the authenticated user ID
		userID := app.session.GetInt(r, "authenticatedUserID")
		if userID == 0 {
			// Invalid user ID, redirect to login
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Fetch user role, ID, and username
		username, loginID, role, memberID, err := app.users.FetchUserRoleAndIDAndUsername(userID)
		if err != nil {
			log.Println("ERROR: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		personName, err := app.users.FetchUserPersonName(memberID)
		if err != nil {
			log.Println("ERROR: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Set user information in the application struct
		app.Username = username
		app.PersonName = personName
		app.MemberID = memberID
		app.LoginID = loginID

		// Redirect to the appropriate route based on role
		switch role {
		case 1: // Assuming role 1 is for normal users
			http.Redirect(w, r, "/add-notice", http.StatusSeeOther)
		case 2: // Assuming role 2 is for admin users
			http.Redirect(w, r, "/student", http.StatusSeeOther)
		case 3: // Assuming role 3 is for guard users
			http.Redirect(w, r, "/guard", http.StatusSeeOther)
		default:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	// User is not authenticated, render the login form
	data := &templateData{
		IsAuthenticated: false, // Set IsAuthenticated to false for the login form
	}

	ts, err := template.ParseFiles("./ui/html/login.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) verification(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	// check the web form fields to validity
	errorsUser := make(map[string]string)
	id, err := app.users.Authenticate(username, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			errorsUser["default"] = "Email or Password is incorrect"
			// rerender the login form
			ts, err := template.ParseFiles("./ui/html/login.tmpl")
			if err != nil {
				log.Println("ERROR: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			err = ts.Execute(w, &templateData{
				ErrorsFromForm:  errorsUser,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
			})

			if err != nil {
				log.Println("ERROR: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	loginID, role, memberID, err := app.users.FetchUserRoleAndID(id)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	personName, err := app.users.FetchUserPersonName(memberID)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Set PersonName in the application struct
	app.PersonName = personName
	app.Username = username
	app.MemberID = memberID
	app.LoginID = loginID
	// Redirect to the appropriate route based on role
	switch role {
	case 1: // Assuming role 1 is for normal users
		app.session.Put(r, "authenticatedUserID", id)
		http.Redirect(w, r, "/add-notice", http.StatusSeeOther)
	case 2: // Assuming role 2 is for admin users
		app.session.Put(r, "authenticatedUserID", id)
		http.Redirect(w, r, "/student", http.StatusSeeOther)
	case 3: // Assuming role 3 is for guard users
		app.session.Put(r, "authenticatedUserID", id)
		http.Redirect(w, r, "/guard", http.StatusSeeOther)
	default:
		app.session.Put(r, "authenticatedUserID", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) addContact(w http.ResponseWriter, r *http.Request) {

	ts, err := template.ParseFiles("./ui/admin/addContact.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		Notifications:   not,
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})
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
		log.Println("ERROR:", err.Error())
	}

	err = ts.Execute(w, &templateData{
		Notifications:   not,
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		DATA:            profiles.DATA, // Pass the actual profile data here
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		DATA:            profiles.DATA, // Pass the actual profile data here
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	ts, err := template.ParseFiles("./ui/student/profile.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	err = ts.Execute(w, &templateData{
		Notifications:   not,
		DATA:            profiles.DATA, // Pass the actual profile data here
		IsAuthenticated: app.isAuthenticated(r),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
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
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) createuser(w http.ResponseWriter, r *http.Request) {

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
		ts, err := template.ParseFiles("./ui/admin/adduser.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

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
	err = ts.Execute(w, &templateData{
		Reports:         q,
		IsAuthenticated: app.isAuthenticated(r),
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) viewContact(w http.ResponseWriter, r *http.Request) {
	q, err := app.ubcs.ReadContact()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Display the quotes using a template
	ts, err := template.ParseFiles("./ui/student/call.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, &templateData{
		Contacts:        q,
		Notifications:   not,
		IsAuthenticated: app.isAuthenticated(r),
	})
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
	err = ts.Execute(w, &templateData{
		Reports:         q,
		IsAuthenticated: app.isAuthenticated(r),
	})

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

	err = ts.Execute(w, &templateData{
		Logs:            q,
		IsAuthenticated: app.isAuthenticated(r),
	})

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
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
func (app *application) createReport(w http.ResponseWriter, r *http.Request) {

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
		ts, err := template.ParseFiles("./ui/student/reports.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

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
		ts, err := template.ParseFiles("./ui/guard/guard-report.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

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
		ts, err := template.ParseFiles("./ui/guard/checkinout.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

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
		ts, err := template.ParseFiles("./ui/admin/addNotice.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

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

func (app *application) createContact(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pnumber := r.FormValue("number")
	email := r.FormValue("email")

	// check the web form fields for validity
	errors := make(map[string]string)
	// check each field

	if name = strings.TrimSpace(name); name == "" {
		errors["name"] = "This field cannot be left blank"
	} else if len(name) > 75 {
		errors["name"] = "This field is too long (maximum is 75 characters)"
	}
	if pnumber = strings.TrimSpace(pnumber); pnumber == "" {
		errors["number"] = "This field cannot be left blank"
	} else if len(pnumber) > 25 {
		errors["number"] = "This field is too long (maximum is 25 characters)"
	}
	if email = strings.TrimSpace(email); email == "" {
		errors["email"] = "This field cannot be left blank"
	} else if len(email) > 50 {
		errors["email"] = "This field is too long (maximum is 50 characters)"
	}
	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/admin/addContact.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		return
	}

	// Declare err variable
	var err error

	// Insert the report
	_, err = app.ubcs.InsertContact(name, pnumber, email)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/add-contact", http.StatusSeeOther)
}

func (app *application) addMyContact(w http.ResponseWriter, r *http.Request) {

	ts, err := template.ParseFiles("./ui/student/studentContact.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, &templateData{
		Notifications:   not,
		IsAuthenticated: app.isAuthenticated(r),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
func (app *application) viewMyContact(w http.ResponseWriter, r *http.Request) {
	// Get the login ID from the session or request, assuming app.Username holds the login ID.
	loginID := app.LoginID

	// Fetch the contact based on the login ID.
	q, err := app.ubcs.ReadMyContact(loginID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Fetch notifications based on the login ID.
	not, err := app.ubcs.Notification(app.Username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Parse the template file.
	ts, err := template.ParseFiles("./ui/student/mycontact.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Execute the template with the fetched data.
	err = ts.Execute(w, &templateData{
		MyContacts:      q,
		Notifications:   not,
		IsAuthenticated: app.isAuthenticated(r),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) createMyContact(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	pnumber := r.FormValue("number")
	email := r.FormValue("email")

	// check the web form fields for validity
	errors := make(map[string]string)
	// check each field

	if name = strings.TrimSpace(name); name == "" {
		errors["name"] = "This field cannot be left blank"
	} else if len(name) > 75 {
		errors["name"] = "This field is too long (maximum is 75 characters)"
	}
	if pnumber = strings.TrimSpace(pnumber); pnumber == "" {
		errors["number"] = "This field cannot be left blank"
	} else if len(pnumber) > 25 {
		errors["number"] = "This field is too long (maximum is 25 characters)"
	}
	if email = strings.TrimSpace(email); email == "" {
		errors["email"] = "This field cannot be left blank"
	} else if len(email) > 50 {
		errors["email"] = "This field is too long (maximum is 50 characters)"
	}
	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/student/studentContact.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
		})

		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return

		}
		return
	}

	// Declare err variable
	var err error
	var loginid = app.LoginID
	// Insert the report
	_, err = app.ubcs.InsertMyContact(loginid, name, pnumber, email)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/my-contact", http.StatusSeeOther)
}

func (app *application) removeMyContact(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("contactName")
	var contactID = app.LoginID

	// Check the web form fields for validity
	errors := make(map[string]string)

	if name = strings.TrimSpace(name); name == "" {
		errors["name"] = "This field cannot be left blank"
	} else if len(name) > 75 {
		errors["name"] = "This field is too long (maximum is 75 characters)"
	}

	if len(errors) > 0 {
		// Handle errors and return if there are validation issues
		return
	}

	// Remove the contact
	err := app.ubcs.RemoveMyContact(contactID, name)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/my-contact", http.StatusSeeOther)
}
