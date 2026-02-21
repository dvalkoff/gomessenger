package model

import "github.com/google/uuid"

type Contacts struct {
	contacts []Contact
}

type Contact struct {
	Id       uuid.UUID
	Nickname string
}
