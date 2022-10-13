package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type Role int
type Sex int

const (
	SP Role = iota + 1
	Manager
	SuperUser
)

const (
	Male Sex = iota + 1
	Female
)

type Name struct {
	First string `json:"First"`
	Last  string `json:"Last"`
}

func (n Name) Create(full_name string) *Name {
	fullname_split := strings.Split(full_name, " ")
	first, last := fullname_split[0], fullname_split[1]
	return &Name{
		First: first,
		Last:  last,
	}
}

func (n Name) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Name) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &n)
}

type Instructor struct {
	Name  Name   `json:"Name"`
	Title string `json:"Title"`
}

func (i Instructor) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Instructor) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &i)
}

func (i Instructor) Create(full_name string, title string) *Instructor {
	return &Instructor{
		Name:  *Name{}.Create(full_name),
		Title: title,
	}
}

type SpUser struct {
	Id                  uint           `json:"Id"`
	Name                Name           `json:"Name"`
	Username            string         `json:"Username"`
	Role                Role           `json:"Role"`
	Sex                 Sex            `json:"Sex"`
	SessionsAvailable   []*SessionInfo `json:"SessionsAvailable"`
	SessionsUnavailable []*SessionInfo `json:"SessionsUnavailable"`
	SessionsAssigned    []*SessionInfo `json:"SessionsAssigned"`
	Password            string         `json:"Password"`
	Email               string         `json:"Email"`
}

func (sp SpUser) Value() (driver.Value, error) {
	return json.Marshal(sp)
}

func (sp *SpUser) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &sp)
}

func (spUser SpUser) Create(name Name, role Role, sex Sex, email string) *SpUser {
	return &SpUser{
		Name:                name,
		Role:                role,
		Sex:                 sex,
		SessionsAvailable:   make([]*SessionInfo, 0),
		SessionsUnavailable: make([]*SessionInfo, 0),
		SessionsAssigned:    make([]*SessionInfo, 0),
		Password:            "",
		Email:               email,
	}
}

type SpManager struct {
	Id                uint       `json:"Id"`
	Name              Name       `json:"Name"`
	Username          string     `json:"Username"`
	Role              Role       `json:"Role"`
	AssignedPatients  []*SpUser  `json:"AssignedPatients"`
	SessionsManaged   []*Session `json:"SessionsManaged"`
	SessionsUnmanaged []*Session `json:"SessionsUnmanaged"`
	Password          string     `json:"Password"`
	Email             string     `json:"Email"`
}

func (sp SpManager) Value() (driver.Value, error) {
	return json.Marshal(sp)
}

func (sp *SpManager) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &sp)
}

func (spManager SpManager) Create(name Name, role Role, email string) *SpManager {
	return &SpManager{
		Name:              name,
		Role:              role,
		AssignedPatients:  []*SpUser{},
		SessionsManaged:   []*Session{},
		SessionsUnmanaged: []*Session{},
		Password:          "",
		Email:             email,
	}
}

type Session struct {
	Id                  uint          `json:"Id"`
	Information         *SessionInfo  `json:"Information"`
	Instructors         []*Instructor `json:"Instructors"`
	PatientsNeeded      int           `json:"PatientsNeeded"`
	PatientsAssigned    []*SpUser     `json:"PatientsAssigned"`
	PatientsAvailable   []*SpUser     `json:"PatientsAvailable"`
	PatientsUnavailable []*SpUser     `json:"PatientsUnavailable"`
	PatientsNoResponse  []*SpUser     `json:"PatientsNoResponse"`
}

func (s Session) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Session) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}

func (s Session) Create(date string, time string, duration string, location string, description string) *Session {
	return &Session{
		Information: &SessionInfo{Date: date,
			Time:        time,
			Duration:    duration,
			Location:    location,
			Description: description},
		Instructors:         []*Instructor{},
		PatientsNeeded:      0,
		PatientsAssigned:    []*SpUser{},
		PatientsAvailable:   []*SpUser{},
		PatientsUnavailable: []*SpUser{},
		PatientsNoResponse:  []*SpUser{},
	}
}

type SessionInfo struct {
	Date        string `json:"Date"`
	Time        string `json:"Time"`
	Duration    string `json:"Duration"`
	Location    string `json:"Location"`
	Description string `json:"Description"`
}

func (si SessionInfo) Value() (driver.Value, error) {
	return json.Marshal(si)
}

func (si *SessionInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &si)
}
