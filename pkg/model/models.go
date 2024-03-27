package models

type LoginCredentials struct {
	MemberID      int
	Username      string
	Password      string
	Role_Name     int
	PastPasswords string
}

type Report struct {
	TypeOfIncident string
	Location       string
	Description    string
	Anonymous      string
	DeviceLocation string
	FilePath       string
}
