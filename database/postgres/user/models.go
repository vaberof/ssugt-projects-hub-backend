package user

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type DbUser struct {
	Id           int       `db:"id"`
	RoleId       int       `db:"role_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Profile      DbUserProfile
}

type DbUserProfile struct {
	Id           int            `db:"id"`
	UserId       int            `db:"user_id"`
	PersonalInfo DbPersonalInfo `db:"personal_info"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type DbPersonalInfo struct {
	HasOrganisation bool           `json:"hasOrganisation"`
	Organisation    DbOrganisation `json:"organisation"`
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
