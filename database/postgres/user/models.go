package user

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type DbUser struct {
	Id                     int       `db:"id"`
	Email                  string    `db:"email"`
	PasswordHash           string    `db:"password_hash"`
	FullName               string    `db:"full_name"`
	PhoneNumber            string    `db:"phone_number"`
	IsEmailConfirmed       bool      `db:"is_email_confirmed"`
	IsPhoneNumberConfirmed bool      `db:"is_phone_number_confirmed"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
	Roles                  []DbRole
	Profile                DbUserProfile
}

type DbRole struct {
	Id     int    `db:"id"`
	UserId int    `db:"user_id"`
	Name   string `db:"name"`
}

type DbUserProfile struct {
	Id           int            `db:"id"`
	UserId       int            `db:"user_id"`
	PersonalInfo DbPersonalInfo `db:"personal_info"`
	Settings     DbSettings     `db:"settings"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type DbPersonalInfo struct {
	HasOrganisation bool           `json:"hasOrganisation"`
	Organisation    DbOrganisation `json:"organisation"`
	HasEducation    bool           `json:"hasEducation"`
	Education       []DbEducation  `json:"education"`
}

func (p *DbPersonalInfo) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), p)
}

func (p *DbPersonalInfo) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type DbOrganisation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (o *DbOrganisation) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), o)
}

func (o *DbOrganisation) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type DbEducation struct {
	Degree string `json:"degree"`
	Course string `json:"course"`
	Group  string `json:"group"`
}

func (e *DbEducation) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), e)
}

func (e *DbEducation) Value() (driver.Value, error) {
	return json.Marshal(e)
}

type DbSettings struct {
}

func (s *DbSettings) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s *DbSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}
