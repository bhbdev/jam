package models

import (
	"strings"

	"gorm.io/gen/helper"
)

/*
type User struct {

		ID             uint64 `gorm:"primaryKey"`
		Username       string `gorm:"uniqueIndex"`
		FirstName 	   string
		LastName  	   string
		Email          string `gorm:"uniqueIndex"`
		Password 	   string
		IsAdmin        bool       `gorm:"index"`
		Status         UserStatus `gorm:"index"` // active, disabled

		CreatedAt   time.Time `gorm:"index"`
		UpdatedAt   time.Time `gorm:"index"`
		LastLoginAt time.Time `gorm:"index"`
	}
*/

func init() {
	generateModel(&ModelSpec{
		structName: "User",
		tableName:  "users",
		fields: []helper.Field{
			idFieldSpec,
			createdAtFieldSpec,
			deletedAtFieldSpec,
			updatedAtFieldSpec,
			&FieldSpec{name: "Username", typ: "string", gormTag: "uniqueIndex", comment: ""},
			&FieldSpec{name: "FirstName", typ: "string", gormTag: "", comment: ""},
			&FieldSpec{name: "LastName", typ: "string", gormTag: "", comment: ""},
			&FieldSpec{name: "Email", typ: "string", gormTag: "uniqueIndex", comment: ""},
			&FieldSpec{name: "Password", typ: "string", gormTag: "", comment: ""},
			&FieldSpec{name: "IsAdmin", typ: "bool", gormTag: "index", comment: ""},
			&FieldSpec{name: "Status", typ: "UserStatus", gormTag: "default:active", comment: "status: active, disabled"},
			&FieldSpec{name: "LastLoginAt", typ: "time.Time", gormTag: "index", comment: ""},
			&FieldSpec{name: "ProfileImage", typ: "string", gormTag: "", comment: ""},
		},
	})
}

// UserStatus options.
type UserStatus string

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

func (u *User) Validate() (errs map[string]string) {
	errs = make(map[string]string)

	if u.FirstName == "" {
		errs["FirstName"] = "First name is required"
	}
	if u.LastName == "" {
		errs["LastName"] = "Last name is required"
	}
	if u.Email == "" {
		errs["Email"] = "Email is required"
	}
	if !strings.ContainsRune(u.Email, '@') {
		errs["Email"] = "Invalid email address"
	}
	if u.Password == "" {
		errs["Password"] = "Password is required"
	}

	return
}
