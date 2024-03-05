package domain

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	NoneRole     Role = "none"
	DirectorRole Role = "director"
	ExpertRole   Role = "expert"
)

type Director struct {
	ID   uuid.UUID
	User User
}

type Expert struct {
	ID          uuid.UUID
	isConfirmed bool
	User        User
}

type User struct {
	ID         uuid.UUID
	Username   string
	Password   string
	RegisterAt time.Time
	UserRole   string
}

type ConfirmExpert struct {
	ExpertID    string
	IsConfirmed bool
}
