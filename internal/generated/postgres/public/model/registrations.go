//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Registrations struct {
	UserID      int64
	TgContact   string
	TeamID      *string
	TeamName    string
	CaptainName string
	Phone       string
	GroupName   string
	Amount      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
