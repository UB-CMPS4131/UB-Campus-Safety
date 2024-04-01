package common

import (
	models "amencia.net/ubb-campus-safety-main/pkg/model"
)

type ProfileDATA struct {
	DATA         []*models.Profile
	Notification []*models.Notification
}
