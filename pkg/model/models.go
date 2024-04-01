package models

import (
	"time"
)

type LoginCredentials struct {
	MemberID      int
	Username      string
	Password      string
	Role_Name     int
	PastPasswords string
}

type Report struct {
	PersonName       string
	TypeOfIncident   string
	Location         string
	Description      string
	ImageName        string
	ImageData        []byte // Assuming imagedata is stored as bytea in the database
	EncodedImageData string // EncodedImageData will hold the base64 encoded string
	MimeType         string // MimeType will store the detected MIME type of the image data
}

type Profile struct {
	ID           int
	Image        string
	ImageData    []byte
	Fname        string
	Mname        string
	LName        string
	DOB          time.Time
	Gender       string
	EncodedImage string
	MimeType     string
}

type Log struct {
	PersonName string
	LogDate    time.Time
	LogTime    time.Time
	CheckType  string
}

type Notification struct {
	Notificationid int
	UserID         int
	Title          string
	Message        string
	Created_at     time.Time
}

type Notification_Seen struct {
	notification_id int
	user_id         int
	seen_at         time.Time
}
