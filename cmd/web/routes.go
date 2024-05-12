package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a variable to hold my middleware chain
	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeadersMiddleware,
	)
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
	mux.Get("/reports", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.reports))
	mux.Get("/panic", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.panic))
	mux.Get("/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.profile))
	mux.Post("/report-add", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createReport))
	mux.Post("/guard-report-add", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.guardcreateReport))
	mux.Get("/add-user", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.addNewuser))
	mux.Get("/add-contact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.addContact))
	mux.Post("/create-user", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createuser))
	mux.Post("/create-contact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createContact))
	mux.Get("/view-reports", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.viewreport))
	mux.Get("/view-contact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.viewContact))
	mux.Get("/guard-reports", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.guardreports))
	mux.Get("/view-log", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.viewlog))
	mux.Get("/check-in-out", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.checkinout))
	mux.Get("/my-contact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.viewMyContact))
	mux.Get("/add-mycontact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.addMyContact))
	mux.Post("/create-mycontact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createMyContact))
	mux.Post("/remove-mycontact", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.removeMyContact))
	mux.Get("/guard-view-report", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.view_report))
	mux.Get("/guard-profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.guard_profile))
	mux.Get("/admin-profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.admin_profile))
	mux.Post("/create-log", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createLog))
	mux.Post("/create-notice", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createnotice))
	mux.Get("/add-notice", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.addNotices))
	mux.Get("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Create a fileserver to serve our static content
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
