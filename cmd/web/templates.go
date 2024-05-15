// Names: Alex Peraza & Abner Mencia
// Assignment: Final Project
package main

import (
	"net/url"

	models "amencia.net/ubb-campus-safety-main/pkg/model"
)

type templateData struct {
	Contacts        []*models.Contact
	MyContacts      []*models.MyContact
	DATA            []*models.Profile
	Logs            []*models.Log
	Notifications   []*models.Notification
	Locations       []*models.Map
	Reports         []*models.Report
	ErrorsFromForm  map[string]string
	Flash           string
	FormData        url.Values
	IsAuthenticated bool
}
