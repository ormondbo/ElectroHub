package domain

import (
	"time"
	"uuid"
)

type Gender uint8

const (
	GenderMale Gender = iota
	GenderFemale
	GenderHelicopter
	GenderOther
)

type Client struct {
	ID                uuid.UUID
	Name              string
	Surname           string
	Birthday          time.Time
	Gender            Gender
	Registration_date time.Time
	Address           Address
}
